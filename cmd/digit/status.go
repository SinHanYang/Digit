package digit

import "github.com/spf13/cobra"

func Status(args []string) {
	// TODO: Status implement
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the current branch",
	Long:  `Show the status of the current branch`,
	Run: func(cmd *cobra.Command, args []string) {
		Status(args)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
