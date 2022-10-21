package inventory

import (
	"black-friday/api"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func zeroToNill(n uint64) any {
	// because NULL is good in SQLite for rows that have FK
	// and not have a record to point to
	if n == 0 {
		return nil
	}
	return n
}

func apply(tx *Tx, e proto.Message) error {
	switch t := e.(type) {
	case *api.LocationAdded:
		values := []any{t.Id, t.Name, zeroToNill(t.Parent), t.Id, "Locations"}
		return tx.Exec(`
INSERT INTO Locations(Id, Name, Parent) VALUES (?,?,?);
UPDATE sqlite_sequence SET seq=? WHERE name=?
`, values...)
	case *api.LocationMoved:
		return tx.Exec(`
UPDATE Locations SET Parent=? WHERE Id=?
`, zeroToNill(t.NewParent), t.Id)
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

	default:
		return fmt.Errorf("Unhandled event: %s", e.ProtoReflect().Descriptor().Name())
	}
}
