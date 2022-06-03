package digit

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func CatFile(p bool, t bool, s bool, args []string) {
	if !p {
		fmt.Printf("%s\n", getCatFileStr(p, t, s, args))
	} else {
		fmt.Printf("%s\n", getCatFileStr(p, t, s, args))
	}
}

func getCatFileStr(p bool, t bool, s bool, args []string) string {
	if !p && !t && !s {
		log.Fatal("a -p or -t is needed!")
	}

	//get data from database
	objectSha1 := args[len(args)-1]
	objectDir := filepath.Join(".git", "objects", objectSha1[:2])

	dir, err := ioutil.ReadDir(objectDir)
	if err != nil {
		log.Fatal(err)
	}
	var data []byte
	for _, file := range dir {
		if strings.HasPrefix(file.Name(), objectSha1[2:]) {
			data, err = ioutil.ReadFile(filepath.Join(objectDir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//uncompress
	raw := unCompressData(data)
	i := bytes.IndexByte(raw, ' ')
	j := bytes.IndexByte(raw, '\u0000')

	if t {
		objectType := raw[:i]
		return string(objectType)

	} else if s {
		objectSize := raw[i+1 : j]
		return string(objectSize)

	} else if p {
		objectContent := raw[j+1:]
		return string(objectContent)
	}
	return ""
}

func unCompressData(data []byte) []byte {
	reader := bytes.NewReader(data)
	r, err := zlib.NewReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	var out bytes.Buffer
	io.Copy(&out, r)
	return out.Bytes()
}

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Show object information",
	Long:  "Show object information",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p, _ := cmd.Flags().GetBool("print")
		t, _ := cmd.Flags().GetBool("type")
		s, _ := cmd.Flags().GetBool("size")
		CatFile(p, t, s, args)
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
	catFileCmd.Flags().BoolP("print", "p", false, "print object content")
	catFileCmd.Flags().BoolP("type", "t", false, "print object type")
	catFileCmd.Flags().BoolP("size", "s", false, "print object size")
}
