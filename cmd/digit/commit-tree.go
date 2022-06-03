package digit

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

func CommitTree(treeObjSha1 string, p string, m string, rest []string) *CommitOjbect {
	// Todo: CommitTree implement
	// Below is a sample implementation from Git原理教學
	if treeObjSha1 == "" || len(treeObjSha1) < 4 {
		log.Fatalf("Not a valid object name %s\n", treeObjSha1)
	}
	var commitObj CommitOjbect

	// get the whole sha1 of the tree object
	exist, treeObjSha1 := isObjectExist(treeObjSha1)
	if !exist {
		log.Fatalf("Not a valid object name %s\n", treeObjSha1)
	}
	commitObj.treeObjSha1 = treeObjSha1

	if p != "" {
		exist, parentSha1 := isObjectExist(p)
		if !exist {
			log.Fatalf("The parent commit object is not exist!")
		}
		commitObj.parent = parentSha1
	}

	//in fact, we need to read info below from the config file ψ(._. )>
	commitObj.author = "liuyj24<liuyijun2017@email.szu.edu.cn>"
	commitObj.committer = "liuyj24<liuyijun2017@email.szu.edu.cn>"

	commitObj.date = fmt.Sprintf("%s", time.Now())
	commitObj.message = getMessage(m, rest)

	objSha1, data := getSha1AndRawData(&commitObj)
	commitObj.Sha1 = objSha1
	writeObject(objSha1, data)

	fmt.Printf("%s\n", objSha1)
	return &commitObj
}

func getMessage(m string, rest []string) string {
	for _, s := range rest {
		m += " " + s
	}
	return m
}

var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree",
	Short: "commit tree object",
	Long:  "commit tree object",
	Run: func(cmd *cobra.Command, args []string) {
		treeObjSha1 := args[0]
		rest := args[1:]
		p, _ := cmd.Flags().GetString("parent")
		m, _ := cmd.Flags().GetString("message")
		CommitTree(treeObjSha1, p, m, rest)
	},
}

func init() {
	rootCmd.AddCommand(commitTreeCmd)
	commitTreeCmd.Flags().StringP("parent", "p", "", "indicates the id of a parent commit object")
	commitTreeCmd.Flags().StringP("message", "m", "", "the message of the commit")
}
