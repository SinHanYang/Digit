package diff

import (
	"fmt"
	"sort"
)

// approximately 2 priKeys / chunk
var rollingHashThreshold = 128 * 20 * 2
var EmptyChunk Chunk
var EmptyPriKey PriKey

type PriKey struct {
	key           int
	nextLevelHash ChunkAddress
	data          map[string]int
	nextPriKey    *PriKey
}

func (myKey PriKey) GetKey() int {
	return myKey.key
}

func (myKey PriKey) GetData() map[string]int {
	return myKey.data
}

func (myKey PriKey) GetHash() ChunkAddress {
	return myKey.nextLevelHash
}

func (myKey PriKey) Less(otherKey Value) bool {
	return myKey.GetKey() < otherKey.GetKey()
}
func (myKey PriKey) GetNext() *PriKey { return myKey.nextPriKey }
func NewPrikey(k int, nLHash string, d map[string]int, nPrikey *PriKey) PriKey {
	return PriKey{
		key:           k,
		nextLevelHash: Decode(nLHash),
		data:          d,
		nextPriKey:    nPrikey,
	}
}

// leaves at ProllyTree[0]
// root at ProllyTree[len(ProllyTree)-1]
type ProllyTree struct {
	tree       []map[ChunkAddress]Chunk
	headChunks []Chunk
	Lastid     int
}

func (t ProllyTree) GetTree() []map[ChunkAddress]Chunk { return t.tree }
func (t ProllyTree) GetheadChunks() []Chunk            { return t.headChunks }
func NewLoadProllyTree(t []map[ChunkAddress]Chunk, h []Chunk, Lid int) ProllyTree {
	return ProllyTree{
		tree:       t,
		headChunks: h,
		Lastid:     Lid,
	}
}

type Chunk struct {
	hash        ChunkAddress
	chunkMap    map[ChunkAddress]PriKey
	nextChunk   *Chunk
	parentChunk *Chunk
	headPriKey  *PriKey
}

func (c Chunk) GetHash() ChunkAddress                { return c.hash }
func (c Chunk) GetNext() *Chunk                      { return c.nextChunk }
func (c Chunk) GetParent() *Chunk                    { return c.parentChunk }
func (c Chunk) GetHeadPri() *PriKey                  { return c.headPriKey }
func (c Chunk) GetChunkMap() map[ChunkAddress]PriKey { return c.chunkMap }
func NewLoadChunk(h ChunkAddress, cm map[ChunkAddress]PriKey, nc *Chunk, pc *Chunk, hp *PriKey) Chunk {
	return Chunk{
		hash:        h,
		chunkMap:    cm,
		nextChunk:   nc,
		parentChunk: pc,
		headPriKey:  hp,
	}
}

type ChunkCursor struct {
	hash          ChunkAddress
	isDone        bool
	prollyTree    ProllyTree
	currentChunk  Chunk
	currentPriKey PriKey
}

func (cursor *ChunkCursor) Current() Value {
	return cursor.currentPriKey
}

func (cursor *ChunkCursor) Next() {
	if !(cursor.currentChunk.nextChunk == nil && cursor.currentPriKey.nextPriKey == nil) {
		if cursor.currentPriKey.nextPriKey != nil {
			cursor.currentPriKey = *(cursor.currentPriKey.nextPriKey)
			return
		} else {
			if cursor.currentChunk.nextChunk != nil && cursor.currentPriKey.nextPriKey == nil {
				cursor.currentChunk = *cursor.currentChunk.nextChunk
				cursor.currentPriKey = *cursor.currentChunk.headPriKey
			}
		}
	} else {
		cursor.isDone = true
	}
}

func (cursor *ChunkCursor) Done() bool {
	return cursor.isDone
}

func (cursor *ChunkCursor) Path() []ChunkAddress {
	var l []ChunkAddress
	currentChunk := cursor.currentChunk
	l = append([]ChunkAddress{currentChunk.hash}, l...)
	for currentChunk.parentChunk != nil {
		currentChunk = *currentChunk.parentChunk
		l = append([]ChunkAddress{currentChunk.hash}, l...)
	}
	return l
}

func (cursor *ChunkCursor) GetHash() ChunkAddress {
	return cursor.hash
}

func (cursor *ChunkCursor) GetTree() ProllyTree {
	return cursor.prollyTree
}

