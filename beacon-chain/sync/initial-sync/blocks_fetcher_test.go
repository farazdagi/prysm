package initialsync

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	mock "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/db"
	dbtest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	p2pt "github.com/prysmaticlabs/prysm/beacon-chain/p2p/testing"
	stateTrie "github.com/prysmaticlabs/prysm/beacon-chain/state"
	p2ppb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/sliceutil"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func TestBlocksFetcherInitStartStop(t *testing.T) {
	mc, p2p, beaconDB := initializeTestServices(t, []uint64{}, []*peerData{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fetcher := newBlocksFetcher(
		ctx,
		&blocksFetcherConfig{
			headFetcher: mc,
			p2p:         p2p,
		})

	t.Run("check for leaked goroutines", func(t *testing.T) {
		err := fetcher.start()
		if err != nil {
			t.Error(err)
		}
		fetcher.stop() // should block up until all resources are reclaimed
		select {
		case <-fetcher.requestResponses():
		default:
			t.Error("fetchResponses channel is leaked")
		}
	})

	t.Run("re-starting of stopped fetcher", func(t *testing.T) {
		if err := fetcher.start(); err == nil {
			t.Errorf("expected error not returned: %v", errFetcherCtxIsDone)
		}
	})

	t.Run("multiple stopping attempts", func(t *testing.T) {
		fetcher := newBlocksFetcher(
			context.Background(),
			&blocksFetcherConfig{
				headFetcher: mc,
				p2p:         p2p,
			})
		if err := fetcher.start(); err != nil {
			t.Error(err)
		}

		fetcher.stop()
		fetcher.stop()
	})

	t.Run("cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		fetcher := newBlocksFetcher(
			ctx,
			&blocksFetcherConfig{
				headFetcher: mc,
				p2p:         p2p,
			})
		if err := fetcher.start(); err != nil {
			t.Error(err)
		}

		cancel()
		fetcher.stop()
	})

	dbtest.TeardownDB(t, beaconDB)
}

