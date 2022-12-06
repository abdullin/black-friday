package apply

import (
	"black-friday/env/uid"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func Event(tx fx.Tx, e proto.Message) error {
	switch t := e.(type) {
	case *LocationAdded:

		id := uid.Parse(t.Uid)

		values := []any{id, t.Name, uid.Parse(t.Parent), id, "Locations"}
		return tx.Exec(`
INSERT INTO Locations(Id, Name, Parent) VALUES (?,?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, values...)
	case *LocationMoved:

		return tx.Exec(`
UPDATE Locations SET Parent=? WHERE Id=?
`, uid.Parse(t.NewParent), uid.Parse(t.Uid))
	case *ProductAdded:
		id := uid.Parse(t.Uid)
		return tx.Exec(`
INSERT INTO Products(Id, Sku) VALUES (?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, id, t.Sku, id, "Products")
	case *InventoryUpdated:

		before := t.OnHand - t.OnHandChange
		if t.OnHand == 0 {
			return tx.Exec("DELETE FROM Inventory WHERE Product=? AND Location=?", uid.Parse(t.Product), uid.Parse(t.Location))
		} else if before == 0 {
			return tx.Exec("INSERT INTO Inventory(Product, Location, OnHand) VALUES(?,?,?)", uid.Parse(t.Product), uid.Parse(t.Location), t.OnHand)
		} else {
			return tx.Exec("UPDATE Inventory SET OnHand=? WHERE Product=? AND Location=?", t.OnHand, uid.Parse(t.Product), uid.Parse(t.Location))
		}
	case *Reserved:

		id := uid.Parse(t.Reservation)
		err := tx.Exec(`
INSERT INTO Reservations(Id, Code) VALUES(?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, id, t.Code, id, "Reservations")
		if err != nil {
			return err
		}
		for _, i := range t.Items {
			err = tx.Exec("INSERT INTO Reserves (Reservation, Product, Quantity, Location) VALUES(?,?,?,?)",
				id, uid.Parse(i.Product), i.Quantity, uid.Parse(i.Location),
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
