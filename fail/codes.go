package fail

type Code int

const (
	None          = Code(0)
	SqlLogicError = 1

	Constraint        = 19 // sqlite
	ConstraintUnique  = 20
	ConstraintForeign = 21

	Unknown = 42
)
