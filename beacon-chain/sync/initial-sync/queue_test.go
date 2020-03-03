package initialsync

import (
	"context"
	"testing"

	dbtest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
)

func TestBlocksQueueInit(t *testing.T) {
	mc, p2p, beaconDB := initializeTestServices(t, []uint64{}, []*peerData{})
	fetcher := newBlocksFetcher(&blocksFetcherConfig{
		headFetcher: mc,
		p2p:         p2p,
	})
	ctx, _ := context.WithCancel(context.Background())

	t.Run("empty config", func(t *testing.T) {
		queue := newBlocksQueue(&blocksQueueConfig{})
		if queue.ctx == nil {
			t.Error("unexpected context: empty")
		}
	})

	t.Run("default startSlot", func(t *testing.T) {
		queue := newBlocksQueue(&blocksQueueConfig{
			headFetcher: mc,
		})
		if queue.startSlot == 0 {
			t.Errorf("unexpected startSlot, expected: %v, got: %v", 1, queue.startSlot)
		}
	})

	t.Run("normal init", func(t *testing.T) {
		startSlot, highestExpectedSlot := uint64(24), uint64(131)
		queue := newBlocksQueue(&blocksQueueConfig{
			ctx:                 ctx,
			blocksFetcher:       fetcher,
			headFetcher:         mc,
			startSlot:           startSlot,
			highestExpectedSlot: highestExpectedSlot,
		})
		if queue.ctx == nil {
			t.Error("unexpected context: empty")
		}
		if queue.startSlot != startSlot {
			t.Errorf("unexpected startSlot, expected: %v, got: %v", startSlot, queue.startSlot)
		}
		if queue.highestExpectedSlot != highestExpectedSlot {
			t.Errorf("unexpected highestExpectedSlot, expected: %v, got: %v", highestExpectedSlot, queue.highestExpectedSlot)
		}
	})

	dbtest.TeardownDB(t, beaconDB)
}
