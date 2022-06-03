package digit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "digit",
	Short: "digit - version control for database",
	Long:  `Digit is a super fancy version control system for your database.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello CLI")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
