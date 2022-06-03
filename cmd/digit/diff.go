package digit

import "github.com/spf13/cobra"

func Diff(d bool, s bool, summary bool, args []string) {
	// TODO: Diff implement
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
