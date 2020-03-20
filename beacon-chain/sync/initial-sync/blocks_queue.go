package initialsync

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	"github.com/prysmaticlabs/prysm/shared/mathutil"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

const (
	// queueMaxPendingRequests limits how many concurrent fetch request queue can initiate.
	queueMaxPendingRequests = 8
	// queueMaxCachedBlocks hard limit on how many queue items to cache before forced dequeue.
	queueMaxCachedBlocks = 16 * queueMaxPendingRequests * blockBatchSize
	// queueStopCallTimeout is time allowed for queue to release resources when quitting.
	queueStopCallTimeout = 1 * time.Second
	// queuePollingInterval is an interval at which queue schedules fetch request.
	queuePollingInterval = 1 * time.Second
	// backoffInterval is an interval for which system waits before proceeding with another request.
	backoffInterval = 4 * time.Second
)

var (
	errQueueCtxIsDone          = errors.New("queue's context is done, reinitialize")
	errQueueTakesTooLongToStop = errors.New("queue takes too long to stop")
)

// blocksProvider exposes enough methods for queue to fetch incoming blocks.
type blocksProvider interface {
	waitForMinimumPeers(ctx context.Context) ([]peer.ID, error)
	requestResponses() <-chan *fetchRequestResponse
	scheduleRequest(ctx context.Context, start, count uint64) error
	start() error
	stop()
}

// blocksQueueConfig is a config to setup block queue service.
type blocksQueueConfig struct {
	blocksFetcher       blocksProvider
	headFetcher         blockchain.HeadFetcher
	startSlot           uint64
	highestExpectedSlot uint64
	p2p                 p2p.P2P
}

// blocksQueueState holds internal queue state (for easier management of state transitions).
type blocksQueueState struct {
	scheduler    *schedulerState
	sender       *senderState
	cachedBlocks map[uint64]*cachedBlock
}

// blockState enums possible queue block states.
type blockState uint8

const (
	// pendingBlock is a default block status when just added to queue.
	pendingBlock = iota
	// readyBlock represents a block for which response has been received from the fetcher.
	readyBlock
	// validBlock represents block that can be processed.
	validBlock
	// skippedBlock is a block for a slot that is not found on any available peers.
	skippedBlock
	// failedBlock represents block that can not be processed at the moment.
	failedBlock
	// blockStateLen is a sentinel to know number of possible block states.
	blockStateLen
)

// schedulerState a state of scheduling process.
type schedulerState struct {
	sync.Mutex
	currentSlot     uint64
	requestedBlocks map[blockState]uint64
}

// senderState is a state of block sending process.
type senderState struct {
	sync.Mutex
}

// cachedBlock is a container for signed beacon block.
type cachedBlock struct {
	*eth.SignedBeaconBlock
}

// blocksQueue is a priority queue that serves as a intermediary between block fetchers (producers)
// and block processing goroutine (consumer). Consumer can rely on order of incoming blocks.
type blocksQueue struct {
	ctx                  context.Context
	cancel               context.CancelFunc
	highestExpectedSlot  uint64
	state                *blocksQueueState
	blocksFetcher        blocksProvider
	headFetcher          blockchain.HeadFetcher
	fetchedBlocks        chan *eth.SignedBeaconBlock // output channel for ready blocks
	pendingFetchRequests chan struct{}               // pending requests semaphore
	quit                 chan struct{}               // termination notifier
}

// newBlocksQueue creates initialized priority queue.
func newBlocksQueue(ctx context.Context, cfg *blocksQueueConfig) *blocksQueue {
	ctx, cancel := context.WithCancel(ctx)

	blocksFetcher := cfg.blocksFetcher
	if blocksFetcher == nil {
		blocksFetcher = newBlocksFetcher(ctx, &blocksFetcherConfig{
			headFetcher: cfg.headFetcher,
			p2p:         cfg.p2p,
		})
	}

	return &blocksQueue{
		ctx:                 ctx,
		cancel:              cancel,
		highestExpectedSlot: cfg.highestExpectedSlot,
		state: &blocksQueueState{
			scheduler: &schedulerState{
				currentSlot:     cfg.startSlot,
				requestedBlocks: make(map[blockState]uint64, blockStateLen),
			},
			sender:       &senderState{},
			cachedBlocks: make(map[uint64]*cachedBlock, queueMaxCachedBlocks),
		},
		blocksFetcher:        blocksFetcher,
		headFetcher:          cfg.headFetcher,
		fetchedBlocks:        make(chan *eth.SignedBeaconBlock, blockBatchSize),
		pendingFetchRequests: make(chan struct{}, queueMaxPendingRequests),
		quit:                 make(chan struct{}),
	}
}

// start boots up the queue processing.
func (q *blocksQueue) start() error {
	select {
	case <-q.ctx.Done():
		return errQueueCtxIsDone
	default:
		go q.loop()
		return nil
	}
}

