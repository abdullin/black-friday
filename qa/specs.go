package qa

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"os"
	"runtime/debug"
	. "sdk-go/protos"
	"strings"
)

type Rpc[T any, K any] func(ctx context.Context, req T) K

type QAContext struct {
	problems []string
}

func (q *QAContext) f(format string, a ...any) {
	q.problems = append(q.problems, fmt.Sprintf(format, a...))
}

func (q *QAContext) stop() bool {
	return len(q.problems) > 0
}

func (q *QAContext) Problems() []string {
	return q.problems
}

func (q *QAContext) assert(resp proto.Message, err error) bool {

	if resp == nil {
		q.f("got no response")
	}
	if err != nil {
		q.f("got error: %s", err.Error())
	}

	return q.stop()
}

func Dispose(e io.Closer) {
	err := e.Close()
	if err != nil {
		log.Println(err.Error())
	}
}

func RunCommandDrivenSpec(svc InventoryServiceServer) {

	os.TempDir()

	f, err := os.Create("/tmp/data.txt")
	if err != nil {
		panic(err)
	}

	defer Dispose(f)

	defer func() {
		if r := recover(); r != nil {
			f.WriteString(fmt.Sprintln("Your code doesn't run on my machine :( ", r, "\n```\n", string(debug.Stack()), "\n```"))
		}
	}()

	for _, r := range tests {

		q := NewQA(svc)
		r(q)

		if len(q.fails) == 0 {
			continue
		}

		f.WriteString("Hey! I got error with scenario '" + q.text + "'\n\n")

		for i, s := range q.steps {

			f.WriteString(fmt.Sprintf(" %d. %s\n", i+1, s))
		}

		if len(q.fails) > 0 {
			f.WriteString("\n")
		}

		for _, fail := range q.fails {
			f.WriteString(fail + "\n")
		}
	}
}

type Result struct {
	Text  string   `json:"text,omitempty"`
	Steps []string `json:"steps,omitempty"`
	Fails []string `json:"fails,omitempty"`
}

type QA struct {
	service InventoryServiceServer
	text    string

	steps   []string
	fails   []string
	locs    map[LocationID]string
	producs map[ProductID]string
}

func NewQA(svc InventoryServiceServer) *QA {
	return &QA{
		service: svc,
		producs: map[ProductID]string{},
		locs:    map[LocationID]string{},
	}
}

type ProductID uint64
type LocationID uint64

func typed[T ProductID | LocationID](values []uint64) []T {
	if values == nil {
		return nil
	}
	results := make([]T, len(values))
	for i, id := range values {
		results[i] = T(id)
	}
	return results
}

func (q *QA) AddProducts(skus ...string) ([]ProductID, error) {

	q.step("add SKU %s", strings.Join(skus, ", "))
	prod, err := q.service.AddProducts(nil, &AddProductsReq{Skus: skus})
	if err != nil {
		return nil, err
	}

	return typed[ProductID](prod.Ids), err
}

func (q *QA) AddLocs(names ...string) ([]LocationID, error) {

	q.step("add locations %s", strings.Join(names, ", "))
	locs, err := q.service.AddLocations(nil, &AddLocationsReq{Names: names})
	if err != nil {
		return nil, err
	}
	return typed[LocationID](locs.Ids), err
}

func (q *QA) UpdateQty(l LocationID, p ProductID, qt int64) int64 {

	if qt > 0 {

		q.step("put %d %s at %s", qt, q.producs[p], q.locs[l])
	} else {

		q.step("remove %d %s from %s", qt, q.producs[p], q.locs[l])
	}

	qty, _ := q.service.UpdateQty(nil, &UpdateQtyReq{
		Location: uint64(l),
		Product:  uint64(p),
		Quantity: qt,
	})

	return qty.Total
}

