package digit

import (
	env "Digit/libraries/core/env"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var DIFF_API = "http://localhost:17617/diff"

func Diff(d bool, s bool, summary bool, args []string) {
	// TODO: Diff implement
	if len(args) < 2 {
		log.Fatal("diff command requires two arguments")
	}
	data := DigitRequest{
		Diff_from: args[0],
		Diff_to:   args[1],
	}
	data.DB_NAME, data.DB_USER, data.DB_PASS, _ = env.GetConfig(".")
	requestBody, _ := json.MarshalIndent(data, "", " ")
	req, _ := http.NewRequest("POST", DIFF_API, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response body as json
	var response DigitResponse
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &response)
	if response.Message != "" {
		// log the message
		fmt.Println(response.Message)
	}

}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show the changes between two commits",
	Long:  `Show the changes between two commits`,
	Run: func(cmd *cobra.Command, args []string) {
		d, _ := cmd.Flags().GetBool("data")
		s, _ := cmd.Flags().GetBool("schema")
		summary, _ := cmd.Flags().GetBool("summary")
		Diff(d, s, summary, args)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolP("data", "d", true, "Show only the data changes, do not show the schema changes (Both shown by default).")
	diffCmd.Flags().BoolP("schema", "s", true, "Show only the schema changes, do not show the data changes (Both shown by default).")
	diffCmd.Flags().BoolP("summary", "", false, "Show a summary of the changes.")
}
