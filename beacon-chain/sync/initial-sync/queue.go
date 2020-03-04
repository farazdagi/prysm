package initialsync

import (
	"context"
	"fmt"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/blockchain"
)

const fetchRequestsLimit = 8
const fetchedBlocksBufferSize = 0

// blocksProvider exposes enough methods for queue to schedule and fetch blocks.
type blocksProvider interface {
	requestResponses() <-chan *fetchRequestResponse
	scheduleRequest(ctx context.Context, start, count uint64)
}

// blocksQueueConfig is a config to setup block queue service.
type blocksQueueConfig struct {
	ctx                 context.Context
	blocksFetcher       blocksProvider
	headFetcher         blockchain.HeadFetcher
	startSlot           uint64
	highestExpectedSlot uint64
}

// blocksQueue is a priority queue that serves as a intermediary between block fetchers (producers)
// and block processing goroutine (consumer). Consumer can rely on order of incoming blocks.
type blocksQueue struct {
	ctx                 context.Context
	blocksFetcher       blocksProvider
	headFetcher         blockchain.HeadFetcher
	startSlot           uint64                      // starting slot from which queue will start of
	highestExpectedSlot uint64                      // highest slot to which queue will work towards
	fetchedBlocks       chan *eth.SignedBeaconBlock // ordered blocks (up to highest expected slot) are pushed into this channel
	fetchRequests       chan *fetchRequestParams    // controls fetch request scheduling
	quit                chan struct{}               // termination notifier
}

// newBlocksQueue creates initialized priority queue.
func newBlocksQueue(cfg *blocksQueueConfig) *blocksQueue {
	ctx := cfg.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	startSlot := cfg.startSlot
	if startSlot == 0 && cfg.headFetcher != nil {
		startSlot = cfg.headFetcher.HeadSlot() + 1
	}

	return &blocksQueue{
		ctx:                 ctx,
		blocksFetcher:       cfg.blocksFetcher,
		headFetcher:         cfg.headFetcher,
		startSlot:           startSlot,
		highestExpectedSlot: cfg.highestExpectedSlot,
		fetchedBlocks:       make(chan *eth.SignedBeaconBlock, fetchedBlocksBufferSize),
		fetchRequests:       make(chan *fetchRequestParams, fetchRequestsLimit),
		quit:                make(chan struct{}),
	}
}

func (q *blocksQueue) startFetchRequestLoop() {
	fmt.Printf("here we go")
}