// stop terminates all queue operations.
func (q *blocksQueue) stop() error {
	q.cancel()
	select {
	case <-q.quit:
		return nil
	case <-time.After(queueStopCallTimeout):
		return errQueueTakesTooLongToStop
	}
}

// loop is a main queue loop.
func (q *blocksQueue) loop() {
	defer close(q.quit)

	// Wait for all goroutines to wrap up (forced by cancelled context), and do a cleanup.
	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
		q.blocksFetcher.stop()
		close(q.fetchedBlocks)
	}()

	if err := q.blocksFetcher.start(); err != nil {
		log.WithError(err).Debug("Can not start blocks provider")
	}

	// Reads from semaphore channel, thus allowing next goroutine to grab it and schedule next request.
	releaseTicket := func() {
		select {
		case <-q.ctx.Done():
		case <-q.pendingFetchRequests:
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		q.schedulerLoop(q.ctx)
	}()

	ticker := time.NewTicker(queuePollingInterval)
	defer ticker.Stop()
	for {
		if q.headFetcher.HeadSlot() >= q.highestExpectedSlot {
			log.Debug("Highest expected slot reached")
			q.cancel()
		}

		// Make sure there is are available peers before processing requests.
		if _, err := q.blocksFetcher.waitForMinimumPeers(q.ctx); err != nil {
			log.Error(err)
		}

		select {
		case <-q.ctx.Done():
			log.Debug("Context closed, exiting goroutine (blocks queue)")
			return
		case q.pendingFetchRequests <- struct{}{}:
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Schedule request.
				if err := q.scheduleFetchRequests(q.ctx); err != nil {
					log.Debugf("Error scheduling fetch requests: %v", err)
					time.Sleep(backoffInterval)
					releaseTicket()
				}
			}()
		case response, ok := <-q.blocksFetcher.requestResponses():
			if !ok {
				log.Debug("Fetcher closed output channel")
				q.cancel()
				return
			}
			log.WithFields(logrus.Fields{
				"err":   response.err,
				"start": response.start,
				"count": response.count,
				"len":   len(response.blocks),
			}).Debug("Response received....")

			// Release semaphore ticket.
			go releaseTicket()

			// Process incoming response into blocks.
			wg.Add(1)
			go func() {
				defer wg.Done()

				skippedBlocks, err := q.parseFetchResponse(q.ctx, response)
				if err != nil {
					log.Debugf("Error response received from fetcher: %v", err)
					q.state.scheduler.incrementCounter(failedBlock, response.count)
					return
				}
				q.state.scheduler.incrementCounter(skippedBlock, skippedBlocks)
				log.WithFields(logrus.Fields{
					"skipped": skippedBlocks,
					"start":   response.start,
					"count":   response.count,
					"len":     len(response.blocks),
				}).Debug("Response received")
			}()
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := q.sendFetchedBlocks(q.ctx); err != nil {
					log.Debugf("Error sending received blocks: %v", err)
				}
			}()
		}
	}
}

func (q *blocksQueue) schedulerLoop(ctx context.Context) {
	releaseTicket := func() {
		select {
		case <-q.ctx.Done():
		case <-q.pendingFetchRequests:
		}
	}

	for {
		select {
		case <-q.ctx.Done():
			log.Debug("Context closed, exiting goroutine (blocks queue)")
			return
		case q.pendingFetchRequests <- struct{}{}:
			// Schedule request.
			if err := q.scheduleFetchRequests(q.ctx); err != nil {
				log.Debugf("Error scheduling fetch requests: %v", err)
				time.Sleep(backoffInterval)
				releaseTicket()
			}
		}
	}

}

// scheduleFetchRequests enqueues block fetch requests to block fetcher.
func (q *blocksQueue) scheduleFetchRequests(ctx context.Context) error {
	q.state.scheduler.Lock()
	defer q.state.scheduler.Unlock()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	s := q.state.scheduler
	blocks := q.state.scheduler.requestedBlocks
	headSlot := q.headFetcher.HeadSlot()

	resetStateCounters := func() {
		for i := 0; i < blockStateLen; i++ {
			blocks[blockState(i)] = 0
		}
		s.currentSlot = headSlot
	}

	// Update state's current slot pointer.
	count := blocks[pendingBlock] + blocks[skippedBlock] + blocks[failedBlock] + blocks[validBlock]
	if count == 0 {
		s.currentSlot = headSlot
	}

	// All blocks processed, no pending blocks.
	count = blocks[skippedBlock] + blocks[failedBlock] + blocks[validBlock]
	if count > 0 && blocks[pendingBlock] == 0 && count != blocks[skippedBlock] {
		resetStateCounters()
	}

	// Too many pending blocks, slowdown.
	if count = blocks[pendingBlock]; count > queueMaxPendingRequests*blockBatchSize {
		log.WithField("blocks", blocks).Debug("Too many pending blocks")
		return errors.New("too many pending blocks, slowing down")
	}

	// Too many items in scheduler's queue, adding more is not allowed.
	count = blocks[pendingBlock] + blocks[failedBlock] + blocks[validBlock]
	log.WithField("blocks", blocks).Debug("-------------->>>")
	if count >= queueMaxCachedBlocks {
		resetStateCounters()
		log.WithField("blocks", blocks).Debug("AFTER-------------->>>")
	}

	offset := blocks[pendingBlock] + blocks[skippedBlock] + blocks[failedBlock] + blocks[validBlock]
	start := q.state.scheduler.currentSlot + offset + 1
	count = mathutil.Min(blockBatchSize, q.highestExpectedSlot-start+1)

	log.WithFields(logrus.Fields{
		"blocks": blocks,
		"start":  start,
		"count":  count,
		"offset": offset,
	}).Debug("<<<<<<<<<<<<<<<<<<<<<<<<before sced")
	if count <= 0 {
		return errStartSlotIsTooHigh
	}

	if err := q.blocksFetcher.scheduleRequest(ctx, start, count); err != nil {
		return err
	}
	q.state.scheduler.requestedBlocks[pendingBlock] += count

	return nil
}

