package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
	"sdk-go/inventory"
	. "sdk-go/protos"
	"strings"
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

	resp, err := s.GetInventory(context.Background(), &GetInventoryReq{Location: 1})
	var r DiffReporter
	if diff := cmp.Diff(spec.Expect, resp, cmp.Reporter(&r), cmpopts.IgnoreUnexported(GetInventoryResp{}, GetInventoryResp_Item{})); diff != "" {
		t.Error(r.String())
	}
}

// DiffReporter is a simple custom reporter that only records differences
// detected during comparison.
type DiffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		msg := fmt.Sprintf("Expected %#v to be '%+v' but got '%+v'\n", r.path[1:], vx, vy)
		r.diffs = append(r.diffs, msg)
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}
