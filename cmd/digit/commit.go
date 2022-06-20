package digit

import (
	"log"

	"github.com/spf13/cobra"
)

func Commit(m string, rest []string) {
	// TODO: Commit implement
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
		m, _ := cmd.Flags().GetString("message")
		Commit(m, args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "commit message")
	commitCmd.MarkFlagRequired("message")
}
