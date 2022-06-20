package digit

import (
	env "Digit/libraries/core/env"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var COMMIT_API = "http://localhost:17617/commit"
var commit_message string

func Commit(rest []string) {
	// TODO: Commit implement
	if commit_message == "" {
		log.Fatal("commit message cannot be empty")
	}
	data := DigitRequest{
		Commit_message: commit_message,
	}
	data.DB_NAME, data.DB_USER, data.DB_PASS, _ = env.GetConfig(".")
	requestBody, _ := json.MarshalIndent(data, "", " ")
	req, _ := http.NewRequest("POST", COMMIT_API, bytes.NewBuffer(requestBody))
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
		log.Println(response.Message)
	}

}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit",
	Long:  "commit",
	Run: func(cmd *cobra.Command, args []string) {
		Commit(args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&commit_message, "message", "m", "", "commit message")
	commitCmd.MarkFlagRequired("message")
}
