package helpers

import (
	"reflect"
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/cache"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEpochCommitteeCount_OK(t *testing.T) {
	// this defines the # of validators required to have 1 committee
	// per slot for epoch length.
	validatorsPerEpoch := params.BeaconConfig().SlotsPerEpoch * params.BeaconConfig().TargetCommitteeSize
	tests := []struct {
		validatorCount uint64
		committeeCount uint64
	}{
		{0, params.BeaconConfig().SlotsPerEpoch},
		{1000, params.BeaconConfig().SlotsPerEpoch},
		{2 * validatorsPerEpoch, 2 * params.BeaconConfig().SlotsPerEpoch},
		{5 * validatorsPerEpoch, 5 * params.BeaconConfig().SlotsPerEpoch},
		{16 * validatorsPerEpoch, 16 * params.BeaconConfig().SlotsPerEpoch},
		{32 * validatorsPerEpoch, 16 * params.BeaconConfig().SlotsPerEpoch},
	}
	for _, test := range tests {
		vals := make([]*pb.Validator, test.validatorCount)
		for i := 0; i < len(vals); i++ {
			vals[i] = &pb.Validator{
				ExitEpoch: params.BeaconConfig().FarFutureEpoch,
			}
		}
		s := &pb.BeaconState{
			ValidatorRegistry: vals,
		}
		if test.committeeCount != EpochCommitteeCount(s, 1) {
			t.Errorf("wanted: %d, got: %d",
				test.committeeCount, EpochCommitteeCount(s, 1))
		}
	}
}

func TestEpochCommitteeCount_LessShardsThanEpoch(t *testing.T) {
	validatorCount := uint64(8)
	productionConfig := params.BeaconConfig()
	testConfig := &params.BeaconChainConfig{
		ShardCount:          1,
		SlotsPerEpoch:       4,
		TargetCommitteeSize: 2,
	}
	params.OverrideBeaconConfig(testConfig)
	vals := make([]*pb.Validator, validatorCount)
	for i := 0; i < len(vals); i++ {
		vals[i] = &pb.Validator{
			ExitEpoch: params.BeaconConfig().FarFutureEpoch,
		}
	}
	s := &pb.BeaconState{
		ValidatorRegistry: vals,
	}
	if EpochCommitteeCount(s, 1) != validatorCount/testConfig.TargetCommitteeSize {
		t.Errorf("wanted: %d, got: %d",
			validatorCount/testConfig.TargetCommitteeSize, EpochCommitteeCount(s, 1))
	}
	params.OverrideBeaconConfig(productionConfig)
}

func TestShardDelta_OK(t *testing.T) {
	minShardDelta := params.BeaconConfig().ShardCount -
		params.BeaconConfig().ShardCount/params.BeaconConfig().SlotsPerEpoch
	tests := []struct {
		validatorCount uint64
		shardCount     uint64
	}{
		{0, params.BeaconConfig().SlotsPerEpoch},    // Empty minimum shards
		{1000, params.BeaconConfig().SlotsPerEpoch}, // 1000 Validators minimum shards,
		{100000, 768 /*len(active_validators) // TARGET_COMMITTEE_SIZE*/},
		{500000, minShardDelta}, // 5 Mil, above shard delta
	}
	for _, test := range tests {
		vals := make([]*pb.Validator, test.validatorCount)
		for i := 0; i < len(vals); i++ {
			vals[i] = &pb.Validator{
				ExitEpoch: params.BeaconConfig().FarFutureEpoch,
			}
		}
		s := &pb.BeaconState{
			ValidatorRegistry: vals,
		}
		if test.shardCount != ShardDelta(s, 1) {
			t.Errorf("wanted: %d, got: %d",
				test.shardCount, ShardDelta(s, 1))
		}
	}
}

