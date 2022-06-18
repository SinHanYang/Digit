package digit

import (
	"log"

	"github.com/spf13/cobra"
)

func Commit(a string, m string, rest []string) {
	// TODO: Commit implement
	if a == "" {
		log.Fatal("please tell me who you are")
	}
	if m == "" {
		log.Fatal("commit message cannot be empty")
	}
	// var hash Hash
	// hash.hash_value = "hash"
	// t := time.Now()
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit",
	Long:  "commit",
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := cmd.Flags().GetString("author")
		m, _ := cmd.Flags().GetString("message")
		Commit(a, m, args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("author", "a", "", "commit author")
	commitCmd.Flags().StringP("message", "m", "", "commit message")
	commitCmd.MarkFlagRequired("author")
	commitCmd.MarkFlagRequired("message")
}