// NextAtLevel(0) is equivalent to Next(). NextAtLevel(1) advances the cursor
// past the current entry on the first internal level of the tree---the one
// directly above the leaves. After advancing the cursor at an internal level,
// if the cursor is not `Done()`, the cursor points to the left-most entry
// accessible from the newly selected internal entry.
// TODO
func (cursor *ChunkCursor) NextAtLevel(level int) {
	if level == 0 {
		cursor.Next()
	} else {
		// if not, let cursor point at the left-most entry
		// accessible from the newly selected internal entry.
		// find ChunkAddress at level
		previousHash := cursor.currentChunk.hash
		currentChunk := cursor.currentChunk

		for i := 1; i <= level; i++ {
			previousHash = currentChunk.hash
			currentChunk = *currentChunk.parentChunk
		}

		// move next at level
		// fmt.Println(currentChunk)
		currentPriKey := currentChunk.chunkMap[previousHash]
		// fmt.Println(currentPriKey)
		if currentPriKey.nextPriKey == nil {
			if currentChunk.nextChunk != nil {
				currentChunk = *currentChunk.nextChunk
				currentPriKey = *currentChunk.headPriKey
			} else {
				fmt.Println("Level", level, "can't move next!, To level", level-1)
				/* cursor.currentChunk.nextChunk = nil
				cursor.currentPriKey.nextPriKey = nil */
				cursor.NextAtLevel(level - 1)
				return
			}
		} else {
			currentPriKey = *currentPriKey.nextPriKey
		}

		// move down to leave
		nextLevelHash := currentPriKey.nextLevelHash
		for i := level - 1; i >= 0; i-- {
			currentChunk = cursor.prollyTree.tree[i][nextLevelHash]
			currentPriKey = *currentChunk.headPriKey
			nextLevelHash = currentPriKey.nextLevelHash
		}
		cursor.currentChunk = currentChunk
		cursor.currentPriKey = currentPriKey

	}
}

func NewChunk(priKeys []PriKey) Chunk {
	content := ""
	chunkMap := make(map[ChunkAddress]PriKey, len(priKeys))
	for i, priKey := range priKeys {
		for _, b := range priKey.GetHash() {
			content += string(b)
		}
		if i == len(priKeys)-1 {
			priKeys[i].nextPriKey = nil
		} else {
			priKeys[i].nextPriKey = &priKeys[i+1]
		}
		chunkMap[priKey.nextLevelHash] = priKeys[i]
		// fmt.Println(priKeys[i].nextPriKey)
	}

	hash := hashHash(content)

	return Chunk{
		hash:        hash,
		chunkMap:    chunkMap,
		nextChunk:   nil,
		parentChunk: nil,
		headPriKey:  &priKeys[0],
	}

}

func NewCursor(header []string, rows [][2]int) ChunkCursor {
	t := NewProllyTree(header, rows)
	// _ = t
	// printProllyTree(t.tree)

	/* 	   	for _, v := range t.tree {
		for key, value := range v {
			fmt.Println(key)
			fmt.Println(value)
		}
	} */
	return NewCursorFromProllyTree(t)
}

func NewCursorFromProllyTree(t ProllyTree) ChunkCursor {
	return ChunkCursor{
		hash:          t.headChunks[len(t.headChunks)-1].hash,
		isDone:        false,
		prollyTree:    t,
		currentChunk:  t.headChunks[0],
		currentPriKey: *t.headChunks[0].headPriKey,
	}
}

