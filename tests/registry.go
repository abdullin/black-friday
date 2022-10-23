package tests

import (
	"black-friday/specs"
)

var Specs []*specs.S

func register(s *specs.S) {
	Specs = append(Specs, s)
}
