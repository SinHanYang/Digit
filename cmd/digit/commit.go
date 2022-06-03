package digit

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func Commit(m string, rest []string) {
	// TODO: Commit implement
	// Below is a sample implementation from Git原理教學
	tree := WriteTree()

	refPathBytes, err := ioutil.ReadFile(filepath.Join(".git", "HEAD"))
	if err != nil {
		log.Fatal(err)
	}
	refPathStr := string(refPathBytes)
	i := strings.LastIndex(refPathStr, "/")

	refName := refPathStr[i+1:]
	exist, parentSha1 := isRefExist(refName)
	if !exist {
		log.Fatal("missing initial commit...")
	}
	m = getMessage(m, rest)
	commit := CommitTree(tree.Sha1, parentSha1, m, []string{})

	i = strings.Index(refPathStr, " ")
	refPath := refPathStr[i+1:]
	err = ioutil.WriteFile(filepath.Join(".git", refPath), []byte(commit.Sha1), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
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
