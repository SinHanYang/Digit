package digit

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Symbolic(args []string) {
	if len(args) < 3 || args[1] != "HEAD" {
		log.Fatal("error args...only support update HEAD now!")
		return
	}
	path := args[2]
	content := fmt.Sprintf("ref: %s", path)
	err := ioutil.WriteFile(filepath.Join(".git", "HEAD"), []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

var symbolicRefCmd = &cobra.Command{
	Use:   "symbolic-ref",
	Short: "update HEAD",
	Long:  "update HEAD",
	Run: func(cmd *cobra.Command, args []string) {
		Symbolic(args)
	},
}

func init() {
	rootCmd.AddCommand(symbolicRefCmd)
}
