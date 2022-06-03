package digit

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func UpdateRef(args []string) {
	if len(args) < 3 {
		log.Fatal("missing args\n")
		return
	}
	path := args[1]
	objSha1 := args[2]
	path = filepath.Join(".git", path)

	_, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	_, objSha1 = isObjectExist(objSha1)
	err = ioutil.WriteFile(path, []byte(objSha1), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func isRefExist(ref string) (bool, string) {
	path := filepath.Join(".git", "refs", "heads", ref)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, ""
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return true, string(bytes)
}

var updateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "update commit ref to HEAD",
	Long:  "update commit ref to HEAD",
	Run: func(cmd *cobra.Command, args []string) {
		UpdateRef(args)
	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)
}
