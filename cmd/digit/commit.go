package digit

import (
	"github.com/spf13/cobra"
)

func Commit(m string, rest []string) {
	// TODO: Commit implement

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
}