func TestBlocksFetcherRoundRobin(t *testing.T) {
	tests := []struct {
		name               string
		expectedBlockSlots []uint64
		peers              []*peerData
		requests           []*fetchRequestParams
	}{
		{
			name:               "Single peer with all blocks",
			expectedBlockSlots: makeSequence(1, 5*blockBatchSize+17),
			peers: []*peerData{
				{
					blocks:         makeSequence(1, 192),
					finalizedEpoch: 5,
					headSlot:       192,
				},
			},
			requests: []*fetchRequestParams{
				{
					start: 1,
					count: blockBatchSize,
				},
				{
					start: blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 2*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 3*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 4*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 5*blockBatchSize + 1,
					count: 17,
				},
			},
		},
		{
			name:               "Single peer with all blocks (many small requests)",
			expectedBlockSlots: makeSequence(1, 4*blockBatchSize/2+blockBatchSize/2),
			peers: []*peerData{
				{
					blocks:         makeSequence(1, 131),
					finalizedEpoch: 3,
					headSlot:       131,
				},
			},
			requests: []*fetchRequestParams{
				{
					start: 1,
					count: blockBatchSize / 2,
				},
				{
					start: blockBatchSize/2 + 1,
					count: blockBatchSize / 2,
				},
				{
					start: 2*blockBatchSize/2 + 1,
					count: blockBatchSize / 2,
				},
				{
					start: 3*blockBatchSize/2 + 1,
					count: blockBatchSize / 2,
				},
				{
					start: 4*blockBatchSize/2 + 1,
					count: blockBatchSize / 2,
				},
			},
		},
		{
			name:               "Multiple peers with all blocks",
			expectedBlockSlots: makeSequence(1, 171),
			peers: []*peerData{
				{
					blocks:         makeSequence(1, 171),
					finalizedEpoch: 5,
					headSlot:       171,
				},
				{
					blocks:         makeSequence(1, 171),
					finalizedEpoch: 5,
					headSlot:       171,
				},
				{
					blocks:         makeSequence(1, 171),
					finalizedEpoch: 5,
					headSlot:       171,
				},
			},
			requests: []*fetchRequestParams{
				{
					start: 1,
					count: blockBatchSize,
				},
				{
					start: blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 2*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 3*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 4*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 5*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 6*blockBatchSize + 1,
					count: blockBatchSize,
				},
			},
		},
		{
			name: "Multiple peers with skipped slots",
			// finalizedEpoch(18).slot = 608
			expectedBlockSlots: append(makeSequence(1, 64), makeSequence(500, 640)...), // up to 18th epoch
			peers: []*peerData{
				{
					blocks:         append(makeSequence(1, 64), makeSequence(500, 640)...),
					finalizedEpoch: 18,
					headSlot:       640,
				},
				{
					blocks:         append(makeSequence(1, 64), makeSequence(500, 640)...),
					finalizedEpoch: 18,
					headSlot:       640,
				},
				{
					blocks:         append(makeSequence(1, 64), makeSequence(500, 640)...),
					finalizedEpoch: 18,
					headSlot:       640,
				},
			},
			requests: []*fetchRequestParams{
				{
					start: 1,
					count: blockBatchSize,
				},
				{
					start: blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 2*blockBatchSize + 1,
					count: blockBatchSize,
				},
				{
					start: 400,
					count: 149,
				},
				{
					start: 550,
					count: 200,
				},
			},
		},
		{
			name:               "Multiple peers with failures",
			expectedBlockSlots: makeSequence(1, 2*blockBatchSize),
			peers: []*peerData{
				{
					blocks:         makeSequence(1, 320),
					finalizedEpoch: 8,
					headSlot:       320,
				},
				{
					blocks:         makeSequence(1, 320),
					finalizedEpoch: 8,
					headSlot:       320,
					failureSlots:   makeSequence(1, 32), // first epoch
				},
				{
					blocks:         makeSequence(1, 320),
					finalizedEpoch: 8,
					headSlot:       320,
				},
				{
					blocks:         makeSequence(1, 320),
					finalizedEpoch: 8,
					headSlot:       320,
					failureSlots:   makeSequence(1, 32), // first epoch
				},
			},
			requests: []*fetchRequestParams{
				{
					start: 1,
					count: blockBatchSize,
				},
				{
					start: blockBatchSize + 1,
					count: blockBatchSize,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.initializeRootCache(tt.expectedBlockSlots, t)

			beaconDB := dbtest.SetupDB(t)

			p := p2pt.NewTestP2P(t)
			connectPeers(t, p, tt.peers, p.Peers())
			cache.RLock()
			genesisRoot := cache.rootCache[0]
			cache.RUnlock()

			err := beaconDB.SaveBlock(context.Background(), &eth.SignedBeaconBlock{
				Block: &eth.BeaconBlock{
					Slot: 0,
				}})
			if err != nil {
				t.Fatal(err)
			}

			st, err := stateTrie.InitializeFromProto(&p2ppb.BeaconState{})
			if err != nil {
				t.Fatal(err)
			}
			mc := &mock.ChainService{
				State: st,
				Root:  genesisRoot[:],
				DB:    beaconDB,
			}

			ctx, cancel := context.WithCancel(context.Background())
			fetcher := newBlocksFetcher(
				ctx,
				&blocksFetcherConfig{
					headFetcher: mc,
					p2p:         p,
				})

			err = fetcher.start()
			if err != nil {
				t.Error(err)
			}

			var wg sync.WaitGroup
			wg.Add(len(tt.requests)) // how many block requests we are going to make
			go func() {
				wg.Wait()
				log.Debug("Stopping fetcher")
				fetcher.stop()
			}()

			processFetchedBlocks := func() ([]*eth.SignedBeaconBlock, error) {
				defer cancel()
				var unionRespBlocks []*eth.SignedBeaconBlock

				for {
					select {
					case resp, ok := <-fetcher.requestResponses():
						if !ok { // channel closed, aggregate
							return unionRespBlocks, nil
						}

						if resp.err != nil {
							log.WithError(resp.err).Debug("Block fetcher returned error")
							if resp.err.Error() == "request produced error, retry scheduled: bad" {
								continue
							}
						} else {
							unionRespBlocks = append(unionRespBlocks, resp.blocks...)
							if len(resp.blocks) == 0 {
								log.WithFields(logrus.Fields{
									"start": resp.start,
									"count": resp.count,
								}).Debug("Received empty slot")
							}
						}

						wg.Done()
					}
				}
			}

			maxExpectedBlocks := uint64(0)
			for _, req := range tt.requests {
				pid, err := fetcher.selectPeer(ctx)
				if err != nil {
					t.Error(err)
				}
				err = fetcher.scheduleRequest(context.Background(), pid, req.start, req.count)
				if err != nil {
					t.Error(err)
				}
				maxExpectedBlocks += req.count
			}

			blocks, err := processFetchedBlocks()
			if err != nil {
				t.Error(err)
			}

			sort.Slice(blocks, func(i, j int) bool {
				return blocks[i].Block.Slot < blocks[j].Block.Slot
			})

			slots := make([]uint64, len(blocks))
			for i, block := range blocks {
				slots[i] = block.Block.Slot
			}

			log.WithFields(logrus.Fields{
				"blocksLen": len(blocks),
				"slots":     slots,
			}).Debug("Finished block fetching")

			if len(blocks) > int(maxExpectedBlocks) {
				t.Errorf("Too many blocks returned. Wanted %d got %d", maxExpectedBlocks, len(blocks))
			}
			if len(blocks) != len(tt.expectedBlockSlots) {
				t.Errorf("Processes wrong number of blocks. Wanted %d got %d", len(tt.expectedBlockSlots), len(blocks))
			}
			var receivedBlockSlots []uint64
			for _, blk := range blocks {
				receivedBlockSlots = append(receivedBlockSlots, blk.Block.Slot)
			}
			if missing := sliceutil.NotUint64(sliceutil.IntersectionUint64(tt.expectedBlockSlots, receivedBlockSlots), tt.expectedBlockSlots); len(missing) > 0 {
				t.Errorf("Missing blocks at slots %v", missing)
			}

			dbtest.TeardownDB(t, beaconDB)
		})
	}
}

func TestBlocksFetcherScheduleRequest(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		fetcher := newBlocksFetcher(ctx, &blocksFetcherConfig{
			headFetcher: nil,
			p2p:         nil,
		})
		cancel()
		if err := fetcher.scheduleRequest(ctx, "", 1, blockBatchSize); err == nil {
			t.Errorf("expected error: %v", errFetcherCtxIsDone)
		}
	})
}

