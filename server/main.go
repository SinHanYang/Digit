package main

import (
	"Digit/libraries/core/commit"
	"Digit/libraries/core/diff"
	"Digit/libraries/core/postgres"
	"Digit/libraries/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var cg commit.CommitGraph
var stage diff.Stage

type DigitRequest struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
	// for add
	Add_table []string
	Add_all   bool
	// for sql
	Sql_query string
}

type DigitResponse struct {
	Status  string
	Message string
	Data    interface{}
}

var resp DigitResponse

func initHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/init")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)
	var config DigitRequest
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

func addHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/add")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)
	var data DigitRequest
	err := json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if data.Add_all {
		// TODO: add all tables
	} else {
		// TODO: add tables
		// for _, table := range data.Add_table {
		// 	stage.AddTable(table)
		// }
	}
	w.Write([]byte("OK"))
}

func sqlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/sql")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)
	var data DigitRequest
	err := json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Println("Unable to parse request")
		return
	}
	var connstr = "postgres://" + data.DB_USER + ":" + data.DB_PASS + "@localhost:5432/" + data.DB_NAME
	db, err := sql.Open("pgx", connstr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Println("Connection Failed")
		return
	}
	log.Println("Create Connection")
	// check sql query
	// if select, use Query
	if strings.ToUpper(strings.Split(data.Sql_query, " ")[0]) == "SELECT" {
		rows, err := db.Query(data.Sql_query)
		if err != nil {
			resp.Status = "error"
			resp.Message = err.Error()
			resp_json, _ := json.Marshal(resp)
			w.Write(resp_json)
			log.Println(err)
			return
		}
		resp = DigitResponse{
			Status: "OK",
			Data:   utils.ReadRows(rows),
		}
		resp_json, _ := json.Marshal(resp)
		w.Write(resp_json)
	} else {
		// if insert, use Exec
		_, err := db.Exec(data.Sql_query)
		if err != nil {
			resp.Status = "error"
			resp.Message = err.Error()
			resp_json, _ := json.Marshal(resp)
			w.Write(resp_json)
			log.Println(err)
			return
		}
		resp = DigitResponse{
			Status: "OK",
			Data:   nil,
		}
		resp_json, _ := json.Marshal(resp)
		w.Write(resp_json)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Digit core server!")
	})
	http.HandleFunc("/init", initHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/sql", sqlHandler)
	log.Fatal(http.ListenAndServe(":17617", nil))
}

// func main() {
// 	t := time.Now().Format("20060102150405")
// 	db_user := "fiona"
// 	db_pass := "12345678"
// 	Init(t, db_user, db_pass)
// }
