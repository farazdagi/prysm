package initialsync

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/flags"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	prysmsync "github.com/prysmaticlabs/prysm/beacon-chain/sync"
	p2ppb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

const fetchRequestRetries = 3 // how many times to retry fetching given range

var (
	errNoPeersAvailable   = errors.New("no peers available, waiting for reconnect")
	errFetcherCtxIsDone   = errors.New("fetcher's context is done, reinitialize")
	errStartSlotIsTooHigh = errors.New("start slot is bigger than highest finalized slot")
)

// blocksFetcherConfig is a config to setup the block fetcher.
type blocksFetcherConfig struct {
	headFetcher blockchain.HeadFetcher
	p2p         p2p.P2P
}

// blocksFetcher is a service to fetch chain data from peers.
// On an incoming requests, requested block range is evenly divided
// among available peers (for fair network load distribution).
type blocksFetcher struct {
	ctx            context.Context
	cancel         context.CancelFunc
	headFetcher    blockchain.HeadFetcher
	p2p            p2p.P2P
	rateLimiter    *leakybucket.Collector
	fetchRequests  chan *fetchRequestParams
	fetchResponses chan *fetchRequestResponse
	quit           chan struct{} // termination notifier
}

// fetchRequestParams holds parameters necessary to schedule a fetch request.
type fetchRequestParams struct {
	ctx     context.Context // if provided, it is used instead of global fetcher's context
	pid     peer.ID         // assigned peer
	start   uint64          // starting slot
	count   uint64          // how many slots to receive (fetcher may return fewer slots)
	retries uint8           // number of retries
}

// fetchRequestResponse is a combined type to hold results of both successful executions and errors.
// Valid usage pattern will be to check whether result's `err` is nil, before using `blocks`.
type fetchRequestResponse struct {
	start, count uint64
	blocks       []*eth.SignedBeaconBlock
	err          error
}

// newBlocksFetcher creates ready to use fetcher.
func newBlocksFetcher(ctx context.Context, cfg *blocksFetcherConfig) *blocksFetcher {
	ctx, cancel := context.WithCancel(ctx)
	rateLimiter := leakybucket.NewCollector(
		allowedBlocksPerSecond, /* rate */
		allowedBlocksPerSecond, /* capacity */
		false                   /* deleteEmptyBuckets */)

	return &blocksFetcher{
		ctx:            ctx,
		cancel:         cancel,
		headFetcher:    cfg.headFetcher,
		p2p:            cfg.p2p,
		rateLimiter:    rateLimiter,
		fetchRequests:  make(chan *fetchRequestParams, queueMaxPendingRequests),
		fetchResponses: make(chan *fetchRequestResponse, queueMaxPendingRequests),
		quit:           make(chan struct{}),
	}
}

// start boots up the fetcher, which starts listening for incoming fetch requests.
func (f *blocksFetcher) start() error {
	select {
	case <-f.ctx.Done():
		return errFetcherCtxIsDone
	default:
		go f.loop()
		return nil
	}
}

// stop terminates all fetcher operations.
func (f *blocksFetcher) stop() {
	f.cancel()
	<-f.quit // make sure that loop() is done
}

// requestResponses exposes a channel into which fetcher pushes generated request responses.
func (f *blocksFetcher) requestResponses() <-chan *fetchRequestResponse {
	return f.fetchResponses
}

// loop is a main fetcher loop, listens for incoming requests/cancellations, forwards outgoing responses.
func (f *blocksFetcher) loop() {
	defer close(f.quit)

	// Wait for all loop's goroutines to finish, and safely release resources.
	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
		close(f.fetchResponses)
	}()

	for {
		// Make sure there is are available peers before processing requests.
		if _, err := f.waitForMinimumPeers(f.ctx); err != nil {
			log.Error(err)
		}

		select {
		case <-f.ctx.Done():
			log.Debug("Context closed, exiting goroutine (blocks fetcher)")
			return
		case req := <-f.fetchRequests:
			// Do not overflow any peers, slowdown on a recently used peer.
			if f.rateLimiter.Remaining(req.pid.String()) < int64(blockBatchSize) {
				log.WithField("peer", req.pid).Debug("Slowing down before processing request")
				time.Sleep(f.rateLimiter.TillEmpty(req.pid.String()))
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-f.ctx.Done():
				case f.fetchResponses <- f.handleRequest(req.ctx, req):
				}
			}()
		}
	}
}

// scheduleRequest adds request to incoming queue.
func (f *blocksFetcher) scheduleRequest(ctx context.Context, pid peer.ID, start, count uint64) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if pid == "" {
		pid, err = f.selectPeer(ctx)
		if err != nil {
			return
		}
	}

	request := &fetchRequestParams{
		ctx:   ctx,
		pid:   pid,
		start: start,
		count: count,
	}
	select {
	case <-f.ctx.Done():
		return errFetcherCtxIsDone
	case f.fetchRequests <- request:
	}
	return nil
}

