package initialsync

import (
	"math/rand"
	"sync"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

func TestFetchedBlocksMap(t *testing.T) {
	m := newFetchedBlocksMap()

	t.Run("concurrent access", func(t *testing.T) {
		var wg sync.WaitGroup
		funcs := []func(uint64){
			func(key uint64) {
				m.load(key)
				wg.Done()
			},
			func(key uint64) {
				m.delete(key)
				wg.Done()
			},
			func(key uint64) {
				m.store(key, &eth.SignedBeaconBlock{})
				if m.len() < 2 {
					m.store(key, &eth.SignedBeaconBlock{})
				}
				wg.Done()
			},
			func(key uint64) {
				m.load(key)
				wg.Done()
			},
			func(key uint64) {
				m.store(key, &eth.SignedBeaconBlock{})
				m.delete(key)
				wg.Done()
			},
		}

		wg.Add(5 * len(funcs))
		for i := 0; i < 5*len(funcs); i++ {
			go funcs[i%len(funcs)](rand.Uint64())
		}

		wg.Wait()
	})
}