func TestComputeCommittee_OK(t *testing.T) {
	// TODO(2682): Don't fix this test, this will be removed after merging #2682
	t.Skip()

	validatorsPerEpoch := params.BeaconConfig().SlotsPerEpoch * params.BeaconConfig().TargetCommitteeSize
	committeesPerEpoch := uint64(6)
	// Set epoch total validators count to 6 committees per slot.
	validators := make([]*pb.Validator, committeesPerEpoch*validatorsPerEpoch)
	for i := 0; i < len(validators); i++ {
		validators[i] = &pb.Validator{
			ExitEpoch: params.BeaconConfig().FarFutureEpoch,
		}
	}

	state := &pb.BeaconState{
		ValidatorRegistry:      validators,
		Slot:                   params.BeaconConfig().GenesisSlot + 200,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}

	wantedEpoch := SlotToEpoch(state.Slot)
	shardCount := params.BeaconConfig().ShardCount
	startShard := state.LatestStartShard

	committeesPerSlot := committeesPerEpoch / params.BeaconConfig().SlotsPerEpoch
	offset := state.Slot % params.BeaconConfig().SlotsPerEpoch
	slotStartShard := (startShard + committeesPerSlot + offset) % shardCount
	seed := GenerateSeed(state, wantedEpoch)

	indices := ActiveValidatorIndices(state, wantedEpoch)
	newCommittees := make([]*CrosslinkCommittee, committeesPerSlot)
	committees := []*CrosslinkCommittee{{}}
	for i := uint64(0); i < committeesPerSlot; i++ {
		committee, err := ComputeCommittee(indices, seed, committeesPerSlot*offset+i, committeesPerEpoch)
		if err != nil {
			t.Errorf("could not compute committee: %v", err)
		}
		committees[i] = &CrosslinkCommittee{
			Committee: committee,
			Shard:     (slotStartShard + i) % shardCount,
		}
	}

	if reflect.DeepEqual(committees, newCommittees) {
		t.Error("Committees from different slot shall not be equal")
	}
}

func TestAttestationParticipants_NoCommitteeCache(t *testing.T) {
	if params.BeaconConfig().SlotsPerEpoch != 64 {
		t.Errorf("SlotsPerEpoch should be 64 for these tests to pass")
	}

	validators := make([]*pb.Validator, 2*params.BeaconConfig().SlotsPerEpoch)
	for i := 0; i < len(validators); i++ {
		validators[i] = &pb.Validator{
			ExitEpoch: params.BeaconConfig().FarFutureEpoch,
		}
	}

	state := &pb.BeaconState{
		ValidatorRegistry:      validators,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}

	attestationData := &pb.AttestationData{}

	tests := []struct {
		attestationSlot uint64
		stateSlot       uint64
		shard           uint64
		bitfield        []byte
		wanted          []uint64
	}{
		{
			attestationSlot: params.BeaconConfig().GenesisSlot + 2,
			stateSlot:       params.BeaconConfig().GenesisSlot + 5,
			shard:           3,
			bitfield:        []byte{0x03},
			wanted:          []uint64{37, 100},
		},
		{
			attestationSlot: params.BeaconConfig().GenesisSlot + 1,
			stateSlot:       params.BeaconConfig().GenesisSlot + 10,
			shard:           2,
			bitfield:        []byte{0x01},
			wanted:          []uint64{2, 35},
		},
		{
			attestationSlot: params.BeaconConfig().GenesisSlot + 10,
			stateSlot:       params.BeaconConfig().GenesisSlot + 10,
			shard:           11,
			bitfield:        []byte{0x03},
			wanted:          []uint64{95, 101},
		},
	}

	for _, tt := range tests {
		state.Slot = tt.stateSlot
		attestationData.Slot = tt.attestationSlot
		attestationData.Shard = tt.shard
		attestationData.TargetEpoch = params.BeaconConfig().GenesisEpoch

		result, err := AttestingIndices(state, attestationData, tt.bitfield)
		if err != nil {
			t.Errorf("Failed to get attestation participants: %v", err)
		}

		if !reflect.DeepEqual(tt.wanted, result) {
			t.Errorf(
				"Result indices was an unexpected value. Wanted %d, got %d",
				tt.wanted,
				result,
			)
		}
	}
}

func TestAttestationParticipants_IncorrectBitfield(t *testing.T) {
	if params.BeaconConfig().SlotsPerEpoch != 64 {
		t.Errorf("SlotsPerEpoch should be 64 for these tests to pass")
	}

	validators := make([]*pb.Validator, params.BeaconConfig().DepositsForChainStart)
	for i := 0; i < len(validators); i++ {
		validators[i] = &pb.Validator{
			ExitEpoch: params.BeaconConfig().FarFutureEpoch,
		}
	}

	state := &pb.BeaconState{
		ValidatorRegistry:      validators,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}
	attestationData := &pb.AttestationData{}

	if _, err := AttestingIndices(state, attestationData, []byte{}); err == nil {
		t.Error("attestation participants should have failed with incorrect bitfield")
	}
}

func TestVerifyBitfield_OK(t *testing.T) {
	bitfield := []byte{0xFF}
	committeeSize := 8

	isValidated, err := VerifyBitfield(bitfield, committeeSize)
	if err != nil {
		t.Fatal(err)
	}

	if !isValidated {
		t.Error("bitfield is not validated when it was supposed to be")
	}

	bitfield = []byte{0xff, 0x80}
	committeeSize = 9

	isValidated, err = VerifyBitfield(bitfield, committeeSize)
	if err != nil {
		t.Fatal(err)
	}

	if isValidated {
		t.Error("bitfield is validated when it was supposed to be")
	}

	bitfield = []byte{0xff, 0x03}
	committeeSize = 10
	isValidated, err = VerifyBitfield(bitfield, committeeSize)
	if err != nil {
		t.Fatal(err)
	}

	if !isValidated {
		t.Error("bitfield is not validated when it was supposed to be")
	}
}

