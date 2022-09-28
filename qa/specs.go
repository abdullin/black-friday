package qa

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	. "sdk-go/protos"
	"strings"
)

// TODO: have unique name generator

type Rpc[T any, K any] func(ctx context.Context, req T) K

type QAContext struct {
	problems []string
}

func (q *QAContext) f(format string, a ...any) {
	q.problems = append(q.problems, fmt.Sprintf(format, a))
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

func RunCommandDrivenSpec(svc InventoryServiceServer) {

	q := NewQA(svc)
	teees(q)

	if len(q.fails) == 0 {
		return
	}

	fmt.Println("QA: I got an error when doing " + q.text)

	for i, s := range q.steps {

		fmt.Printf(" %d. %s\n", i+1, s)
	}

	for _, f := range q.fails {
		fmt.Println(f)
	}

}

type QA struct {
	service InventoryServiceServer
	text    string

	steps []string

	fails []string

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

func (q *QA) givenProduct(name string) ProductID {

	q.step("add product %s", name)
	prod, _ := q.service.AddProduct(nil, &AddProductReq{Name: name})

	result := ProductID(prod.Id)

	q.producs[result] = name
	return result

}

func (q *QA) givenLoc(name string) LocationID {

	q.step("add location %s", name)
	loc, _ := q.service.AddLocation(nil, &AddLocationReq{Name: name})

	q.locs[LocationID(loc.Id)] = name
	return LocationID(loc.Id)
}

func (q *QA) givenQty(p ProductID, l LocationID, qt int64) int64 {

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

func (q *QA) assertInventory(l LocationID, vals map[ProductID]int64) {

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
			q.fail("not found %s in stock", q.producs[i], i)
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

func teees(q *QA) {

	q.title("check if quantity is added properly")

	p := q.givenProduct("cola")
	l1 := q.givenLoc("Shelf")

	q.givenQty(p, l1, 2)
	q.givenQty(p, l1, 3)

	q.assertInventory(l1, map[ProductID]int64{p: 6})
}
