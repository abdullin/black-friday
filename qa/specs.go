package qa

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	for _, r := range tests {

		q := NewQA(svc)
		r(q)

		if len(q.fails) == 0 {
			continue
		}

		fmt.Println("QA: I got an error when doing " + q.text)

		for i, s := range q.steps {

			fmt.Printf(" %d. %s\n", i+1, s)
		}

		for _, f := range q.fails {
			fmt.Println(f)
		}

		return
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

func (q *QA) AddProduct(sku string) ProductID {

	q.step("add SKU %s", sku)
	prod, _ := q.service.AddProduct(nil, &AddProductReq{Sku: sku})

	result := ProductID(prod.Id)

	q.producs[result] = sku
	return result

}

func (q *QA) AddLoc(name string) LocationID {

	q.step("add location %s", name)
	loc, _ := q.service.AddLocation(nil, &AddLocationReq{Name: name})

	q.locs[LocationID(loc.Id)] = name
	return LocationID(loc.Id)
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

func (q *QA) expectAddProductErr(sku string, code codes.Code) {
	q.step("add SKU %s", sku)
	prod, err := q.service.AddProduct(nil, &AddProductReq{Sku: sku})

	if prod != nil {
		q.fail("expected no response, but got it")
	}
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
}

func additive_quantity(q *QA) {

	q.title("check if quantity is added properly")

	cola := q.AddProduct("cola")
	shelf := q.AddLoc("shelf")

	q.UpdateQty(shelf, cola, 2)
	q.UpdateQty(shelf, cola, 3)

	q.expectInventory(shelf, map[ProductID]int64{cola: 5})
}

func product_names_are_unique(q *QA) {
	q.title("We can't have duplicate product names")

	_ = q.AddProduct("milk")

	q.expectAddProductErr("milk", codes.AlreadyExists)
}

func negative_qty(q *QA) {
	q.title("quantity can't go negative")

	fanta := q.AddProduct("fanta")
	bar := q.AddLoc("bar")
	q.expectUpdateQtyError(bar, fanta, -1, codes.InvalidArgument)
}
