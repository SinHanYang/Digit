package main

import (
	"Digit/libraries/core/commit"
	"Digit/libraries/core/diff"
	"Digit/libraries/core/postgres"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var cg commit.CommitGraph
var stage diff.Stage

type Config struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/init")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)
	var config Config
	err := json.Unmarshal(body, &config)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	postgres.Init(config.DB_NAME, config.DB_USER, config.DB_PASS)
	cg = commit.NewCommitGraph()
	stage = diff.NewStage()
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Digit core server!")
	})
	http.HandleFunc("/init", initHandler)
	log.Fatal(http.ListenAndServe(":8089", nil))
}

// func main() {
// 	t := time.Now().Format("20060102150405")
// 	db_user := "fiona"
// 	db_pass := "12345678"
// 	Init(t, db_user, db_pass)
// }
