package seq

import (
	"fmt"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"strings"
)

type Delta struct {
	Expected, Actual any
	Path             string
}

func Format(val any) string {
	if val == nil {
		return "<nil>"
	}
	switch v := val.(type) {
	case proto.Message:

		repr := prototext.MarshalOptions{Multiline: false}.Format(v)
		return string(v.ProtoReflect().Descriptor().Name()) + ":" + repr
	case []proto.Message:
		names := []string{}
		for _, m := range v {
			names = append(names, string(m.ProtoReflect().Descriptor().Name()))
		}
		return fmt.Sprintf("[%s]", strings.Join(names, ", "))

	case error:
		return fmt.Sprintf("Error '%v'", v.Error())
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func (d *Delta) String() string {
	return fmt.Sprintf("Expected %v to be %v but got %v",
		d.Path,
		Format(d.Expected),
		Format(d.Actual))
}

func Diff(expected, actual proto.Message, prefix string) (r []*Delta) {

	var path []string
	if prefix != "" {
		path = append(path, prefix)
	}

	enil, anil := expected == nil, actual == nil
	if enil && anil {
		// both are nil. Good
		return nil
	}

	if enil != anil {
		// one of them is nil. Quit now, too
		return []*Delta{
			{
				Expected: expected,

				Actual: actual,
				Path:   prefix,
			},
		}
	}

	return compare(expected.ProtoReflect(), actual.ProtoReflect(), path...)
}

func compare(expected, actual protoreflect.Message, path ...string) (r []*Delta) {
	e, a := expected, actual
	ed, ad := e.Descriptor(), a.Descriptor()
	if ed != ad {

		r = append(r, &Delta{
			Expected: string(e.Descriptor().Name()),
			Actual:   string(a.Descriptor().Name()),
			Path:     strings.Join(append(path, "type"), "."),
		})
		return r
	}

	for i := 0; i < ed.Fields().Len(); i++ {
		field := ed.Fields().Get(i)

		s := field.TextName()

		pth := strings.Join(append(path, s), ".")

		ev := e.Get(field)
		av := a.Get(field)

		switch {
		case field.IsList():
			r = append(r, handleList(field, ev, av, pth)...)
		case field.IsMap():
			panic("maps not handled")
		default:
			deltas := handleSingular(field, ev, av, pth)
			r = append(r, deltas...)
		}
	}

	return r
}

func handleList(field protoreflect.FieldDescriptor, ev protoreflect.Value, av protoreflect.Value, pth string) (r []*Delta) {
	el := ev.List()
	al := av.List()

	if el.Len() != al.Len() {
		return []*Delta{{Expected: el.Len(), Actual: al.Len(), Path: pth + ".length"}}
	} else {
		for i := 0; i < el.Len(); i++ {
			ev, av := el.Get(i), al.Get(i)
			deltas := handleSingular(field, ev, av, fmt.Sprintf("%s[%d]", pth, i))
			r = append(r, deltas...)
		}
	}
	return r
}

func handleSingular(field protoreflect.FieldDescriptor, ev protoreflect.Value, av protoreflect.Value, pth string) []*Delta {
	switch field.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return compare(ev.Message(), av.Message(), pth)
	default:
		if !reflect.DeepEqual(ev.Interface(), av.Interface()) {
			return []*Delta{{
				Expected: ev.Interface(),
				Actual:   av.Interface(),
				Path:     pth,
			}}
		}

	}
	return nil
}
