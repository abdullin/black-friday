package fail

import (
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

// Extract error details into a clean internal enum
// so that we can pattern match in API methods
func Extract(err error) (error, Code) {

	if err == nil {
		return nil, OK
	}

	var sqlErr sqlite3.Error

	if errors.As(err, &sqlErr) {
		switch sqlErr.ExtendedCode {
		case sqlite3.ErrConstraintForeignKey:
			return err, ConstraintForeign
		case sqlite3.ErrConstraintUnique:
			return err, ConstraintUnique
		default:
			panic(fmt.Errorf("Unexpected sql extended code %v", sqlErr.ExtendedCode))
		}
	}

	return err, Unknown

}
