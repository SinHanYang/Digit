package diff

import (
	"Digit/libraries/core/postgres"
	"fmt"
)

type StageMap map[int]Difference

type Stage struct {
	StageAdded StageMap
}

func NewStage() Stage {
	return Stage{
		StageAdded: make(StageMap),
	}
}

// Print Stages
func (s *Stage) PrintStatus() {
	fmt.Println("============== ADDED ============== ")
	for _, v := range s.StageAdded {
		fmt.Println(v.Op, v.Value.GetData())
	}
}

// Update all stages
func (s *Stage) Add(head_cursor ChunkCursor, current_cursor ChunkCursor) []Difference {
	diff := Diff(&head_cursor, &current_cursor)
	for _, dif := range diff {
		s.StageAdded[dif.Value.GetKey()] = dif
	}
	return diff
}

// Unstage All -> Rollback all
func (s *Stage) Unstage(lastid int, view_name string, connstr string) {
	// TODO rollback
	var rollback [][2]int
	postgres.Rollback(lastid, view_name, connstr)
	for key, dif := range s.StageAdded {
		if dif.Op == Add {
			// do nothing ?
		} else if dif.Op == Delete || dif.Op == Edit {
			// postgres.ModifyStatus(view_name, , coconnstr)
			rollback = append(rollback, [2]int{int(key), 1})
		}
	}
	postgres.ModifyStatus(view_name, rollback, connstr)
	s.StageAdded = make(StageMap)
}

// Commit : Clear all staged
func (s *Stage) Commit() {
	// check StageAdded is not empty
	if len(s.StageAdded) == 0 {
		fmt.Println("There's nothing to commit!")
		return
	}

	// TODO real commit
	// current prolly tree to new commit
	s.StageAdded = make(StageMap)
}
