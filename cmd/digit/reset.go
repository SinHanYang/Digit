package digit

import "github.com/spf13/cobra"

func Reset(hard bool, args []string) {
	// TODO: Reset implement
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database to specified commit",
	Long:  `Reset the database to specified commit`,
	Run: func(cmd *cobra.Command, args []string) {
		hard, _ := cmd.Flags().GetBool("hard")
		Reset(hard, args)

	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().BoolP("hard", "", false, "commit to reset to")
}
