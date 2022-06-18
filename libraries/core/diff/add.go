package diff

type StageList []Difference

type Stage struct {
	StageAdded StageList
	// StageModified  StageList
	StageUntracked StageList
}

func NewStage() Stage {
	return Stage{
		StageAdded: StageList{},
		// StageModified:  StageList{},
		StageUntracked: StageList{},
	}
}

//
func (s *Stage) Status(head_cursor ChunkCursor, current_cursor ChunkCursor) {
	if len(s.StageAdded) == 0 && len(s.StageUntracked) == 0 {
		diff := Diff(&head_cursor, &current_cursor)
		s.StageUntracked = diff
	}
}

// Stage All
func (s *Stage) Add(head_cursor ChunkCursor, current_cursor ChunkCursor) {
	s.Status(head_cursor, current_cursor)
	s.StageAdded = s.StageUntracked
	s.StageUntracked = StageList{}
}

// Unstage All
func (s *Stage) Unstage() {
	s.StageUntracked = s.StageAdded
	s.StageAdded = StageList{}
}

func (s *Stage) Commit() {
	s.StageAdded = StageList{}
	s.StageUntracked = StageList{}
}