func NewProllyTree(header []string, rows [][2]int) ProllyTree {

	if len(rows) == 1 {
		var newProllyTreeLevel [][]Chunk
		// init all priKeys
		var priKeys []PriKey
		data := make(map[string]int, len(rows))
		row := rows[0]
		key := row[0]
		for i := range row {
			data[header[i]] = row[i]
		}
		priKey := PriKey{
			key:           key,
			data:          data,
			nextPriKey:    nil,
			nextLevelHash: hashMap(data),
		}
		priKeys = append(priKeys, priKey)
		chunk := NewChunk(priKeys)
		newProllyTreeLevel = append(newProllyTreeLevel, []Chunk{chunk})

		var t []map[ChunkAddress]Chunk
		var heads []Chunk
		for _, v := range newProllyTreeLevel {
			m := make(map[ChunkAddress]Chunk)
			heads = append(heads, v[0])
			for i := range v {
				m[v[i].hash] = v[i]
			}
			t = append(t, m)
		}
		return ProllyTree{
			tree:       t,
			headChunks: heads,
			Lastid:     len(rows),
		}
	}

	var newProllyTreeLevel [][]Chunk
	// init all priKeys
	var priKeys []PriKey
	for _, row := range rows {
		data := make(map[string]int, len(rows))
		key := row[0]
		for i := range row {
			data[header[i]] = row[i]
		}
		priKey := PriKey{
			key:           key,
			data:          data,
			nextPriKey:    nil,
			nextLevelHash: hashMap(data),
		}
		priKeys = append(priKeys, priKey)
	}
	// sort by Less()
	sort.Slice(priKeys, func(i, j int) bool {
		return priKeys[i].key < priKeys[j].key
	})

	for _, v := range priKeys {
		fmt.Print(v.key, ",")
	}
	fmt.Println()

	// split by rolling hasher
	for len(priKeys) > 1 {
		var newPrikeys []PriKey
		var newChunks []Chunk

		currentHashSum := 0
		startAt := 0
		for i, priKey := range priKeys {
			currentHashSum += rollingHash(priKey.nextLevelHash)
			if currentHashSum >= rollingHashThreshold {
				newChk := NewChunk(priKeys[startAt : i+1])
				newChunks = append(newChunks, newChk)
				newPrikeys = append(newPrikeys,
					PriKey{
						key:           priKeys[startAt].key,
						data:          nil,
						nextPriKey:    nil,
						nextLevelHash: newChk.hash,
					})
				startAt = i + 1
				currentHashSum = 0
			}
		}
		// fmt.Println(startAt)
		// fmt.Println(len(priKeys))

		if startAt < len(priKeys) {
			newChk := NewChunk(priKeys[startAt:])
			newChunks = append(newChunks, newChk)
			newPrikeys = append(newPrikeys, PriKey{
				key:           priKeys[startAt].key,
				data:          nil,
				nextPriKey:    nil,
				nextLevelHash: newChk.hash,
			})
		}

		for i, _ := range newChunks {
			if i == len(newChunks)-1 {
				newChunks[i].nextChunk = nil
			} else {
				newChunks[i].nextChunk = &newChunks[i+1]
			}
		}
		/* for i := range newChunks {
			fmt.Println(i)
			for _, v := range newChunks[i].chunkMap {
				fmt.Println(v.key)
			}
			fmt.Println()
		} */
		newProllyTreeLevel = append(newProllyTreeLevel, newChunks)
		priKeys = newPrikeys

	}

	// maintain parent
	assignParents(newProllyTreeLevel)

	// insert in map
	var t []map[ChunkAddress]Chunk
	var heads []Chunk
	for _, v := range newProllyTreeLevel {
		m := make(map[ChunkAddress]Chunk)
		heads = append(heads, v[0])
		for i := range v {
			m[v[i].hash] = v[i]
		}
		t = append(t, m)
	}

	/* printProllyTree(t)
	fmt.Println() */

	return ProllyTree{
		tree:       t,
		headChunks: heads,
		Lastid:     len(rows),
	}
}

// print from level 0 to top
func printProllyTree(t []map[ChunkAddress]Chunk) {
	for addr := range t[0] {
		currentChunk := t[0][addr]
		fmt.Println(currentChunk.hash)
		fmt.Println(currentChunk.headPriKey.key)
		for currentChunk.parentChunk != nil {
			currentChunk = *currentChunk.parentChunk
			fmt.Println(currentChunk.hash)
			fmt.Println(currentChunk.headPriKey.key)
		}
		fmt.Println()
	}
}

func assignParents(t [][]Chunk) {
	for level := len(t) - 1; level >= 1; level-- {
		for i, chunk := range t[level] {
			for hash, _ := range chunk.chunkMap {
				for j, nextLevelChunk := range t[level-1] {
					if nextLevelChunk.hash == hash {
						// fmt.Println("Parent:")
						// fmt.Println(chunk)
						nextLevelChunk.parentChunk = &t[level][i]
						t[level-1][j] = nextLevelChunk
						// fmt.Println("Next:")
						// fmt.Println(nextLevelChunk)
						break
					}
				}
			}
		}
	}
}
