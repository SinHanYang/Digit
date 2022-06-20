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

var ADD_API = "http://localhost:17617/add"

func Add(a bool, args []string) {
	//TODO: Add implement
	// TODO: fetch data and call diff/add.go api

	if !env.HasDigitDir("./") {
		log.Fatal("This isn't a Digit repo, please `digit init` first.")
	} else {
		if len(args) == 0 && !a {
			log.Fatal("Nothing specified, nothing added.\n Maybe you wanted to say 'dolt add .'?")
		} else if a || len(args) == 1 && args[0] == "." {
			//TODO: stage all data current table
			data := DigitRequest{
				Add_all:   true,
				Add_table: args,
			}
			data.DB_NAME, data.DB_USER, data.DB_PASS, _ = env.GetConfig(".")
			requestBody, _ := json.MarshalIndent(data, "", " ")
			req, _ := http.NewRequest("POST", ADD_API, bytes.NewBuffer(requestBody))
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
			if response.Status == "error" && response.Message != "" {
				// log the message
				log.Fatal(response.Message)
			}
			fmt.Println("Add table success!")
		} else {
			// TODO: stage tables
			data := DigitRequest{
				Add_all:   false,
				Add_table: args,
			}
			data.DB_NAME, data.DB_USER, data.DB_PASS, _ = env.GetConfig(".")
			requestBody, _ := json.MarshalIndent(data, "", " ")
			req, _ := http.NewRequest("POST", ADD_API, bytes.NewBuffer(requestBody))
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
			if response.Status == "error" && response.Message != "" {
				// log the message
				log.Fatal(response.Message)
			}
			fmt.Println("Add table success!")
		}
	}
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Stage table",
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := cmd.Flags().GetBool("all")
		Add(a, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("all", "A", false, "add all table in the working root")
}
