package diff

import "fmt"

type Value interface {
	Less(Value) bool
	GetKey() string
	GetData() map[string]string
	GetHash() ChunkAddress
}

type ChunkAddress [20]byte

type Cursor interface {
	Done() bool
	Current() Value
	// update Current, Path
	Next()

	// path to specific chunk, start at root of tree
	Path() []ChunkAddress

	NextAtLevel(int)

	GetHash() ChunkAddress
	GetTree() ProllyTree
}

type OpType string

const (
	Add    OpType = "ADD"
	Delete        = "DELETE"
	Edit          = "EDIT"
)

type Difference struct {
	Op    OpType
	Value Value
}

func Diff(left, right Cursor) []Difference {
	var res []Difference
	counter := 0
	for !left.Done() && !right.Done() {
		lv, rv := left.Current(), right.Current()
		fmt.Println("Left:", lv.GetKey())
		fmt.Println("Right:", rv.GetKey())
		counter++
		if lv.Less(rv) {
			res = append(res, Difference{Op: Delete, Value: lv})
			left.Next()
		} else if rv.Less(lv) {
			res = append(res, Difference{Op: Add, Value: rv})
			right.Next()
		} else {
			if lv.GetHash() == rv.GetHash() {
				fmt.Println("FastForward!")
				FastForwardUntilUnequal(left, right)
			} else {
				res = append(res, Difference{Op: Edit, Value: rv})
				left.Next()
				right.Next()
			}
		}
	}

	for !left.Done() {
		res = append(res, Difference{Op: Delete, Value: left.Current()})
		left.Next()
		counter++
	}
	for !right.Done() {
		res = append(res, Difference{Op: Add, Value: right.Current()})
		right.Next()
		counter++
	}
	fmt.Println(counter)
	return res
}

// Advance until one of them is Done or unequal
func FastForwardUntilUnequal(left, right Cursor) {
	for !left.Done() && !right.Done() {
		lv, rv := left.Current(), right.Current()
		if lv.Less(rv) || rv.Less(lv) || lv.GetHash() != rv.GetHash() {
			return
		}
		level := GreatestMatchingLevelForPaths(left.Path(), right.Path())
		fmt.Println("At Level:", level+1)
		left.NextAtLevel(level + 1)
		right.NextAtLevel(level + 1)
		// fmt.Println("Left After Next At Level", level+1, ":", left.Current().GetKey())
		// fmt.Println("Right After Next At Level", level+1, ":", right.Current().GetKey())
	}
}

// Return the highest level in the tree at paths match
// Return -1 if there is no chunk address that matches
func GreatestMatchingLevelForPaths(left, right []ChunkAddress) int {
	level := -1
	for li, ri := len(left)-1, len(right)-1; li >= 0 && ri >= 0; li, ri, level = li-1, ri-1, level+1 {
		// fmt.Println(left[li])
		// fmt.Println(right[ri])
		if left[li] != right[ri] {
			break
		}
	}
	return level
}
