package seq

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"strings"
)

type Delta struct {
	Expected, Actual any
	Path             string
}

func (d *Delta) String() string {
	return fmt.Sprintf("Expected %v to be %v but got %v", d.Path, d.Expected, d.Actual)
}

func Diff(expected, actual proto.Message) (r []*Delta) {
	return compare(expected.ProtoReflect(), actual.ProtoReflect())
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
			el := ev.List()
			al := av.List()

			if el.Len() != al.Len() {
				r = append(r, &Delta{
					Expected: int(el.Len()),
					Actual:   int(al.Len()),
					Path:     pth + ".length",
				})
			} else {
				for i := 0; i < el.Len(); i++ {
					ev, av := el.Get(i), al.Get(i)
					if deltas := handleSingular(field, ev, av, fmt.Sprintf("%s[%d]", pth, i)); len(deltas) > 0 {
						r = append(r, deltas...)
					}
				}
			}
		case field.IsMap():
			panic("mapos not handled")
		default:
			if deltas := handleSingular(field, ev, av, pth); len(deltas) > 0 {
				r = append(r, deltas...)
			}
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
