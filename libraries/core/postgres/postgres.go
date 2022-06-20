package postgres

import (
	"Digit/libraries/core/diff"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// TODO
var db_name = "KoreanShowbiz"
var connstr = "postgres://noidname:test@localhost:5432/" + db_name

// func main() {
// 	status := getStatus("teammembers", connstr)
// 	fmt.Println(status)

// 	var newstatus [][2]int
// 	fmt.Println(status)
// 	for _, ss := range status {
// 		ss[1] = 1
// 		newstatus = append(newstatus, ss)
// 	}
// 	modifyStatus("teammembers", newstatus, connstr)

// 	rollback(4, "teammembers", connstr)
// 	fmt.Println(getStatus("teammembers", connstr))
// }

func Init() {
	readSqlFile("func.sql")
	readSqlFile("start.sql")
	readSqlFile("test.sql")
}

func readSqlFile(filename string) {
	path := filepath.Join("libraries", "core", "postgres", filename)
	c, ioErr := ioutil.ReadFile(path)
	if ioErr != nil {
		fmt.Println(ioErr)
	}
	conn, err := pgx.Connect(context.Background(), connstr)
	if err != nil {
		fmt.Println(err)
	}
	_, err = conn.Exec(context.Background(), string(c))
	if err != nil {
		fmt.Println(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func GetStatus(view_name string, connstr string) [][2]diff.KeyType {
	db, err := sql.Open("pgx", connstr)
	check(err)
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), "select digitsnid, digitstatus from "+view_name+"__backend order by digitsnid asc")
	check(err)
	defer rows.Close()
	cols, err := rows.Columns()
	check(err)
	fmt.Println(cols)
	var status [][2]diff.KeyType
	for rows.Next() {
		var ss [2]diff.KeyType
		err := rows.Scan(&ss[0], &ss[1])
		check(err)
		status = append(status, ss)
	}
	return status
}

func modifyStatus(view_name string, status [][2]diff.KeyType, connstr string) {
	db, err := sql.Open("pgx", connstr)
	check(err)
	defer db.Close()
	stmt, err := db.PrepareContext(context.Background(),
		"DROP TRIGGER hook_update_"+view_name+" ON "+view_name+"__backend")
	check(err)
	_, err = stmt.ExecContext(context.Background())
	check(err)
	// fmt.Println(res.RowsAffected())
	for _, ss := range status {
		ctx := context.Background()
		stmt, err = db.PrepareContext(ctx, fmt.Sprint("UPDATE ", view_name, "__backend", " SET digitstatus=", ss[1], " WHERE digitsnid=", ss[0]))
		check(err)
		_, err := stmt.ExecContext(ctx)
		check(err)
		// fmt.Println(res.RowsAffected())
		// fmt.Println(ss[0], ss[1])
	}
	stmt, err = db.PrepareContext(context.Background(),
		"CREATE TRIGGER hook_update_"+view_name+" BEFORE UPDATE OR DELETE ON "+view_name+"__backend FOR EACH ROW EXECUTE FUNCTION hook_update()")
	check(err)
	_, err = stmt.ExecContext(context.Background())
	check(err)
	// fmt.Println(res.RowsAffected())
}

func Rollback(lastid diff.KeyType, view_name string, connstr string) {
	var setvalint diff.KeyType
	var setvalbool bool
	setvalbool = true
	setvalint = lastid
	db, err := sql.Open("pgx", connstr)
	check(err)
	defer db.Close()
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx,
		"DROP TRIGGER hook_update_"+view_name+" ON "+view_name+"__backend")
	check(err)
	_, err = stmt.ExecContext(ctx)
	check(err)
	stmt, err = db.PrepareContext(ctx, fmt.Sprint("DELETE FROM ", view_name, "__backend", " WHERE digitsnid>", lastid))
	check(err)
	res, err := stmt.ExecContext(ctx)
	check(err)
	fmt.Println(res.RowsAffected())
	if lastid == 0 {
		setvalint = 1
		setvalbool = false
	}
	stmt, err = db.PrepareContext(ctx, fmt.Sprint("SELECT setval(pg_get_serial_sequence('", view_name, "__backend'", ",'digitsnid'),", setvalint, ",", setvalbool, ")"))
	check(err)
	_, err = stmt.ExecContext(ctx)
	check(err)
	stmt, err = db.PrepareContext(context.Background(),
		"CREATE TRIGGER hook_update_"+view_name+" BEFORE UPDATE OR DELETE ON "+view_name+"__backend FOR EACH ROW EXECUTE FUNCTION hook_update()")
	check(err)
	_, err = stmt.ExecContext(context.Background())
	check(err)
}
