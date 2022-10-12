package tests

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"sdk-go/inventory"
	. "sdk-go/protos"
	"sdk-go/seq"
	"testing"
)

type Spec struct {
	Name         string
	Given        []proto.Message
	When         proto.Message
	ThenResponse proto.Message
	ThenError    codes.Code
	ThenEvents   []proto.Message
}

func common_spec() *Spec {
	return &Spec{
		Name: "query inventory",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Cola"},
			&ProductAdded{Id: 2, Sku: "Fanta"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&QuantityUpdated{Location: 1, Product: 2, Quantity: 1, Total: 2, Before: 0},
		},
		When:         &GetInventoryReq{Location: 1},
		ThenResponse: &GetInventoryResp{Items: []*GetInventoryResp_Item{{Product: 2, Quantity: 2}}},
	}
}

func add_locations() *Spec {
	return &Spec{
		Name: "query locations",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf1"},
			&LocationAdded{Id: 2, Name: "Shelf2"},
		},
		When: &ListLocationsReq{},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Id: 1, Name: "Shelf1"},
			{Id: 2, Name: "Shelf2"},
		}},
	}
}

func run_spec(t *testing.T, spec *Spec) {
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

	check(s.ApplyEvents(spec.Given))

	actual, err := s.Dispatch(context.Background(), spec.When)

	deltas1 := seq.Diff(spec.ThenResponse, actual, "response")

	actualStatus, _ := status.FromError(err)

	if spec.ThenError != actualStatus.Code() {
		deltas1 = append(deltas1, &seq.Delta{
			Expected: spec.ThenError,
			Actual:   actualStatus.Code(),
			Path:     "status",
		})
	}

	for _, d := range deltas1 {
		t.Error(d.String())
	}
}

func Test_Spec(t *testing.T) {

	specs := []*Spec{add_locations(), common_spec()}

	for _, s := range specs {
		t.Run(s.Name, func(t *testing.T) {
			run_spec(t, s)
		})
	}

}
