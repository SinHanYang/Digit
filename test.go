package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	. "Digit/libraries/core/commit"
	. "Digit/libraries/core/diff"
	. "Digit/libraries/core/postgres"

	_ "github.com/go-sql-driver/mysql"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Test() {

	Init()
	cg := NewCommitGraph()
	stage := NewStage()
	_ = stage

	db_name := "KoreanShowbiz"
	tb_name := "Member"
	username := "noidname"
	password := "test"
	var connstr = "postgres://" + username + ":" + password + "@localhost:5432/" + db_name
	db, err := sql.Open("pgx", connstr)
	check(err)
	status := GetStatus(tb_name, connstr)
	fmt.Println(status)

	headers := []string{"digitsnid", "digitstatus"}

	cursor_orig := NewCursor(headers, status)
	// initial commit
	cg.NewCommit(time.Now(), "Tim", Encode(cursor_orig.GetHash()), cursor_orig.GetTree(), "Initial Commit")
	fmt.Println(cg.GetHeadCommit().Message)

	// Insert Statement
	stmt, err := db.PrepareContext(context.Background(), "INSERT INTO "+tb_name+" VALUES ($1, $2, $3, $4, $5, $6)")
	check(err)
	_, err = stmt.ExecContext(context.Background(), "1", "2", "3", "4", "5", "6")
	check(err)
	// Update Statement
	stmt, err = db.PrepareContext(context.Background(), "UPDATE "+tb_name+" SET GroupName=$1 WHERE GroupName=$2")
	check(err)
	_, err = stmt.ExecContext(context.Background(), "OFF", "ONF")
	check(err)
	// Delete Statement
	stmt, err = db.PrepareContext(context.Background(), "DELETE FROM "+tb_name+" WHERE GroupName=$1")
	check(err)
	_, err = stmt.ExecContext(context.Background(), "OFF")
	check(err)

	status2 := GetStatus(tb_name, connstr)
	fmt.Println(status)

	cursor_modified1 := NewCursor(headers, status2)

	// Add
	stage.Add(NewCursorFromProllyTree(cg.GetHeadCommit().Value), cursor_modified1)
	stage.PrintStatus()

	// Commit
	/* cg.NewCommit(time.Now(), "Tim", Encode(cursor_modified1.GetHash()), cursor_modified1.GetTree(), "Second Commit")
	stage.Commit()
	stage.PrintStatus() */

	rows, err := db.QueryContext(context.Background(), "select * from "+tb_name)
	check(err)
	fmt.Println(readRows(rows))
	// (Unstage)RollBack
	stage.Unstage(25, tb_name, connstr)

	rows, err = db.QueryContext(context.Background(), "select * from "+tb_name)
	check(err)
	fmt.Println(readRows(rows))

	stmt, err = db.PrepareContext(context.Background(), "DROP VIEW "+tb_name)
	check(err)
	_, err = stmt.ExecContext(context.Background())
	check(err)

	stmt, err = db.PrepareContext(context.Background(), "DROP TABLE "+tb_name+"__backend")
	check(err)
	_, err = stmt.ExecContext(context.Background())
	check(err)

	defer db.Close()

}

// func Test() {
// 	// define database name and table name
// 	db_name := "KoreanShowbiz"
// 	tb_name := "Member"
// 	db_password := ""

// 	// connect mysql
// 	db, err := sql.Open("mysql", "root:"+db_password+"@tcp(127.0.0.1:3306)/"+db_name+"?charset=utf8&parseTime=True")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// read schemas
// 	// rows, err := db.Query("describe " + tb_name)

// 	// get column names
// 	rows, err := db.Query("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE `TABLE_SCHEMA`='" + db_name + "' AND `TABLE_NAME`='" + tb_name + "'")

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	tmp := readRows(rows)

// 	var headers []string
// 	for _, val := range tmp {
// 		headers = append(headers, val[0])
// 	}

// 	// fmt.Println(headers)

// 	// get rows
// 	rows, err = db.Query("select * from " + tb_name)

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	datas := readRows(rows)
// 	orig_data := make([][]string, len(datas))
// 	modified_1 := make([][]string, len(datas))

// 	c := copy(orig_data, datas)
// 	c = copy(modified_1, datas)
// 	_ = c
// 	fmt.Println(orig_data[6])
// 	orig_data[6] = []string{"1", "2", "HI", "4", "5", "6"}
// 	fmt.Println(orig_data[11])
// 	modified_1[11] = []string{"1", "2", "XD", "4", "5", "6"}
// 	modified_1 = append(modified_1, []string{"1", "2", "包豪斯", "4", "5", "6"})
// 	// modified_1 = append(modified_1, []string{"1", "2", "任炫植", "4", "5", "6"})
// 	test := make([]string, len(modified_1[8]))
// 	c = copy(test, modified_1[8])
// 	_ = c
// 	test[1] = "SHINEE"
// 	modified_1[8] = test

// 	cg := NewCommitGraph()
// 	stage := NewStage()

// 	// fmt.Println(orig_data)
// 	// fmt.Println(modified_1)

// 	// first commit
// 	cursor_orig := NewCursor(headers, orig_data)
// 	cg.NewCommit(time.Now(), "Tim", Encode(cursor_orig.GetHash()), cursor_orig.GetTree(), "orig_data")

// 	fmt.Println(cg.GetHeadCommit().Author)
// 	fmt.Println(cg.GetHeadCommit().Message)

// 	cursor_modified1 := NewCursor(headers, modified_1)

// 	stage.Add(NewCursorFromProllyTree(cg.GetHeadCommit().Value), cursor_modified1)
// 	stage.PrintStatus()
// 	stage.Commit()
// 	stage.PrintStatus()

// 	// Init()

// 	var connstr = "postgres://noidname:test@localhost:5432/" + db_name
// 	status := GetStatus("Member", connstr)
// 	fmt.Println(status)

// 	/* cursor_modified1 := NewCursor(headers, modified_1)

// 	stage.Add(cg., cursor_modified1)
// 	stage.PrintStatus() */

// 	/* cursor_modified1 := NewCursor(headers, modified_1)

// 	stage := NewStage()

// 	stage.Add(cursor_orig, cursor_modified1)
// 	stage.PrintStatus()

// 	modified_2 := append(modified_1, []string{"1", "2", "BTS", "4", "5", "6"})
// 	modified_2[11] = []string{"1", "2", "XD", "4", "666", "6"}
// 	cursor_modified2 := NewCursor(headers, modified_2)

// 	stage.Status(cursor_orig, cursor_modified2)
// 	stage.PrintStatus()
// 	stage.Add(cursor_orig, cursor_modified2)
// 	stage.PrintStatus()
// 	stage.Commit()
// 	stage.PrintStatus()

// 	commit_hash := Encode(cursor_orig.GetHash())
// 	fmt.Println(commit_hash)
// 	fmt.Println(cursor_orig.GetHash())
// 	fmt.Println(Decode(commit_hash)) */

// 	defer db.Close()
// }

func readRows(rows *sql.Rows) [][]string {
	var row_s [][]string

	cols, err := rows.Columns()

	if err != nil {
		fmt.Println(err)
		return row_s
	}

	vals := make([]interface{}, len(cols))

	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			return row_s
		}
		var row []string
		rowString := ""
		for _, col := range vals {
			switch col.(type) {
			case *sql.RawBytes:
				rowString = string(*col.(*sql.RawBytes))
				row = append(row, rowString)
				break
			}
		}
		row_s = append(row_s, row)
	}
	return row_s
}
