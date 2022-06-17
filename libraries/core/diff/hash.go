package diff

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
	"sort"
)

var encoding = base32.NewEncoding("0123456789abcdefghijklmnopqrstuv")

func encode(data []byte) string {
	return encoding.EncodeToString(data)
}

func decode(s string) []byte {
	slice, err := encoding.DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	return slice
}

func hashMap(data map[string]string) ChunkAddress {
	var keys []string
	b := new(bytes.Buffer)
	for key, _ := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(b, "%s=\"%s\"\n", k, data[k])
	}
	// fmt.Println(b.String())
	r := sha512.Sum512([]byte(b.String()))
	h := ChunkAddress{}
	copy(h[:], r[:20])
	return h
}

func rollingHash(hash [20]byte) int {
	sum := 0
	for _, v := range hash {
		sum += int(v)
	}
	return sum
}

func hashHash(content string) ChunkAddress {
	r := sha512.Sum512([]byte(content))
	h := ChunkAddress{}
	copy(h[:], r[:20])
	return h
}
