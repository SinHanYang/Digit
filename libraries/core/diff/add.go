package main

import "fmt"

type StageMap map[string]Difference

type Stage struct {
	StageAdded     StageMap
	StageModified  StageMap
	StageUntracked StageMap
}

func NewStage() Stage {
	return Stage{
		StageAdded:     make(StageMap),
		StageModified:  make(StageMap),
		StageUntracked: make(StageMap),
	}
}

// Print Stages
func (s *Stage) PrintStatus() {
	fmt.Println("============== ADDED ============== ")
	for _, v := range s.StageAdded {
		fmt.Println(v.Op, v.Value.GetData())
	}
	fmt.Println("============== MODIFIED ============== ")
	for _, v := range s.StageModified {
		fmt.Println(v.Op, v.Value.GetData())
	}
	fmt.Println("============== UNTRACKED ==============")
	for _, v := range s.StageUntracked {
		fmt.Println(v.Op, v.Value.GetData())
	}
}

// Update all stages
func (s *Stage) Status(head_cursor ChunkCursor, current_cursor ChunkCursor) {
	diff := Diff(&head_cursor, &current_cursor)
	for _, dif := range diff {
		if val, ok := s.StageAdded[dif.Value.GetKey()]; ok {
			if dif.Value.GetHash() != val.Value.GetHash() {
				// move to modified
				s.StageModified[dif.Value.GetKey()] = dif
				delete(s.StageAdded, dif.Value.GetKey())
			}
			continue
		}
		if val, ok := s.StageModified[dif.Value.GetKey()]; ok {
			if dif.Value.GetHash() != val.Value.GetHash() {
				// replace
				s.StageModified[dif.Value.GetKey()] = dif
			}
			continue
		}
		// untrack still untrack
		s.StageUntracked[dif.Value.GetKey()] = dif
	}
}

// Stage All
func (s *Stage) Add(head_cursor ChunkCursor, current_cursor ChunkCursor) {
	s.Status(head_cursor, current_cursor)
	for key, val := range s.StageModified {
		s.StageAdded[key] = val
	}
	for key, val := range s.StageUntracked {
		s.StageAdded[key] = val
	}
	s.StageModified = make(StageMap)
	s.StageUntracked = make(StageMap)
}

// Unstage All
func (s *Stage) Unstage(head_cursor ChunkCursor, current_cursor ChunkCursor) {
	s.Status(head_cursor, current_cursor)
	for key, val := range s.StageAdded {
		s.StageModified[key] = val
	}
	s.StageAdded = make(StageMap)
}

// Commit : Clear all staged
func (s *Stage) Commit() {
	// TODO real commit
	s.StageAdded = make(StageMap)
	s.StageModified = make(StageMap)
	s.StageUntracked = make(StageMap)
}
