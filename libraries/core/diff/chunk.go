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
	key           string
	nextLevelHash ChunkAddress
	data          map[string]string
	nextPriKey    *PriKey
}

func (myKey PriKey) GetKey() string {
	return myKey.key
}

func (myKey PriKey) GetData() map[string]string {
	return myKey.data
}

func (myKey PriKey) GetHash() ChunkAddress {
	return myKey.nextLevelHash
}

func (myKey PriKey) Less(otherKey Value) bool {
	return myKey.GetKey() < otherKey.GetKey()
}

// leaves at ProllyTree[0]
// root at ProllyTree[len(ProllyTree)-1]
type ProllyTree struct {
	tree       []map[ChunkAddress]Chunk
	headChunks []Chunk
}

type Chunk struct {
	hash        ChunkAddress
	chunkMap    map[ChunkAddress]PriKey
	nextChunk   *Chunk
	parentChunk *Chunk
	headPriKey  *PriKey
}

type ChunkCursor struct {
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

func newChunk(priKeys []PriKey) Chunk {
	content := ""
	chunkMap := make(map[ChunkAddress]PriKey, len(priKeys))
	for i, priKey := range priKeys {
		content += priKey.key
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

func newCursor(header []string, rows [][]string) ChunkCursor {
	t := newProllyTree(header, rows)
	// _ = t
	// printProllyTree(t.tree)

	/* 	   	for _, v := range t.tree {
		for key, value := range v {
			fmt.Println(key)
			fmt.Println(value)
		}
	} */
	return ChunkCursor{
		isDone:        false,
		prollyTree:    t,
		currentChunk:  t.headChunks[0],
		currentPriKey: *t.headChunks[0].headPriKey,
	}
	// return ChunkCursor{}
}

func newProllyTree(header []string, rows [][]string) ProllyTree {

	var newProllyTreeLevel [][]Chunk
	// init all priKeys
	var priKeys []PriKey
	for _, row := range rows {
		data := make(map[string]string, len(rows))
		key := row[2]
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
				newChk := newChunk(priKeys[startAt : i+1])
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
			newChk := newChunk(priKeys[startAt:])
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
