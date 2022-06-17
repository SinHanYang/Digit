package digit

import (
	env "Digit/libraries/core/env"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var FILE_MODE = os.ModePerm

func Init(path string) {
	if env.HasDigitDir(path) {
		log.Fatal("This directory has already been initialized.")
		return
	}

	if path == "" {
		return
	}

	//confirm the path is vaild and the dir is empty
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, FILE_MODE)
		if err != nil {
			println("Path: " + path)
			log.Fatal(err)
		}
	} else {
		if stat.IsDir() {
			dir, err := ioutil.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			} else {
				if len(dir) != 0 {
					log.Fatalf("dir %v is not empty", path)
				}
			}
		} else {
			log.Fatalf("%v is a file", path)
		}
	}

	// create dir .digit
	digitPath := filepath.Join(path, ".digit")
	err := os.Mkdir(digitPath, FILE_MODE)
	if err != nil {
		log.Fatal(err)
	}

	// create file: config, description, HEAD
	os.Create(filepath.Join(digitPath, "config"))
	os.Create(filepath.Join(digitPath, "description"))
	head, _ := os.Create(filepath.Join(digitPath, "HEAD"))
	head.Write([]byte("ref: refs/heads/master"))

	// create dir: hooks, info, object, refs
	os.Mkdir(filepath.Join(digitPath, "hooks"), FILE_MODE)
	os.Mkdir(filepath.Join(digitPath, "info"), FILE_MODE)
	os.Mkdir(filepath.Join(digitPath, "object"), FILE_MODE)
	os.Mkdir(filepath.Join(digitPath, "refs"), FILE_MODE)

	// create dir tags and heads in refs
	refsPath := filepath.Join(digitPath, "refs")
	os.Mkdir(filepath.Join(refsPath, "tags"), FILE_MODE)
	os.Mkdir(filepath.Join(refsPath, "heads"), FILE_MODE)

	fmt.Printf("Successfully initialized Digit repository.")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a empty Digit data repository",
	Long: `This command creates an empty Digit data repository in the chosen directory. The directory must be empty. 
Running digit init in an already initialized directory will fail.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Initializing a empty Digit data repository...")
		if len(args) > 0 {
			Init(args[0])
		} else {
			Init(".")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
