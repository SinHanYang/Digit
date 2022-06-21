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
func Test1() {

	db_name := "testt"
	tb_name := "Member"
	username := "billy"
	password := "test"
	Init(db_name, username, password)
	cg := NewCommitGraph()
	stage := NewStage()
	_ = stage
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
	cg.NewCommit(time.Now(), "Tim", Encode(cursor_modified1.GetHash()), cursor_modified1.GetTree(), "Second Commit")
	stage.Commit()
	stage.PrintStatus()

	rows, err := db.QueryContext(context.Background(), "select * from "+tb_name)
	check(err)
	fmt.Println(readRows(rows))
	// (Unstage)RollBack
	stage.Unstage(cg.GetHeadCommit().Value.Lastid, tb_name, connstr)

	// print table
	rows, err = db.QueryContext(context.Background(), "select * from "+tb_name)
	check(err)
	fmt.Println(readRows(rows))

	// delete this table
	stmt, err = db.PrepareContext(context.Background(), "DROP VIEW "+tb_name)
	check(err)
	_, err = stmt.ExecContext(context.Background())
	check(err)
	// delete backend table
	stmt, err = db.PrepareContext(context.Background(), "DROP TABLE "+tb_name+"__backend")
	check(err)
	_, err = stmt.ExecContext(context.Background())
	check(err)

	defer db.Close()

}
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
