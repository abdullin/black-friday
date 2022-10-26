package specs

var Specs []*S

func Add(s *S) {
	Specs = append(Specs, s)
}
