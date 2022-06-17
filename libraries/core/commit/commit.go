package main

import "time"
//struct

type Hash struct{
	hash_value [20]byte
	// Need to discuss hash struct details
}

type Commit struct {
	//prev *Node
	//next *Node
	time time.Time
	author string
	hash Hash
	value string// type Value when sync
}

type Head struct{
	// head record head commit's hash
	head Hash
	// branch implementation
	// parent *[]
}

//global variable

var commit_graph map[Hash]Commit
var Head_pt Head

// functions

//init map
func init(){
	commit_graph=make(map[Hash]Commit)
}

// set head
func set_head(hash Hash){
	Head_pt.head=hash
}
// new commit
func new_commit(time time.Time, author string, hash Hash, value string){
	node:=Commit{
		time: time,
		author: author,
		hash: hash,
		value: value,
	}
	//fmt.Printf("finish")
	commit_graph[hash]=node
	set_head(hash)
}

// find head
func find_head () Commit{
	return commit_graph[Head_pt.head]
}
// find commit
func find_commit (hash Hash) Commit{
	return commit_graph[hash]
}
//reset
func reset(hash Hash){
	Head_pt.head=hash
}