// rescheduleRequest tries to reschedule previously scheduled request.
func (f *blocksFetcher) rescheduleRequest(ctx context.Context, req *fetchRequestParams) error {
	pid, err := f.selectPeer(ctx)
	if err != nil {
		return errors.Wrap(err, "error re-scheduling request")
	}
	req.retries++
	return f.scheduleRequest(ctx, pid, req.start, req.count)
}

// handleRequest parses fetch request and forwards it to response builder.
func (f *blocksFetcher) handleRequest(ctx context.Context, req *fetchRequestParams) *fetchRequestResponse {
	ctx, span := trace.StartSpan(ctx, "initialsync.handleRequest")
	defer span.End()

	response := &fetchRequestResponse{
		start:  req.start,
		count:  req.count,
		blocks: []*eth.SignedBeaconBlock{},
	}

	if ctx.Err() != nil {
		response.err = ctx.Err()
		return response
	}

	headEpoch := helpers.SlotToEpoch(f.headFetcher.HeadSlot())
	root, finalizedEpoch, peers := f.p2p.Peers().BestFinalized(params.BeaconConfig().MaxPeersToSync, headEpoch)

	if len(peers) == 0 {
		response.err = errNoPeersAvailable
		return response
	}

	// Short circuit start far exceeding the highest finalized epoch in some infinite loop.
	highestFinalizedSlot := helpers.StartSlot(finalizedEpoch + 1)
	if req.start > highestFinalizedSlot {
		response.err = errStartSlotIsTooHigh
		return response
	}

	blocks, err := f.requestBlocks(ctx, req.pid, root, req.start, req.count)
	if err != nil {
		// When there are no retries left, just send error response.
		if req.retries >= fetchRequestRetries {
			response.err = err
			return response
		}

		// Try to resend request (to another peer).
		if err := f.rescheduleRequest(ctx, req); err != nil {
			response.err = err
			return response
		}

		response.err = errors.Wrap(err, "request produced error, retry scheduled")
		log.WithError(response.err).Debug("Retrying request")
		return response
	}

	response.blocks = blocks
	return response
}

// requestBlocks is a wrapper for handling BeaconBlocksByRangeRequest requests/streams.
func (f *blocksFetcher) requestBlocks(
	ctx context.Context,
	pid peer.ID,
	root []byte,
	start, count uint64,
) ([]*eth.SignedBeaconBlock, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if f.rateLimiter.Remaining(pid.String()) < int64(count) {
		log.WithField("peer", pid).Debug("Slowing down for rate limit")
		time.Sleep(f.rateLimiter.TillEmpty(pid.String()))
	}
	f.rateLimiter.Add(pid.String(), int64(count))
	log.WithFields(logrus.Fields{
		"peer":  pid,
		"range": []uint64{start, start + count - 1},
		"count": count,
		"head":  fmt.Sprintf("%#x", root),
	}).Debug("Requesting blocks")

	// Send request.
	req := &p2ppb.BeaconBlocksByRangeRequest{
		HeadBlockRoot: root,
		StartSlot:     start,
		Count:         count,
		Step:          1,
	}
	stream, err := f.p2p.Send(ctx, req, pid)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	resp := make([]*eth.SignedBeaconBlock, 0, req.Count)
	for {
		blk, err := prysmsync.ReadChunkedBlock(stream, f.p2p)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		resp = append(resp, blk)
	}

	return resp, nil
}

// selectPeer spins up until there is a minimum required number of available peers.
// At that point randomly selects a peer that is ready to respond.
func (f *blocksFetcher) selectPeer(ctx context.Context) (peer.ID, error) {
	peers, err := f.waitForMinimumPeers(ctx)
	if err != nil {
		return "", err
	}

	if len(peers) == 0 {
		return "", errNoPeersAvailable
	}

	// Randomize list of peers.
	randGenerator := rand.New(rand.NewSource(time.Now().Unix()))
	randGenerator.Shuffle(len(peers), func(i, j int) {
		peers[i], peers[j] = peers[j], peers[i]
	})

	return peers[0], nil
}

// waitForMinimumPeers spins and waits up until enough peers are available.
func (f *blocksFetcher) waitForMinimumPeers(ctx context.Context) ([]peer.ID, error) {
	required := params.BeaconConfig().MaxPeersToSync
	if flags.Get().MinimumSyncPeers < required {
		required = flags.Get().MinimumSyncPeers
	}
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		headEpoch := helpers.SlotToEpoch(f.headFetcher.HeadSlot())
		_, _, peers := f.p2p.Peers().BestFinalized(params.BeaconConfig().MaxPeersToSync, headEpoch)
		if len(peers) >= required {
			return peers, nil
		}
		log.WithFields(logrus.Fields{
			"suitable": len(peers),
			"required": required}).Info("Waiting for enough suitable peers before syncing")
		time.Sleep(handshakePollingInterval)
	}
}
