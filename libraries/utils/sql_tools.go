package utils

import (
	"database/sql"
	"fmt"
)

func ReadRows(rows *sql.Rows) [][]string {
	var row_s [][]string

	cols, err := rows.Columns()
	row_s = append(row_s, cols)
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
