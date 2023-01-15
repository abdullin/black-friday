package mem

import (
	"fmt"
	"strings"
)

// Stocks is a transactional in-memory model
// that tries to be ECS-friendly

// we apply changes in sequence as events. But if a change fails - we can always roll things back
type Stocks map[int32]*Stock

type Stock struct {
	lines []Line
}
type Line struct {
	loc int32
	// root has parent index of -1
	reserved  int32
	onHand    int32
	parentIdx int16
}

func NewStocks() Stocks {
	return make(Stocks)
}

func ToTestString(lines []Line) string {
	var b strings.Builder

	for _, l := range lines {
		b.WriteString(fmt.Sprintf("%2d %2d %2d %2d\n", l.loc, l.parentIdx, l.onHand, l.reserved))
	}
	return b.String()
}

func Ensure(lines []Line, branch []int32, qty int32, reserve int32) []Line {

	var bi, ti int16

	path := make([]int16, 0, MAX_DEPTH)

	// first of all, advance as deep in the locs, as possible
	// we always advance at least one step
	size := len(lines)
	for {
		//fmt.Println(fmt.Sprintf("bi %d ti %d\n", bi, ti))
		// check if we are still in the locs
		if branch[bi] == lines[ti].loc {

			path = append(path, ti)

			bi += 1

			// are we at the branch end?
			if bi >= int16(len(branch)) {
				break
			}
		}

		ti += 1

		// are we at the locs end?
		if ti >= int16(size) {
			break
		}
	}

	parentIdx := path[len(path)-1]

	// extend
	extender := branch[bi:]

	extendSize := len(extender)

	if extendSize > 0 {
		// Make sure there is space to append n elements without re-allocating:
		lines = append(make([]Line, 0, len(lines)+len(extender)), lines...)

		for _, b := range extender {

			l := Line{
				parentIdx: parentIdx,
				reserved:  reserve,
				loc:       b,
				onHand:    qty,
			}
			lines = append(lines, l)

			parentIdx = int16(len(lines) - 1)
		}
	}

	// now increment the quantity

	for _, i := range path {
		lines[i].onHand += qty
		lines[i].reserved += reserve
	}
	return lines

}

// Update loc hierarchy in forward
// pass (chance to reuse cache lines)
func Update(lines []Line, loc int32, qty int32, reserve int32) (outQty, outReserve int32) {
	count := int16(len(lines))
	for i := count - 1; i >= 0; i-- {
		if lines[i].loc != loc {
			continue
		}
		outQty = lines[i].onHand + qty
		outReserve = lines[i].reserved + reserve
		lines[i].onHand = outQty
		lines[i].reserved = outReserve

		i = lines[i].parentIdx
	}
	return
}
