package apply

import (
	"black-friday/env/uid"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"black-friday/inventory/features/graphs"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func setInventory(tx fx.Tx, product, location, onHand, delta int64) error {

	before := onHand - delta
	if onHand == 0 {
		return tx.Exec("DELETE FROM Inventory WHERE Product=? AND Location=?", product, location)
	} else if before == 0 {
		return tx.Exec("INSERT INTO Inventory(Product, Location, OnHand) VALUES(?,?,?)", product, location, onHand)
	} else {
		return tx.Exec("UPDATE Inventory SET OnHand=? WHERE Product=? AND Location=?", onHand, product, location)
	}
}

func Event(tx fx.Tx, e proto.Message) error {
	switch t := e.(type) {
	case *LocationAdded:

		id := uid.Parse(t.Uid)

		graphs.Cache = nil
		values := []any{id, t.Name, uid.Parse(t.Parent), id, "Locations"}
		return tx.Exec(`
INSERT INTO Locations(Id, Name, Parent) VALUES (?,?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, values...)
	case *LocationMoved:

		graphs.Cache = nil
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
		return setInventory(tx, uid.Parse(t.Product), uid.Parse(t.Location), t.OnHand, t.OnHandChange)
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
	case *Cancelled:
		rid := uid.Parse(t.Reservation)
		return tx.Exec(`
DELETE FROM Reserves WHERE Reservation=?; 
DELETE FROM Reservations WHERE Id=?;
`, rid, rid)
	case *Fulfilled:
		rid := uid.Parse(t.Reservation)

		for _, i := range t.Items {
			err := setInventory(tx,
				uid.Parse(i.Product),
				uid.Parse(i.Location),
				i.OnHand,
				-i.Removed)
			if err != nil {
				return nil
			}
		}

		return tx.Exec(`
DELETE FROM Reserves WHERE Reservation=?; 
DELETE FROM Reservations WHERE Id=?;
`, rid, rid)
	default:
		return fmt.Errorf("Unhandled event: %s", e.ProtoReflect().Descriptor().Name())
	}
}
