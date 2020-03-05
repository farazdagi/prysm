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

func TestFetchedBlocksMapBlockRangeIsReady(t *testing.T) {
	tests := []struct {
		name         string
		start, count uint64
		slots        []uint64
		wantSlotsLen int
		want         bool
	}{
		{
			name:         "empty map",
			start:        0,
			count:        10,
			slots:        []uint64{},
			wantSlotsLen: 0,
			want:         false,
		},
		{
			name:         "duplicate keys",
			start:        0,
			count:        10,
			slots:        []uint64{0, 0, 1, 1, 2, 3, 3, 3, 2, 2, 2, 4, 5, 6, 7, 8, 9},
			wantSlotsLen: 10,
			want:         true,
		},
		{
			name:         "unordered keys",
			start:        5,
			count:        6,
			slots:        []uint64{22, 10, 9, 7, 4, 8, 1, 6, 3, 5, 2, 0, 1},
			wantSlotsLen: 12,
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newFetchedBlocksMap()
			for _, slot := range tt.slots {
				m.store(slot, &eth.SignedBeaconBlock{})
			}
			if got := m.len(); got != tt.wantSlotsLen {
				t.Errorf("invalid map len = %v, want %v", got, tt.wantSlotsLen)
			}

			if got := m.populated(tt.start, tt.count); got != tt.want {
				t.Errorf("populated() = %v, want %v", got, tt.want)
			}
		})
	}

}