func TestBlocksFetcherHandleRequest(t *testing.T) {
	chainConfig := struct {
		expectedBlockSlots []uint64
		peers              []*peerData
	}{
		expectedBlockSlots: makeSequence(1, blockBatchSize),
		peers: []*peerData{
			{
				blocks:         makeSequence(1, 320),
				finalizedEpoch: 8,
				headSlot:       320,
			},
			{
				blocks:         makeSequence(1, 320),
				finalizedEpoch: 8,
				headSlot:       320,
			},
		},
	}

	mc, p2p, beaconDB := initializeTestServices(t, chainConfig.expectedBlockSlots, chainConfig.peers)
	defer dbtest.TeardownDB(t, beaconDB)

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		fetcher := newBlocksFetcher(ctx, &blocksFetcherConfig{
			headFetcher: mc,
			p2p:         p2p,
		})

		cancel()
		response := fetcher.handleRequest(ctx, &fetchRequestParams{})
		if response.err == nil {
			t.Errorf("expected error: %v", errFetcherCtxIsDone)
		}
	})

	t.Run("receive blocks", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		fetcher := newBlocksFetcher(ctx, &blocksFetcherConfig{
			headFetcher: mc,
			p2p:         p2p,
		})

		requestCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		go func() {
			pid, err := fetcher.selectPeer(ctx)
			if err != nil {
				t.Error(err)
			}
			request := &fetchRequestParams{
				ctx:   requestCtx,
				pid:   pid,
				start: 1,
				count: blockBatchSize,
			}
			response := fetcher.handleRequest(requestCtx, request)
			select {
			case <-ctx.Done():
			case fetcher.fetchResponses <- response:
			}
		}()

		var blocks []*eth.SignedBeaconBlock
		select {
		case <-ctx.Done():
			t.Error(ctx.Err())
		case resp := <-fetcher.requestResponses():
			if resp.err != nil {
				t.Error(resp.err)
			} else {
				blocks = resp.blocks
			}
		}
		if len(blocks) != blockBatchSize {
			t.Errorf("incorrect number of blocks returned, expected: %v, got: %v", blockBatchSize, len(blocks))
		}

		var receivedBlockSlots []uint64
		for _, blk := range blocks {
			receivedBlockSlots = append(receivedBlockSlots, blk.Block.Slot)
		}
		if missing := sliceutil.NotUint64(sliceutil.IntersectionUint64(chainConfig.expectedBlockSlots, receivedBlockSlots), chainConfig.expectedBlockSlots); len(missing) > 0 {
			t.Errorf("Missing blocks at slots %v", missing)
		}
	})
}

