package specs

import (
	"black-friday/inventory/api"
	"bufio"
	"fmt"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"log"
	"strings"
)

func msgToString(m proto.Message) string {
	body := prototext.MarshalOptions{
		Multiline: false,
	}.Format(m)

	return fmt.Sprintf("%s %s", m.ProtoReflect().Descriptor().Name(), body)
}

func stringsToMsgs(lines []string) ([]proto.Message, error) {
	var msgs []proto.Message
	for _, l := range lines {
		m, err := stringToMsg(l)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func stringToMsg(s string) (proto.Message, error) {

	source := strings.TrimSpace(s)
	splits := strings.SplitN(source, " ", 2)

	name := splits[0]

	mt, found := lookups[name]
	if !found {
		return nil, fmt.Errorf("Unknown type: %s", name)
	}

	instance := mt.New().Interface()

	if len(splits) == 2 {
		body := []byte(splits[1])

		if err := prototext.Unmarshal(body, instance); err != nil {
			return nil, err
		}
	}

	return instance, nil

}

func SpecToParseableString(s *api.Spec) (string, error) {

	var b strings.Builder
	ln := func(text string, args ...any) {
		_, err := b.WriteString(fmt.Sprintf(text, args...) + "\n")
		if err != nil {
			panic(err)
		}
	}

	ln(s.Name)
	ln("GIVEN:")

	for _, e := range s.Given {
		ln("  %s", msgToString(e))
	}
	ln("WHEN:\n%s", msgToString(s.When))
	if s.ThenResponse != nil {
		ln("THEN:\n%s", msgToString(s.ThenResponse))
	}

	if len(s.ThenEvents) > 0 {
		ln("EVENTS:")
		for _, e := range s.ThenEvents {
			ln("  %s", msgToString(e))
		}
	}
	if s.ThenError != nil {
		ln("ERROR: %s", s.ThenError.Error())
	}

	return b.String(), nil

}

var lookups = make(map[string]protoreflect.MessageType)

func init() {
	msgs := api.File_inventory_api_api_proto.Messages()

	for i := 0; i < msgs.Len(); i++ {
		m := msgs.Get(i)

		mt, err := protoregistry.GlobalTypes.FindMessageByName(m.FullName())
		if err != nil {
			log.Panicln(err)
		}

		lookups[string(m.Name())] = mt
	}
}

func groupLines(lines []string) (map[string][]string, error) {
	current := ""
	result := make(map[string][]string)
	for i, l := range lines {
		trimmed := strings.TrimSpace(l)
		if len(trimmed) == 0 {
			continue
		}
		switch trimmed {
		case GIVEN, WHEN, THEN, ERROR, EVENTS:
			current = trimmed
			continue
		default:
			if len(current) == 0 {
				return nil, fmt.Errorf("Expected suffix on line %d, got: %s", i, l)
			}
			array, _ := result[current]
			array = append(array, trimmed)
			result[current] = array
		}

	}
	return result, nil
}

const (
	GIVEN  = "GIVEN:"
	WHEN   = "WHEN:"
	THEN   = "THEN:"
	ERROR  = "ERROR:"
	EVENTS = "EVENTS:"
)

func SpecFromParseableString(s string) (*api.Spec, error) {

	scanner := bufio.NewScanner(strings.NewReader(s))

	var lines []string
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) > 0 {
			lines = append(lines, scanner.Text())
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	spec := &api.Spec{}

	spec.Name = lines[0]

	groups, err := groupLines(lines[1:])
	if err != nil {
		return nil, err
	}

	if lines, ok := groups[GIVEN]; ok {
		spec.Given, err = stringsToMsgs(lines)
		if err != nil {
			return nil, err
		}
	}
	if lines, found := groups[EVENTS]; found {
		spec.ThenEvents, err = stringsToMsgs(lines)
		if err != nil {
			return nil, err
		}
	}

	if lines, found := groups[WHEN]; found {
		if len(lines) > 1 {
			return nil, fmt.Errorf("must be only one message in %s. Got %v", WHEN, lines)
		}
		spec.When, err = stringToMsg(lines[0])
		if err != nil {
			return nil, err
		}
	}
	if lines, found := groups[THEN]; found {
		if len(lines) > 1 {
			return nil, fmt.Errorf("must be only one message in %s. Got %d", THEN, len(lines))
		}
		spec.ThenResponse, err = stringToMsg(lines[0])
		if err != nil {
			return nil, err
		}
	}
	if lines, found := groups[ERROR]; found {
		if len(lines) > 1 {
			return nil, fmt.Errorf("must be only one message in %s. Got %v", ERROR, lines)
		}
	}

	return spec, nil
}