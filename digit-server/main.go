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
	"strconv"
	"strings"
	"time"
)

var cg commit.CommitGraph
var stage diff.Stage
var cursor diff.ChunkCursor
var headers = []string{"digitsnid", "digitstatus"}

type DigitRequest struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
	// for add
	Add_table []string
	Add_all   bool
	// for sql
	Sql_query string
	// for commit
	Commit_message string
	// for diff
	Diff_from string
	Diff_to   string
	// for reset
	Reset_hash string
	Reset_tb   string
}

type DigitResponse struct {
	Status  string
	Message string
	Data    interface{}
}

var resp DigitResponse
var init_flag bool = false

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
	var connstr = "postgres://" + data.DB_USER + ":" + data.DB_PASS + "@localhost:5432/" + data.DB_NAME
	// db, err := sql.Open("pgx", connstr)
	// if err != nil {
	// 	http.Error(w, "Bad request", http.StatusBadRequest)
	// 	fmt.Println("Connection Failed")
	// 	return
	// }

	if data.Add_all {
		// TODO: add all tables
		// get all tables
		// rows, _ := db.Query("select table_name from INFORMATION_SCHEMA.views WHERE table_schema = ANY (current_schemas(false));")
		// tables := utils.ReadRows(rows)
		// for _, table := range tables {

	} else {
		// TODO: add tables
		for _, table := range data.Add_table {
			if !init_flag {
				status := postgres.GetStatus(table, connstr)
				cursor = diff.NewCursor(headers, status)
				cg.NewCommit(time.Now(), data.DB_USER, diff.Encode(cursor.GetHash()), cursor.GetTree(), "Initial Commit")
				init_flag = true
			} else {
				status := postgres.GetStatus(table, connstr)
				cursor = diff.NewCursor(headers, status)
				stage.Add(diff.NewCursorFromProllyTree(cg.GetHeadCommit().Value), cursor)
				stage.PrintStatus()
			}
		}
	}
	w.Write([]byte("OK"))
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/commit")
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
	t := time.Now()
	cg.NewCommit(t, data.DB_USER, diff.Encode(cursor.GetHash()), cursor.GetTree(), data.Commit_message)
	stage.Commit()
	resp.Status = "ok"
	resp.Message = "Commit " + diff.Encode(cursor.GetHash()) + "\n" + "Author: " + data.DB_USER + "\n" + "Date: " + t.Format("2006-01-02 15:04:05") + "\n" + "Message: " + data.Commit_message
	resp_json, _ := json.Marshal(resp)
	w.Write(resp_json)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/log")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	logs := cg.ListCommits()
	resp.Status = "ok"
	resp.Message = ""
	for _, log := range logs {
		resp.Message += "Commit " + log.Hash + "\n" + "Author: " + log.Author + "\n" + "Date: " + log.Time.Format("2006-01-02 15:04:05") + "\n\n\t" + log.Message + "\n\n"
	}
	resp_json, _ := json.Marshal(resp)
	w.Write(resp_json)
}

func diffHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/diff")
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
	cursor_left := diff.NewCursorFromProllyTree(cg.GetCommit(data.Diff_from).Value)
	cursor_right := diff.NewCursorFromProllyTree(cg.GetCommit(data.Diff_to).Value)

	d := diff.Diff(&cursor_left, &cursor_right)

	resp.Status = "ok"
	resp.Message = ""

	var add_count int = 0
	var del_count int = 0
	var edit_count int = 0

	for _, v := range d {

		if v.Op == "ADD" {
			// resp.Message += "[ADD] " + fmt.Sprint(v.Value.GetData()) + "\n"
			add_count++
		} else if v.Op == "DELETE" {
			// resp.Message += "[DELETE] " + fmt.Sprint(v.Value.GetData()) + "\n"
			del_count++
		} else if v.Op == "EDIT" {
			// resp.Message += "[EDIT] " + fmt.Sprint(v.Value.GetData()) + "\n"
			edit_count++
		}
	}
	resp.Message += strconv.Itoa(add_count) + " row(s) added. "
	resp.Message += strconv.Itoa(del_count) + " row(s) deleted. "
	resp.Message += strconv.Itoa(edit_count) + " row(s) edited."
	resp_json, _ := json.Marshal(resp)
	w.Write(resp_json)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/reset")
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
	cg.SetHead(data.Reset_hash)
	connstr := "postgres://" + data.DB_USER + ":" + data.DB_PASS + "@localhost:5432/" + data.DB_NAME
	stage.Unstage(cg.GetHeadCommit().Value.Lastid, data.Reset_tb, connstr)
	resp.Status = "ok"
	resp.Message = "Reset to " + data.Reset_hash
	resp_json, _ := json.Marshal(resp)
	w.Write(resp_json)
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
	cg = commit.NewCommitGraph()
	stage = diff.NewStage()
	_ = stage
	emp := [][2]int{{}}
	cursor = diff.NewCursor(headers, emp)
	cg.NewCommit(time.Now(), "Digit", diff.Encode(cursor.GetHash()), cursor.GetTree(), "Initial Commit")
	// fmt.Println(cg.GetHeadCommit().Value.Lastid)
	init_flag = true
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Digit core server!")
	})
	http.HandleFunc("/init", initHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/commit", commitHandler)
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/diff", diffHandler)
	http.HandleFunc("/reset", resetHandler)
	http.HandleFunc("/sql", sqlHandler)
	log.Fatal(http.ListenAndServe(":17617", nil))
}

// func main() {
// 	t := time.Now().Format("20060102150405")
// 	db_user := "fiona"
// 	db_pass := "12345678"
// 	Init(t, db_user, db_pass)
// }
