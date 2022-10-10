package inventory

import (
	"database/sql"
	"google.golang.org/protobuf/proto"
	"log"
	"sdk-go/protos"
)

func must(r sql.Result, err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func Apply(tx *sql.Tx, e proto.Message) {

	switch t := e.(type) {
	case *protos.LocationAdded:

		must(tx.Exec("INSERT INTO Locations(Id, Name) VALUES (?,?)", t.Id, t.Name))
		must(tx.Exec("UPDATE sqlite_sequence SET seq=? WHERE name=?", t.Id, "Locations"))

	case *protos.ProductAdded:

		must(tx.Exec("INSERT INTO Products(Id, Sku) VALUES (?,?)", t.Id, t.Sku))
		must(tx.Exec("UPDATE sqlite_sequence SET seq=? WHERE name=?", t.Id, "Products"))
	case *protos.QuantityUpdated:
		if t.Total == 0 {
			must(tx.Exec("DELETE FROM Inventory WHERE Product=? AND Location=?", t.Product, t.Location))
		} else if t.Before == 0 {
			must(tx.Exec("INSERT INTO Inventory(Product, Location, Quantity) VALUES(?,?,?)", t.Product, t.Location, t.Total))
		} else {
			must(tx.Exec("UPDATE Inventory SET Quantity=? WHERE Product=? AND Location=?", t.Total, t.Product, t.Location))
		}

	default:
		panic("UNKNOWN EVENT")

	}
}
