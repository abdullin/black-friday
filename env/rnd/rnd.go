package rnd

import (
	"math"
	"math/bits"
)

type Rand struct {
	a uint64
	b uint64
	c uint64
	w uint64
}

func (s *Rand) init(a uint64, b uint64, c uint64) {
	s.a = a
	s.b = b
	s.c = c
	s.w = 1
	for i := 0; i < 12; i++ {
		s.next64()
	}
}

func (s *Rand) init0() {
	s.a = rand64()
	s.b = rand64()
	s.c = rand64()
	s.w = 1
}

func (s *Rand) init1(u uint64) {
	s.a = u
	s.b = u
	s.c = u
	s.w = 1
	for i := 0; i < 12; i++ {
		s.next64()
	}
}

func (s *Rand) init3(a uint64, b uint64, c uint64) {
	s.a = a
	s.b = b
	s.c = c
	s.w = 1
	for i := 0; i < 18; i++ {
		s.next64()
	}
}

func (s *Rand) next64() (out uint64) { // named return value lowers inlining cost
	out = s.a + s.b + s.w
	s.w++
	s.a, s.b, s.c = s.b^(s.b>>11), s.c+(s.c<<3), bits.RotateLeft64(s.c, 24)+out // single assignment lowers inlining cost
	return
}

func New() *Rand {

	s := &Rand{}
	s.init0()
	return s
}

func (r *Rand) Uint32n(n uint32) uint32 {
	// much faster 32-bit version of Uint64n(); result is unbiased with probability 1 - 2^-32.
	// detecting possible bias would require at least 2^64 samples, which we consider acceptable
	// since it matches 2^64 guarantees about period length and distance between different seeds.
	// note that 2^64 is probably a very conservative estimate: scaled down 16-bit version of this
	// algorithm passes chi-squared test for at least 2^42 (instead of 2^32) values, so
	// 32-bit version will likely require north of 2^80 values to detect non-uniformity.
	res, _ := bits.Mul64(uint64(n), r.next64())
	return uint32(res)
}

// Int63n returns, as an int64, a uniformly distributed non-negative pseudo-random number
// in the half-open interval [0, n). It panics if n <= 0.
func (r *Rand) Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	return int64(r.Uint64n(uint64(n)))
}
func (r *Rand) Uint64n(n uint64) uint64 {
	// "An optimal algorithm for bounded random integers" by Stephen Canon, https://github.com/apple/swift/pull/39143
	res, frac := bits.Mul64(n, r.next64())
	if n <= math.MaxUint32 {
		// we don't use frac <= -n check from the original algorithm, since the branch is unpredictable.
		// instead, we effectively fall back to Uint32n() for 32-bit n
		return res
	}
	hi, _ := bits.Mul64(n, r.next64())
	_, carry := bits.Add64(frac, hi, 0)
	return res + carry
}