// parseFetchResponse processes incoming responses.
func (q *blocksQueue) parseFetchResponse(ctx context.Context, response *fetchRequestResponse) (uint64, error) {
	q.state.sender.Lock()
	defer q.state.sender.Unlock()

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	if response.err != nil {
		return 0, response.err
	}

	if q.headFetcher.HeadSlot() > response.start {
		return 0, errors.New("outdated blocks")
	}

	// Extract beacon blocks.
	responseBlocks := make(map[uint64]*eth.SignedBeaconBlock, len(response.blocks))
	for _, blk := range response.blocks {
		responseBlocks[blk.Block.Slot] = blk
	}

	// Cache blocks in [start, start + count) range, include skipped blocks.
	var skippedBlocks uint64
	end := response.start + mathutil.Max(response.count, uint64(len(response.blocks)))
	for slot := response.start; slot < end; slot++ {
		if block, ok := responseBlocks[slot]; ok {
			q.state.cachedBlocks[slot] = &cachedBlock{
				SignedBeaconBlock: block,
			}
			delete(responseBlocks, slot)
			continue
		}
		q.state.cachedBlocks[slot] = &cachedBlock{}
		skippedBlocks++
	}

	// If there are any items left in incoming response, cache them too.
	for slot, block := range responseBlocks {
		q.state.cachedBlocks[slot] = &cachedBlock{
			SignedBeaconBlock: block,
		}
	}

	return skippedBlocks, nil
}

// sendFetchedBlocks analyses available blocks, and sends them downstream in a correct slot order.
// Blocks are checked starting from the current head slot, and up until next consecutive block is available.
func (q *blocksQueue) sendFetchedBlocks(ctx context.Context) error {
	q.state.sender.Lock()
	defer q.state.sender.Unlock()

	ctx, span := trace.StartSpan(ctx, "initialsync.sendFetchedBlocks")
	defer span.End()

	startSlot := q.headFetcher.HeadSlot() + 1
	for slot := startSlot; slot <= q.highestExpectedSlot; slot++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		blockData, ok := q.state.cachedBlocks[slot]
		if !ok {
			break
		}
		if blockData.SignedBeaconBlock != nil && blockData.Block != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case q.fetchedBlocks <- blockData.SignedBeaconBlock:
			}
		}
	}

	// Remove processed blocks.
	for slot := range q.state.cachedBlocks {
		if slot < startSlot {
			delete(q.state.cachedBlocks, slot)
		}
	}
	log.WithField("LEN", len(q.state.cachedBlocks)).Debug("QUEUE CACHE SIZE")

	return nil
}

// deleteCachedBlock removes cached block.
func (q *blocksQueue) deleteCachedBlock(key uint64) {
	q.state.sender.Lock()
	defer q.state.sender.Unlock()
	delete(q.state.cachedBlocks, key)
}

// incrementCounter increments particular scheduler counter.
func (s *schedulerState) incrementCounter(counter blockState, n uint64) {
	s.Lock()
	defer s.Unlock()

	// Assert that counter is within acceptable boundaries.
	if counter < 1 || counter >= blockStateLen {
		return
	}

	n = mathutil.Min(s.requestedBlocks[pendingBlock], n)
	s.requestedBlocks[counter] += n
	s.requestedBlocks[pendingBlock] -= n
}

// updateCounters moves n tokens from source to target counters.
func (s *schedulerState) updateCounters(source, target blockState, n uint64) uint64 {
	s.Lock()
	defer s.Unlock()

	n = mathutil.Min(s.requestedBlocks[source], n)
	s.requestedBlocks[source] -= n
	s.requestedBlocks[target] += n
	return n
}
