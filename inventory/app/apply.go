package app

import (
	"black-friday/fail"
	"black-friday/inventory/api"
	"black-friday/inventory/db"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func (c *Context) Apply(e proto.Message) (error, fail.Code) {

	err := applyInner(c, e)

	if err != nil {
		extracted, failCode := fail.Extract(err)
		return fmt.Errorf("apply %s: %w", reflect.TypeOf(e).String(), extracted), failCode
	}

	c.events = append(c.events, e)
	return nil, fail.None

}

func (c *Context) TestClear() {
	c.events = nil
}

func (c *Context) TestGet() []proto.Message {
	return c.events
}

func applyInner(tx *Context, e proto.Message) error {
	switch t := e.(type) {
	case *api.LocationAdded:
		values := []any{t.Id, t.Name, db.ZeroToNil(t.Parent), t.Id, "Locations"}
		return tx.Exec(`
INSERT INTO Locations(Id, Name, Parent) VALUES (?,?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, values...)
	case *api.LocationMoved:
		return tx.Exec(`
UPDATE Locations SET Parent=? WHERE Id=?
`, db.ZeroToNil(t.NewParent), t.Id)
	case *api.ProductAdded:
		return tx.Exec(`
INSERT INTO Products(Id, Sku) VALUES (?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, t.Id, t.Sku, t.Id, "Products")
	case *api.InventoryUpdated:

		before := t.OnHand - t.OnHandChange
		if t.OnHand == 0 {
			return tx.Exec("DELETE FROM Inventory WHERE Product=? AND Location=?", t.Product, t.Location)
		} else if before == 0 {
			return tx.Exec("INSERT INTO Inventory(Product, Location, OnHand) VALUES(?,?,?)", t.Product, t.Location, t.OnHand)
		} else {
			return tx.Exec("UPDATE Inventory SET OnHand=? WHERE Product=? AND Location=?", t.OnHand, t.Product, t.Location)
		}
	case *api.Reserved:
		err := tx.Exec(`
INSERT INTO Reservations(Id, Code) VALUES(?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, t.Reservation, t.Code, t.Reservation, "Reservations")
		if err == nil {
			return err
		}
		for _, i := range t.Items {
			err = tx.Exec("INSERT INTO Reserves (Reservation, Product, Quantity, Location) VALUES(?,?,?,?)",
				t.Reservation, i.Product, i.Quantity, i.Location,
			)
			if err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("Unhandled event: %s", e.ProtoReflect().Descriptor().Name())
	}
}
