package digit

import (
	env "Digit/libraries/core/env"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var LOG_API = "http://localhost:17617/log"

func Log(args []string) {
	// Todo: Log implementation
	data := DigitRequest{}
	data.DB_NAME, data.DB_USER, data.DB_PASS, _ = env.GetConfig(".")
	requestBody, _ := json.MarshalIndent(data, "", " ")
	req, _ := http.NewRequest("POST", LOG_API, bytes.NewBuffer(requestBody))
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

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "show commit log",
	Long:  "show commit log",
	Run: func(cmd *cobra.Command, args []string) {
		Log(args)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
