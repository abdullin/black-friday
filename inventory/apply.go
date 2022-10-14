package inventory

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"sdk-go/protos"
)

func apply(tx *Tx, e proto.Message) error {
	lift := func(err error) error {
		if err != nil {
			name := e.ProtoReflect().Descriptor().Name()
			return fmt.Errorf("problem applying '%s': %w", name, err)
		}
		return nil
	}

	switch t := e.(type) {
	case *protos.WarehouseCreated:
		return lift(tx.Exec(`
INSERT INTO Warehouses(Id, Name) VALUES(?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?`,
			t.Id, t.Name, t.Id, "Warehouses"))
	case *protos.LocationAdded:
		return lift(tx.Exec(`
INSERT INTO Locations(Id, Name, Warehouse) VALUES (?,?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, t.Id, t.Name, t.Warehouse, t.Id, "Locations"))
	case *protos.ProductAdded:
		return lift(tx.Exec(`
INSERT INTO Products(Id, Sku) VALUES (?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, t.Id, t.Sku, t.Id, "Products"))
	case *protos.InventoryUpdated:

		before := t.OnHand - t.OnHandChange
		if t.OnHand == 0 {
			return lift(tx.Exec("DELETE FROM Inventory WHERE Product=? AND Location=?", t.Product, t.Location))
		} else if before == 0 {
			return lift(tx.Exec("INSERT INTO Inventory(Product, Location, OnHand) VALUES(?,?,?)", t.Product, t.Location, t.OnHand))
		} else {
			return lift(tx.Exec("UPDATE Inventory SET OnHand=? WHERE Product=? AND Location=?", t.OnHand, t.Product, t.Location))
		}

	default:
		return fmt.Errorf("Unhandled event: %s", e.ProtoReflect().Descriptor().Name())
	}
}