func TestCommitteeAssignment_CanRetrieve(t *testing.T) {
	// Initialize test with 128 validators, each slot and each shard gets 2 validators.
	validators := make([]*pb.Validator, 2*params.BeaconConfig().SlotsPerEpoch)
	for i := 0; i < len(validators); i++ {
		validators[i] = &pb.Validator{
			ExitEpoch: params.BeaconConfig().FarFutureEpoch,
		}
	}
	state := &pb.BeaconState{
		ValidatorRegistry:      validators,
		Slot:                   params.BeaconConfig().SlotsPerEpoch + params.BeaconConfig().GenesisSlot,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}

	tests := []struct {
		index      uint64
		slot       uint64
		committee  []uint64
		shard      uint64
		isProposer bool
	}{
		{
			index:      0,
			slot:       161,
			committee:  []uint64{0, 107},
			shard:      97,
			isProposer: false,
		},
		{
			index:      105,
			slot:       156,
			committee:  []uint64{88, 105},
			shard:      92,
			isProposer: false,
		},
		{
			index:      64,
			slot:       172,
			committee:  []uint64{64, 31},
			shard:      108,
			isProposer: false,
		},
		{
			index:      11,
			slot:       169,
			committee:  []uint64{13, 11},
			shard:      105,
			isProposer: false,
		},
	}

	for _, tt := range tests {
		committee, shard, slot, isProposer, err := CommitteeAssignment(state, tt.slot/params.BeaconConfig().SlotsPerEpoch, tt.index)
		if err != nil {
			t.Fatalf("failed to execute NextEpochCommitteeAssignment: %v", err)
		}
		if shard != tt.shard {
			t.Errorf("wanted shard %d, got shard %d for validator index %d",
				tt.shard, shard, tt.index)
		}
		if slot != tt.slot {
			t.Errorf("wanted slot %d, got slot %d for validator index %d",
				tt.slot, slot, tt.index)
		}
		if isProposer != tt.isProposer {
			t.Errorf("wanted isProposer %v, got isProposer %v for validator index %d",
				tt.isProposer, isProposer, tt.index)
		}
		if !reflect.DeepEqual(committee, tt.committee) {
			t.Errorf("wanted committee %v, got committee %v for validator index %d",
				tt.committee, committee, tt.index)
		}
	}
}

