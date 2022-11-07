package rnd

import "hash/maphash"

func rand64() uint64 {
	return maphash.Bytes(maphash.MakeSeed(), nil)
}
