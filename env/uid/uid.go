package uid

import (
	"fmt"
	"log"
	"strconv"
)

var TestMode = false

func ToTestString(v int64) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012x", v)
}

func Str(v int64) string {
	if TestMode {
		return fmt.Sprintf("00000000-0000-0000-0000-%012x", v)
	} else {
		return strconv.FormatInt(v, 16)
	}
}

func Parse(s string) int64 {
	if len(s) == 0 {
		panic("GUID can't be empty")
	}
	if TestMode {
		return ParseTestUuid(s)
	}
	i, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		log.Panicf("Can't parse id %q in prod mode", s)
	}
	return i

}

func ParseTestUuid(s string) int64 {
	if len(s) != 36 {
		log.Panicf("Id in test mode is not guid-formatted: %q", s)
	}
	val, err := strconv.ParseInt(s[24:], 16, 64)
	if err != nil {
		log.Panicf("Can't parse guid in test mode: %q", s)
	}
	return val
}
