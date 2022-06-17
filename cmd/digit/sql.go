package digit

import "github.com/spf13/cobra"

func Sql(q string, args []string) {
	// TODO: Sql implement
}

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Runs a SQL query",
	Long:  `Runs a SQL query`,
	Run: func(cmd *cobra.Command, args []string) {
		q, _ := cmd.Flags().GetString("data")
		Sql(q, args)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringP("query", "q", "", "Runs a single query and return the response.")
}
