package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// define database name and table name
	db_name := "KoreanShowbiz"
	tb_name := "Member"
	db_password := "107332021" // TODO modify to your password

	// connect mysql
	db, err := sql.Open("mysql", "root:"+db_password+"@tcp(127.0.0.1:3306)/"+db_name+"?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err)
		return
	}

	// read schemas
	// rows, err := db.Query("describe " + tb_name)

	// get column names
	rows, err := db.Query("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE `TABLE_SCHEMA`='" + db_name + "' AND `TABLE_NAME`='" + tb_name + "'")

	if err != nil {
		fmt.Println(err)
		return
	}

	tmp := readRows(rows)

	var headers []string
	for _, val := range tmp {
		headers = append(headers, val[0])
	}

	// fmt.Println(headers)

	// get rows
	rows, err = db.Query("select * from " + tb_name)

	if err != nil {
		fmt.Println(err)
		return
	}

	datas := readRows(rows)
	orig_data := make([][]string, len(datas))
	modified_1 := make([][]string, len(datas))

	c := copy(orig_data, datas)
	c = copy(modified_1, datas)
	_ = c
	fmt.Println(orig_data[6])
	orig_data[6] = []string{"1", "2", "HI", "4", "5", "6"}
	fmt.Println(orig_data[11])
	modified_1[11] = []string{"1", "2", "XD", "4", "5", "6"}
	modified_1 = append(modified_1, []string{"1", "2", "包豪斯", "4", "5", "6"})
	// modified_1 = append(modified_1, []string{"1", "2", "任炫植", "4", "5", "6"})
	test := make([]string, len(modified_1[8]))
	c = copy(test, modified_1[8])
	_ = c
	test[1] = "SHINEE"
	modified_1[8] = test

	// fmt.Println(orig_data)
	// fmt.Println(modified_1)
	cursor_orig := newCursor(headers, orig_data)
	cursor_modified1 := newCursor(headers, modified_1)

	stage := NewStage()

	stage.Add(cursor_orig, cursor_modified1)
	stage.PrintStatus()

	modified_2 := append(modified_1, []string{"1", "2", "BTS", "4", "5", "6"})
	modified_2[11] = []string{"1", "2", "XD", "4", "666", "6"}
	cursor_modified2 := newCursor(headers, modified_2)

	stage.Status(cursor_orig, cursor_modified2)
	stage.PrintStatus()
	stage.Add(cursor_orig, cursor_modified2)
	stage.PrintStatus()
	stage.Commit()
	stage.PrintStatus()

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
