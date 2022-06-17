package diff

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// define database name and table name
	db_name := "KoreanShowbiz"
	tb_name := ""

	// connect mysql
	db, err := sql.Open("mysql", "root:107332021@tcp(127.0.0.1:3306)/"+db_name+"?charset=utf8&parseTime=True")
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
	datas1 := make([][]string, len(datas))
	datas2 := make([][]string, len(datas))

	c := copy(datas1, datas)
	c = copy(datas2, datas)
	_ = c
	fmt.Println(datas1[6])
	datas1[6] = []string{"1", "2", "HI", "4", "5", "6"}
	fmt.Println(datas1[11])
	datas2[11] = []string{"1", "2", "XD", "4", "5", "6"}
	datas2 = append(datas2, []string{"1", "2", "包豪斯", "4", "5", "6"})
	// datas2 = append(datas2, []string{"1", "2", "任炫植", "4", "5", "6"})
	test := make([]string, len(datas2[8]))
	c = copy(test, datas2[8])
	_ = c
	test[1] = "SHINEE"
	datas2[8] = test

	// fmt.Println(datas1)
	// fmt.Println(datas2)
	cursor_left := newCursor(headers, datas1)
	cursor_right := newCursor(headers, datas2)

	diff := Diff(&cursor_left, &cursor_right)
	for _, v := range diff {

		if v.Op == Add {
			fmt.Print("ADD, ")
		} else if v.Op == Delete {
			fmt.Print("DELETE, ")
		} else if v.Op == Edit {
			fmt.Print("EDIT, ")
		}
		fmt.Println(v.Value.GetData())
	}

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
