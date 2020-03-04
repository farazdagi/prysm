package initialsync

import (
	"sync"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/sirupsen/logrus"
)

// fetchedBlocksMap is a storage where incoming/fetched blocks are aggregated in a thread-safe way.
type fetchedBlocksMap struct {
	sync.RWMutex
	blocks map[uint64]*eth.SignedBeaconBlock
}

// newFetchedBlocksMap constructs fully initialized map.
func newFetchedBlocksMap() *fetchedBlocksMap {
	return &fetchedBlocksMap{
		blocks: make(map[uint64]*eth.SignedBeaconBlock),
	}
}

// load returns value (if present) and it status (whether it is present).
func (m *fetchedBlocksMap) load(key uint64) (*eth.SignedBeaconBlock, bool) {
	m.RLock()
	val, ok := m.blocks[key]
	m.RUnlock()
	return val, ok
}

// store saves value for a given key.
func (m *fetchedBlocksMap) store(key uint64, value *eth.SignedBeaconBlock) {
	m.Lock()
	log.WithFields(logrus.Fields{
		"map": m.blocks,
		"key": key,
	}).Debug("store requested")
	m.blocks[key] = value
	m.Unlock()
}

// delete removes value (if it exists).
func (m *fetchedBlocksMap) delete(key uint64) {
	log.WithFields(logrus.Fields{
		"map": m.blocks,
		"key": key,
	}).Debug("delete requested")
	if _, ok := m.load(key); ok {
		m.Lock()
		log.WithFields(logrus.Fields{
			"map": m.blocks,
			"key": key,
		}).Debug("delete processing")
		delete(m.blocks, key)
		m.Unlock()
	}
}

// len returns size of the map.
func (m *fetchedBlocksMap) len() int {
	m.RLock()
	size := len(m.blocks)
	m.RUnlock()
	return size
}
