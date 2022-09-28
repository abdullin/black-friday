package qa

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	. "sdk-go/protos"
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

func RunCommandDrivenSpec(svc InventoryServiceServer, qa *QAContext) {

	// what about dropping the grpc generator and switching to plain classes?

	// code will exec tests via a lib (loading in generic def)

	resp, err := svc.AddLocation(nil, &AddLocationReq{Name: "rand1"})

	if qa.assert(resp, err) {
		return
	}

	resp3, err3 := svc.AddProduct(nil, &AddProductReq{Name: "something"})

	if qa.assert(resp3, err3) {
		return
	}

	r5, err5 := svc.UpdateQty(nil, &UpdateQtyReq{
		Location: resp.Id,
		Product:  resp3.Id,
		Quantity: -1,
	})

	if qa.assert(r5, err5) {
		return
	}

	// TODO: structural compare via text to response

	r2, err2 := svc.GetInventory(nil, &GetInventoryReq{Location: resp.Id})

	if qa.assert(r2, err2) {
		return
	}

}
