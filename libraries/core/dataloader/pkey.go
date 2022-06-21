package dataloader

import (
	. "Digit/libraries/core/diff"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type pkeySaver struct {
	Key           int
	NextLevelHash string
	Data          map[string]int
	NextPkey      string //filename
}

type chunkSaver struct {
	Hash        string
	Chunkmap    []string
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
		nxtlvhash = Encode(pkey.GetHash())
		s := pkeySaver{
			Key:           int(pkey.GetKey()),
			NextLevelHash: nxtlvhash,
			Data:          pkey.GetData(),
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
		a := make([]string, 0)
		for k, _ := range chunk.GetChunkMap() {
			a = append(a, Encode(k))
		}
		s := chunkSaver{
			Hash:        curhash,
			Chunkmap:    a,
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

func SaveTree(tree ProllyTree) []byte {
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
	return b

}

func loadPrikey(b64 string, encodedmap map[string]string, resultmap map[ChunkAddress]PriKey) *PriKey {
	b, err := base64.URLEncoding.DecodeString(b64)
	var s pkeySaver
	check(err)
	json.Unmarshal(b, &s)
	ca := Decode(s.NextLevelHash)
	if p, ok := resultmap[ca]; ok {
		return &p
	} else {
		if s.NextPkey != "nil" {
			p := NewPrikey(s.Key, s.NextLevelHash, s.Data,
				loadPrikey(
					encodedmap[s.NextPkey],
					encodedmap,
					resultmap,
				),
			)
			resultmap[ca] = p
			return &p
		} else {
			p := NewPrikey(s.Key, s.NextLevelHash, s.Data, nil)
			resultmap[ca] = p
			return nil
		}
	}
}

func loadChunk(b64 string, prikeymap map[ChunkAddress]PriKey, encChunkmap map[string]string, resultmap map[ChunkAddress]Chunk) *Chunk {
	b, err := base64.URLEncoding.DecodeString(b64)
	check(err)
	var s chunkSaver
	json.Unmarshal(b, &s)
	hs := Decode(s.Hash)
	if c, ok := resultmap[hs]; ok {
		return &c
	} else {
		var ncp, npp *Chunk
		if s.NextChunk != "nil" {
			ncp = loadChunk(encChunkmap[s.NextChunk], prikeymap, encChunkmap, resultmap)
		} else {
			ncp = nil
		}
		if s.ParentChunk != "nil" {
			npp = loadChunk(encChunkmap[s.ParentChunk], prikeymap, encChunkmap, resultmap)
		} else {
			npp = nil
		}
		mpkey := make(map[ChunkAddress]PriKey, 0)
		for _, k := range s.Chunkmap {
			dk := Decode(k)
			mpkey[dk] = prikeymap[dk]
		}
		pk := prikeymap[Decode(s.HeadPrikey)]
		c := NewLoadChunk(
			hs,
			mpkey,
			ncp,
			npp,
			&pk,
		)
		resultmap[hs] = c
		return &c
	}
}

func LoadTree(b []byte) ProllyTree {
	var s treeSaver
	json.Unmarshal(b, &s)
	tree := make([]map[ChunkAddress]Chunk, 0)
	headchunks := make([]Chunk, 0)
	for index, headChunkaddr := range s.Heads {
		pkeymap := make(map[ChunkAddress]PriKey, len(s.Pkeys[index]))
		eleintree := make(map[ChunkAddress]Chunk, len(s.Tree[index]))
		pkeyhash_encpkey := s.Pkeys[index]
		for _, pkeyenc := range pkeyhash_encpkey {
			loadPrikey(pkeyenc, pkeyhash_encpkey, pkeymap)
		}
		Chunkaddr_encChunk := s.Tree[index]
		for Chunkaddr, encChunk := range Chunkaddr_encChunk {
			eleintree[Decode(Chunkaddr)] = *(loadChunk(encChunk, pkeymap, Chunkaddr_encChunk, eleintree))
		}
		tree = append(tree, eleintree)
		headchunks = append(headchunks, eleintree[Decode(headChunkaddr)])
	}
	p := NewLoadProllyTree(
		tree,
		headchunks,
		s.Lastid,
	)
	return p
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