func (q *QA) expectInventory(l LocationID, vals map[ProductID]int64) {

	lines := []string{}

	for i, v := range vals {
		lines = append(lines, fmt.Sprintf("%d x %s", v, q.producs[i]))
	}

	q.step("check inventory at %s: %s", q.locs[l], strings.Join(lines, ", "))

	resp, _ := q.service.GetInventory(nil, &GetInventoryReq{Location: uint64(l)})

	counts := map[uint64]int64{}

	for _, line := range resp.Items {
		counts[line.Product] = line.Quantity
	}

	for i, expected := range vals {

		actual, found := counts[uint64(i)]
		if !found {
			q.fail("not found %s in stock", q.producs[i])
		} else {
			if actual != expected {
				q.fail("expected %d x %s to be in stock at %s but got %d", expected, q.producs[i], q.locs[l], actual)
			}
		}

	}

}

func (q *QA) title(s string) {
	q.text = s
}

func (q *QA) step(format string, args ...any) {
	q.steps = append(q.steps, fmt.Sprintf(format, args...))
}

func (q *QA) fail(format string, args ...any) {

	stepNum := len(q.steps)

	err := fmt.Sprintf("Problem at step %d: ", stepNum)

	q.fails = append(q.fails, err+fmt.Sprintf(format, args...))

}

func (q *QA) expectUpdateQtyError(l LocationID, p ProductID, qt int64, c codes.Code) {
	if qt > 0 {

		q.step("put %d %s at %s", qt, q.producs[p], q.locs[l])
	} else {
		q.step("remove %d %s from %s", -qt, q.producs[p], q.locs[l])
	}

	qty, err := q.service.UpdateQty(nil, &UpdateQtyReq{
		Location: uint64(l),
		Product:  uint64(p),
		Quantity: qt,
	})

	if qty != nil {
		q.fail("expected no response, but got result with quantity %d", qty.Total)
	}
	if err == nil {
		q.fail("Expected error, but got nothing")
	} else {
		st, ok := status.FromError(err)
		if !ok {
			q.fail("got unexpected error %s", err.Error())
		} else {
			if st.Code() != c {
				q.fail("expected error %v got %v", c, st.Code())
			}

		}
	}

}
func (q *QA) expectOK(result any, err error) {
	if err != nil {
		q.fail("expected no error, but got %s", err.Error())
		return
	}
	if result == nil {
		q.fail("Expected result but got nothing")
	}
}

func (q *QA) expectErr(err error, code codes.Code) {
	if err == nil {
		q.fail("Expected error %v, but got nothing", code)
	} else {
		st, ok := status.FromError(err)
		if !ok {
			q.fail("got unexpected error %s", err.Error())
		} else {
			if st.Code() != code {
				q.fail("expected error %v got %v", code, st.Code())
			}

		}
	}

}

type Test func(q *QA)

var tests = []Test{
	additive_quantity,
	negative_qty,
	product_names_are_unique,
	batch_operation_rolls_back,
}

func additive_quantity(q *QA) {

	q.title("check if quantity is added properly")

	cola, _ := q.AddProducts("cola")
	shelf, _ := q.AddLocs("shelf")

	q.UpdateQty(shelf[0], cola[0], 2)
	q.UpdateQty(shelf[0], cola[0], 3)

	q.expectInventory(shelf[0], map[ProductID]int64{cola[0]: 5})
}

func product_names_are_unique(q *QA) {
	q.title("We can't have duplicate product names")
	_, err := q.AddProducts("milk", "milk")
	q.expectErr(err, codes.AlreadyExists)
}

func batch_operation_rolls_back(q *QA) {
	// good candidate for CH1
	q.title("batch operation rolls back as one")
	_, _ = q.AddProducts("water", "water")

	res, err := q.AddProducts("water")
	q.expectOK(res, err)

}

func negative_qty(q *QA) {
	q.title("quantity can't go negative")

	fanta, _ := q.AddProducts("fanta")
	bar, _ := q.AddLocs("bar")
	q.expectUpdateQtyError(bar[0], fanta[0], -1, codes.FailedPrecondition)
}
