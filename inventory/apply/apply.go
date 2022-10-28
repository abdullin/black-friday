package apply

import (
	"black-friday/fx"
	"black-friday/inventory/api"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func Event(tx fx.Tx, e proto.Message) error {
	switch t := e.(type) {
	case *api.LocationAdded:
		values := []any{t.Id, t.Name, t.Parent, t.Id, "Locations"}
		return tx.Exec(`
INSERT INTO Locations(Id, Name, Parent) VALUES (?,?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, values...)
	case *api.LocationMoved:
		return tx.Exec(`
UPDATE Locations SET Parent=? WHERE Id=?
`, t.NewParent, t.Id)
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