func TestBlocksFetcherRequestBlocks(t *testing.T) {
	chainConfig := struct {
		expectedBlockSlots []uint64
		peers              []*peerData
	}{
		expectedBlockSlots: makeSequence(1, 320),
		peers: []*peerData{
			{
				blocks:         makeSequence(1, 320),
				finalizedEpoch: 8,
				headSlot:       320,
			},
			{
				blocks:         makeSequence(1, 320),
				finalizedEpoch: 8,
				headSlot:       320,
			},
		},
	}

	hook := logTest.NewGlobal()
	mc, p2p, beaconDB := initializeTestServices(t, chainConfig.expectedBlockSlots, chainConfig.peers)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fetcher := newBlocksFetcher(
		ctx,
		&blocksFetcherConfig{
			headFetcher: mc,
			p2p:         p2p,
		})

	root, _, peers := p2p.Peers().BestFinalized(params.BeaconConfig().MaxPeersToSync, helpers.SlotToEpoch(mc.HeadSlot()))

	blocks, err := fetcher.requestBlocks(context.Background(), peers[0], root, 1, blockBatchSize)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(blocks) != blockBatchSize {
		t.Errorf("incorrect number of blocks returned, expected: %v, got: %v", blockBatchSize, len(blocks))
	}

	// Test request fail over (success).
	err = fetcher.p2p.Disconnect(peers[0])
	if err != nil {
		t.Error(err)
	}
	hook.Reset()
	request := &fetchRequestParams{
		ctx:   ctx,
		pid:   peers[0],
		start: 1,
		count: blockBatchSize,
	}
	response := fetcher.handleRequest(request.ctx, request)
	testutil.AssertLogsContain(t, hook, "Retrying request")
	expectedErr := fmt.Sprintf("failed to dial %v: no addresses", peers[0])
	if response.err != nil && !strings.Contains(response.err.Error(), expectedErr) {
		t.Errorf("unexpected error: want: %v, got: %v", expectedErr, response.err)
	}

	// Test context cancellation.
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	blocks, err = fetcher.requestBlocks(ctx, peers[0], root, 1, blockBatchSize)
	if err == nil || err.Error() != "context canceled" {
		t.Errorf("expected context closed error, got: %v", err)
	}

	dbtest.TeardownDB(t, beaconDB)
}

func initializeTestServices(t *testing.T, blocks []uint64, peers []*peerData) (*mock.ChainService, *p2pt.TestP2P, db.Database) {
	cache.initializeRootCache(blocks, t)
	beaconDB := dbtest.SetupDB(t)

	p := p2pt.NewTestP2P(t)
	connectPeers(t, p, peers, p.Peers())
	cache.RLock()
	genesisRoot := cache.rootCache[0]
	cache.RUnlock()

	err := beaconDB.SaveBlock(context.Background(), &eth.SignedBeaconBlock{
		Block: &eth.BeaconBlock{
			Slot: 0,
		}})
	if err != nil {
		t.Fatal(err)
	}

	st, err := stateTrie.InitializeFromProto(&p2ppb.BeaconState{})
	if err != nil {
		t.Fatal(err)
	}

	return &mock.ChainService{
		State: st,
		Root:  genesisRoot[:],
		DB:    beaconDB,
	}, p, beaconDB
}
