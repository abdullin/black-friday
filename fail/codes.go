package fail

type Code int

const (
	OK = Code(0)

	Constraint        = 19 // sqlite
	ConstraintUnique  = 20
	ConstraintForeign = 21

	Unknown = 42
)