func TestCommitteeAssignment_CantFindValidator(t *testing.T) {
	state := &pb.BeaconState{
		Slot:                   params.BeaconConfig().GenesisSlot + params.BeaconConfig().SlotsPerEpoch,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}
	index := uint64(10000)
	_, _, _, _, err := CommitteeAssignment(state, 1, index)
	statusErr, ok := status.FromError(err)
	if !ok {
		t.Fatal(err)
	}
	if statusErr.Code() != codes.NotFound {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAttestationParticipants_CommitteeCacheHit(t *testing.T) {
	// TODO(2682): Don't fix this test, this will be removed after merging #2682
	t.Skip()

	slotOffset := uint64(1111)
	csInSlot := &cache.CommitteesInSlot{
		Slot: params.BeaconConfig().GenesisSlot + slotOffset,
		Committees: []*cache.CommitteeInfo{
			{Shard: 123, Committee: []uint64{55, 105}},
			{Shard: 234, Committee: []uint64{11, 14}},
		}}

	if err := committeeCache.AddCommittees(csInSlot); err != nil {
		t.Fatal(err)
	}

	attestationData := &pb.AttestationData{
		Shard: 234,
		Slot:  params.BeaconConfig().GenesisSlot + uint64(slotOffset),
	}
	result, err := AttestingIndices(&pb.BeaconState{}, attestationData, []byte{0x03})
	if err != nil {
		t.Fatal(err)
	}

	wanted := []uint64{11, 14}
	if !reflect.DeepEqual(wanted, result) {
		t.Errorf(
			"Result indices was an unexpected value. Wanted %d, got %d",
			wanted,
			result,
		)
	}
}

func TestAttestationParticipants_CommitteeCacheMissSaved(t *testing.T) {
	// TODO(2682): Don't fix this test, this will be removed after merging #2682
	t.Skip()

	validators := make([]*pb.Validator, 2*params.BeaconConfig().SlotsPerEpoch)
	for i := 0; i < len(validators); i++ {
		validators[i] = &pb.Validator{
			ExitEpoch: params.BeaconConfig().FarFutureEpoch,
		}
	}

	slotOffset := uint64(10)
	state := &pb.BeaconState{
		Slot:                   params.BeaconConfig().GenesisSlot + slotOffset,
		ValidatorRegistry:      validators,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}

	attestationData := &pb.AttestationData{
		Shard: 11,
		Slot:  params.BeaconConfig().GenesisSlot + slotOffset,
	}
	result, err := AttestingIndices(state, attestationData, []byte{0x03})
	if err != nil {
		t.Fatal(err)
	}

	wanted := []uint64{49, 92}
	if !reflect.DeepEqual(wanted, result) {
		t.Errorf(
			"Result indices was an unexpected value. Wanted %d, got %d",
			wanted,
			result,
		)
	}

	// Verify the committee for offset slot was cached.
	fetchedCommittees, err := committeeCache.CommitteesInfoBySlot(params.BeaconConfig().GenesisSlot + slotOffset)
	if err != nil {
		t.Fatal(err)
	}
	wanted = []uint64{92, 49}
	if !reflect.DeepEqual(wanted, fetchedCommittees.Committees[0].Committee) {
		t.Errorf(
			"Result indices was an unexpected value. Wanted %d, got %d",
			wanted,
			fetchedCommittees.Committees[0].Committee,
		)
	}
}

func TestCommitteeAssignment_CommitteeCacheHit(t *testing.T) {
	// TODO(2682): Don't fix this test, this will be removed after merging #2682
	t.Skip()

	slotOffset := uint64(1111)
	csInSlot := &cache.CommitteesInSlot{
		Slot: params.BeaconConfig().GenesisSlot + slotOffset,
		Committees: []*cache.CommitteeInfo{
			{Shard: 123, Committee: []uint64{55, 105}},
			{Shard: 234, Committee: []uint64{11, 14}},
		}}

	if err := committeeCache.AddCommittees(csInSlot); err != nil {
		t.Fatal(err)
	}

	beaconState := &pb.BeaconState{
		Slot:                   csInSlot.Slot,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}

	committee, shard, _, isProposer, err :=
		CommitteeAssignment(beaconState, csInSlot.Slot, 105)
	if err != nil {
		t.Fatal(err)
	}

	wanted := []uint64{55, 105}
	if !reflect.DeepEqual(wanted, committee) {
		t.Errorf(
			"Result indices was an unexpected value. Wanted %d, got %d",
			wanted,
			committee,
		)
	}
	if shard != csInSlot.Committees[0].Shard {
		t.Errorf(
			"Result shard was an expected value. Wanted %d, got %d",
			csInSlot.Committees[0].Shard,
			shard,
		)
	}
	if !isProposer {
		t.Error("Wanted proposer true")
	}
}

func TestCommitteeAssignment_CommitteeCacheMissSaved(t *testing.T) {
	// TODO(2682): Don't fix this test, this will be removed after merging #2682
	t.Skip()

	validators := make([]*pb.Validator, 2*params.BeaconConfig().SlotsPerEpoch)
	for i := 0; i < len(validators); i++ {
		validators[i] = &pb.Validator{
			ExitEpoch: params.BeaconConfig().FarFutureEpoch,
		}
	}

	slotOffset := uint64(10)
	state := &pb.BeaconState{
		Slot:                   params.BeaconConfig().GenesisSlot + slotOffset,
		ValidatorRegistry:      validators,
		LatestRandaoMixes:      make([][]byte, params.BeaconConfig().LatestRandaoMixesLength),
		LatestActiveIndexRoots: make([][]byte, params.BeaconConfig().LatestActiveIndexRootsLength),
	}

	committee, shard, _, isProposer, err :=
		CommitteeAssignment(state, slotOffset, 105)
	if err != nil {
		t.Fatal(err)
	}

	wantedCommittee := []uint64{44, 105}
	if !reflect.DeepEqual(wantedCommittee, committee) {
		t.Errorf(
			"Result indices was an unexpected value. Wanted %d, got %d",
			wantedCommittee,
			committee,
		)
	}

	wantedShard := uint64(43)
	if shard != wantedShard {
		t.Errorf(
			"Result shard was an expected value. Wanted %d, got %d",
			wantedShard,
			shard,
		)
	}
	if isProposer {
		t.Error("Wanted proposer false")
	}

	// Verify the committee for offset slot was cached.
	fetchedCommittees, err := committeeCache.CommitteesInfoBySlot(params.BeaconConfig().GenesisSlot + slotOffset)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(wantedCommittee, fetchedCommittees.Committees[0].Committee) {
		t.Errorf(
			"Result indices was an unexpected value. Wanted %d, got %d",
			wantedCommittee,
			fetchedCommittees.Committees[0].Committee,
		)
	}
}
