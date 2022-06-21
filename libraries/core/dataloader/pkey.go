package dataloader

import (
	. "Digit/libraries/core/diff"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
)

type pkeySaver struct {
	Key           int
	NextLevelHash string
	Data          map[string]string
	NextPkey      string //filename
}

type chunkSaver struct {
	Hash        string
	Chunkmap    map[string]string
	NextChunk   string //encoded Hash
	ParentChunk string //encoded Hash
	HeadPrikey  string
}

type treeSaver struct {
	Tree   []map[string]string
	Heads  []string //base64 of chunks
	Pkeys  []map[string]string
	Lastid int
}

// prikeymap[nextlevelhash] = (base64 of pkeysaver)
func savePrikey(pkey *PriKey, prikeymap map[string]string) (nxtlvhash string) {
	if pkey == nil {
		return "nil"
	} else if _, ok := prikeymap[Encode(pkey.GetHash())]; ok {
		return Encode(pkey.GetHash())
	} else {
		maptosave := make(map[string]string, len(pkey.GetData()))
		for key, ele := range pkey.GetData() {
			maptosave[key] = strconv.Itoa(ele)
		}
		nxtlvhash = Encode(pkey.GetHash())
		s := pkeySaver{
			Key:           int(pkey.GetKey()),
			NextLevelHash: nxtlvhash,
			Data:          maptosave,
			NextPkey:      savePrikey(pkey.GetNext(), prikeymap),
		}

		b, err := json.Marshal(s)
		check(err)
		encstr := base64.URLEncoding.EncodeToString(b)
		prikeymap[nxtlvhash] = encstr

		return nxtlvhash
	}

}

// totalmap[ Encode(Chunk.hash) ] = (base64 of chunkSaver)
func saveChunk(chunk *Chunk, chunkmap map[string]string, prikeymap map[string]string) (hs string) {
	if chunk == nil {
		return "nil"
	} else if _, ok := chunkmap[Encode(chunk.GetHash())]; ok {
		return Encode(chunk.GetHash())
	} else {
		curhash := Encode(chunk.GetHash())
		s := chunkSaver{
			Hash:        curhash,
			NextChunk:   saveChunk(chunk.GetNext(), chunkmap, prikeymap),
			ParentChunk: saveChunk(chunk.GetParent(), chunkmap, prikeymap),
			HeadPrikey:  savePrikey(chunk.GetHeadPri(), prikeymap),
		}

		b, err := json.Marshal(s)
		check(err)
		encstr := base64.URLEncoding.EncodeToString(b)
		chunkmap[curhash] = encstr
		return curhash
	}
}

func SaveTree(tree ProllyTree) {
	tosaveChunks := make([]map[string]string, 0)
	tosaveHeads := make([]string, 0)
	tosavePkeys := make([]map[string]string, 0)
	for _, chunk := range tree.GetheadChunks() {
		mchunk := make(map[string]string, 1)
		mpkey := make(map[string]string, 1)
		hs := saveChunk(&chunk, mchunk, mpkey)
		tosaveChunks = append(tosaveChunks, mchunk)
		tosaveHeads = append(tosaveHeads, hs)
		tosavePkeys = append(tosavePkeys, mpkey)
	}
	s := treeSaver{
		Tree:   tosaveChunks,
		Heads:  tosaveHeads,
		Pkeys:  tosavePkeys,
		Lastid: tree.Lastid,
	}
	b, err := json.Marshal(s)
	check(err)
	fmt.Println(string(b))

}

// func loadTree(b []byte) ProllyTree {
// 	var s treeSaver
// 	json.Unmarshal(b, &s)

// 	return p
// }
func check(e error) {
	if e != nil {
		panic(e)
	}
}
