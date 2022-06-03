package digit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func LsFiles(s bool) {
	if !s {
		fmt.Printf("ls-files only support --stage now!")
		return
	}
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return
	}
	//read entry-list from index
	entryList := getEntryListFromIndex()

	for _, entry := range entryList.List {
		fmt.Printf("%s %s %d	%s\n", entry.Mode, entry.Sha1, entry.Num, entry.Path)
	}
}

var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "list files in the index",
	Long:  `list files in the index`,
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := cmd.Flags().GetBool("stage")
		LsFiles(s)
	},
}

func init() {
	rootCmd.AddCommand(lsFilesCmd)
	lsFilesCmd.Flags().BoolP("stage", "s", false, "show staged contents' mode bits, object name and stage number")
}
