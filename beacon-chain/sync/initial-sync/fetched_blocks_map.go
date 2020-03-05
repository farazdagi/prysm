package initialsync

import (
	"sync"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// fetchedBlock is a container for signed beacon block.
type fetchedBlock struct {
	slot    uint64
	skipped bool
	block   *eth.SignedBeaconBlock
}

// fetchedBlocksMap is a storage where incoming blocks are aggregated in a thread-safe way.
type fetchedBlocksMap struct {
	sync.RWMutex
	blocks map[uint64]*fetchedBlock
}

// newFetchedBlocksMap constructs fully initialized map.
func newFetchedBlocksMap() *fetchedBlocksMap {
	return &fetchedBlocksMap{
		blocks: make(map[uint64]*fetchedBlock),
	}
}

// load returns value (if present) and its status (whether value is present or not).
func (m *fetchedBlocksMap) load(key uint64) (*fetchedBlock, bool) {
	m.RLock()
	val, ok := m.blocks[key]
	m.RUnlock()
	return val, ok
}

// pop removes and returns a value and its status (whether value is present or not).
func (m *fetchedBlocksMap) pop(key uint64) (*fetchedBlock, bool) {
	m.Lock()
	val, ok := m.blocks[key]
	if ok {
		delete(m.blocks, key)
	}
	m.Unlock()
	return val, ok
}

// store saves value for a given key.
func (m *fetchedBlocksMap) store(key uint64, value *fetchedBlock) {
	m.Lock()
	m.blocks[key] = value
	m.Unlock()
}

// storeBlocks saves list of blocks within a single locked operation.
func (m *fetchedBlocksMap) storeBlocks(blocks []*fetchedBlock) {
	m.Lock()
	for _, block := range blocks {
		m.blocks[block.slot] = block
	}
	m.Unlock()
}

// delete removes value (if one exists).
func (m *fetchedBlocksMap) delete(key uint64) {
	if _, ok := m.load(key); ok {
		m.Lock()
		delete(m.blocks, key)
		m.Unlock()
	}
}

// delete removes value (if one exists).
func (m *fetchedBlocksMap) deleteRange(start, end uint64) {
	m.Lock()
	for key := start; key < end; key++ {
		delete(m.blocks, key)
	}
	m.Unlock()
}

// len returns size of the map.
func (m *fetchedBlocksMap) len() int {
	m.RLock()
	size := len(m.blocks)
	m.RUnlock()
	return size
}

// populated checks whether a given [start, start + count) range of keys has been populated.
func (m *fetchedBlocksMap) populated(start, count uint64) bool {
	m.RLock()
	defer m.RUnlock()

	for i := start; i < start+count; i++ {
		if _, ok := m.blocks[i]; !ok {
			return false
		}
	}

	return true
}
