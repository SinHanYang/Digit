package main

import (
	"time"
	"sort"
	//"fmt"
	//"Digit/libraries/core/diff"
)
//struct

type Commit struct {
	time time.Time
	author string
	hash string
	value ProllyTree // type Value when sync
	message string
}

type Head struct{
	// head record head commit's hash
	CommitHash string //hash encode->string
	// branch implementation
	// parent *[]
}

//global variable

var commit_graph map[string]Commit
var Head_pt Head

// functions

//init map
func init(){
	commit_graph=make(map[string]Commit)
}

// set head
func set_head(hash string){
	Head_pt.CommitHash=hash
}

// new commit
func new_commit(time time.Time, author string, hash string, value ProllyTree, message string){
	copyprollytree:=value
	node:=Commit{
		time: time,
		author: author,
		hash: hash,
		value: copyprollytree,
		message: message,
	}
	//fmt.Printf("finish")
	commit_graph[hash]=node
	set_head(hash)
}

// find head
func find_head () Commit{
	return commit_graph[Head_pt.CommitHash]
}

// find commit
func find_commit (hash string) Commit{
	return commit_graph[hash]
}

//reset
func reset(hash string){
	Head_pt.CommitHash=hash
}

// log , sort by time
func listcommit() []Commit{
	list:=make([]Commit,len(commit_graph))
	for _,val :=range commit_graph{
		list=append(list,val)
	}
	sort.Slice(list, func(i, j int) bool {return list[i].time.After(list[j].time)})
	return list
}