package tests

import (
	"context"
	"database/sql"
	"google.golang.org/protobuf/proto"
	"sdk-go/inventory"
	. "sdk-go/protos"
	"sdk-go/seq"
	"testing"
)

type Spec struct {
	Given  []proto.Message
	When   proto.Message
	Expect proto.Message
}

func Test_Spec(t *testing.T) {

	spec := &Spec{
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Cola"},
			&ProductAdded{Id: 2, Sku: "Fanta"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&QuantityUpdated{Location: 1, Product: 2, Quantity: 1, Total: 2, Before: 0},
		},
		When:   &GetInventoryReq{Location: 1},
		Expect: &GetInventoryResp{Items: []*GetInventoryResp_Item{{Product: 2, Quantity: -1}}},
	}

	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("sqlite3", ":memory:")
	check(err)
	defer db.Close()

	check(inventory.CreateSchema(db))

	s := inventory.NewService(db)

	check(s.Apply(spec.Given))

	actual, err := s.GetInventory(context.Background(), &GetInventoryReq{Location: 1})

	deltas := seq.Diff(spec.Expect, actual)

	for _, d := range deltas {
		t.Error(d.String())
	}
}
