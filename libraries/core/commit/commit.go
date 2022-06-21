package commit

import (
	. "Digit/libraries/core/dataloader"
	"Digit/libraries/core/diff"
	"sort"
	"time"
	//"fmt"
)

//struct

type Commit struct {
	Time    time.Time
	Author  string
	Hash    string
	Value   diff.ProllyTree // type Value when sync
	Message string
}

type CommitGraph struct {
	commit_graph map[string]Commit
	head_hash    string
}

// functions

//init map
func NewCommitGraph() CommitGraph {
	return CommitGraph{
		commit_graph: make(map[string]Commit),
		head_hash:    "",
	}
}

// set head
func (cg *CommitGraph) SetHead(hash string) {
	cg.head_hash = hash
}

// new commit
func (cg *CommitGraph) NewCommit(time time.Time, author string, hash string, value diff.ProllyTree, message string) {
	copyprollytree := value

	copyprollytree = LoadTree(SaveTree(value))

	node := Commit{
		Time:    time,
		Author:  author,
		Hash:    hash,
		Value:   copyprollytree,
		Message: message,
	}
	//fmt.Printf("finish")
	cg.commit_graph[hash] = node
	cg.SetHead(hash)
}

// find head
func (cg *CommitGraph) GetHeadCommit() Commit {
	return cg.GetCommit(cg.head_hash)
}

// find commit
func (cg *CommitGraph) GetCommit(hash string) Commit {
	return cg.commit_graph[hash]
}

// log , sort by time
func (cg *CommitGraph) ListCommits() []Commit {
	list := make([]Commit, 0)
	for _, val := range cg.commit_graph {
		list = append(list, val)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Time.After(list[j].Time) })
	return list
}
