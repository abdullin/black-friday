package graphs

// map - loc/id -> leaf node -> chain up to root.

type Tree map[int32][]int32

func DenormaliseLocTree(in map[int32]int32) Tree {
	t := make(Tree)

	for id, parent := range in {

		var parents []int32

		for {
			parents = append(parents, parent)
			if parent == 0 {
				break
			} else {
				parent = in[parent]
			}

		}
		t[id] = parents

	}
	return t

}
