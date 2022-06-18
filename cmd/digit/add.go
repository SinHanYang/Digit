package digit

import (
	env "Digit/libraries/core/env"
	"log"

	"github.com/spf13/cobra"
)

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
		} else {
			// TODO: stage tables
			// roots, err = actions.StageTables(ctx, roots, docs, tables)
		}
	}
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a file to the index",
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := cmd.Flags().GetBool("all")
		Add(a, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("all", "A", false, "add all files in the working directory")
}
