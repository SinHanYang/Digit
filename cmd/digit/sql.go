package digit

import (
	"Digit/libraries/core/env"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var SQL_API = "http://localhost:17617/sql"
var query string

func Sql(q string, args []string) {
	// TODO: Sql implement
	if q == "" {
		log.Fatal("query cannot be empty")
	}
	if !env.HasDigitDir("./") {
		log.Fatal("This isn't a Digit repo, please `digit init` first.")
	} else {
		// call SQL api
		data := DigitRequest{
			Sql_query: q,
		}
		data.DB_NAME, data.DB_USER, data.DB_PASS, _ = env.GetConfig(".")
		// fmt.Println(data)
		requestBody, _ := json.MarshalIndent(data, "", " ")
		req, _ := http.NewRequest("POST", SQL_API, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// read response body as json
		var response DigitResponse
		body, err := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &response)
		if response.Status == "error" && response.Message != "" {
			// log the message
			log.Fatal(response.Message)
		}

		if response.Data != nil {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			for i, v := range response.Data.([]interface{}) {
				if i == 0 {
					t.AppendHeader(v.([]interface{}))
				} else {
					t.AppendRow(v.([]interface{}))
				}
			}
			t.Render()
		}

	}
}

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Runs a SQL query",
	Long:  `Runs a SQL query`,
	Run: func(cmd *cobra.Command, args []string) {
		Sql(query, args)
	},
}

func init() {
	rootCmd.AddCommand(sqlCmd)
	sqlCmd.Flags().StringVarP(&query, "query", "q", "", "Runs a single query and return the response.")
}
