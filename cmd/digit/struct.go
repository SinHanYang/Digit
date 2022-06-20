package digit

type ConfigFile struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
}

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
