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

var RESET_API = "http://localhost:17617/reset"
var reset_tb string

func Reset(args []string) {
	// TODO: Reset implement
	if len(args) != 1 {
		log.Fatal("reset command requires one argument")
	}
	data := DigitRequest{
		Reset_hash: args[0],
		Reset_tb:   reset_tb,
	}
	data.DB_NAME, data.DB_USER, data.DB_PASS, _ = env.GetConfig(".")
	requestBody, _ := json.MarshalIndent(data, "", " ")
	req, _ := http.NewRequest("POST", RESET_API, bytes.NewBuffer(requestBody))
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

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database to specified commit",
	Long:  `Reset the database to specified commit`,
	Run: func(cmd *cobra.Command, args []string) {
		Reset(args)
	},
}

func init() {
	resetCmd.Flags().StringVarP(&reset_tb, "table", "t", "", "table to reset")
	rootCmd.AddCommand(resetCmd)
}
