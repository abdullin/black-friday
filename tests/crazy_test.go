package tests

import "testing"

type Command struct{}
type Message interface{}

func command() (m1 *Command, m2 *Command, m3 Message) {
	return nil, nil, nil
}
func dispatch() (m1 *Command, m2 Message, m3 Message) {
	return command()
}

func TestHalfNil(t *testing.T) {
	m1, m2, m3 := dispatch()
	t.Log(m1 == nil)
	t.Log(m2 == nil)
	t.Log(m3 == nil)
}
