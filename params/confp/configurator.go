// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package confp

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"sort"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

var (
	// compatibleProtocolNameSchemes define matchable naming schemes used by configuration methods
	// that are not incompatible with configuration either having or lacking them.
	compatibleProtocolNameSchemes = []string{
		"ECBP", // "Ethereum Classic Best Practice"
		"EBP",  // "Ethereum Best Practice"
	}
)

func nameSignalsCompatibility(name string) bool {
	for _, s := range compatibleProtocolNameSchemes {
		if regexp.MustCompile(s).MatchString(name) {
			return true
		}
	}
	return false
}

var (
	blockTransitionNamePattern = regexp.MustCompile(`(?m)^Get.+Transition$`)
	timeTransitionNamePattern  = regexp.MustCompile(`(?m)^Get.+TransitionTime$`)
)

func nameSignalsBlockBasedFork(name string) bool {
	return blockTransitionNamePattern.MatchString(name)
}

func nameSignalsTimeBasedFork(name string) bool {
	return timeTransitionNamePattern.MatchString(name)
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string

	// block numbers of the stored and new configurations if block based forking
	StoredBlock, NewBlock *big.Int

	// timestamps of the stored and new configurations if time based forking
	StoredTime, NewTime *uint64

	// the block number to which the local chain must be rewound to correct the error
	RewindToBlock uint64

	// the timestamp to which the local chain must be rewound to correct the error
	RewindToTime uint64
}

// func NewCompatError(what string, storedblock, newblock *uint64) *ConfigCompatError {
// 	var rew *uint64
// 	switch {
// 	case storedblock == nil:
// 		rew = newblock
// 	case newblock == nil || *storedblock < *newblock:
// 		rew = storedblock
// 	default:
// 		rew = newblock
// 	}
// 	err := &ConfigCompatError{what, storedblock, newblock, 0}
// 	if rew != nil && *rew > 0 {
// 		err.RewindTo = *rew - 1
// 	}
// 	return err
// }

func newBlockCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{
		What:          what,
		StoredBlock:   storedblock,
		NewBlock:      newblock,
		RewindToBlock: 0,
	}
	if rew != nil && rew.Sign() > 0 {
		err.RewindToBlock = rew.Uint64() - 1
	}
	return err
}

func newTimestampCompatError(what string, storedtime, newtime *uint64) *ConfigCompatError {
	var rew *uint64
	switch {
	case storedtime == nil:
		rew = newtime
	case newtime == nil || *storedtime < *newtime:
		rew = storedtime
	default:
		rew = newtime
	}
	err := &ConfigCompatError{
		What:         what,
		StoredTime:   storedtime,
		NewTime:      newtime,
		RewindToTime: 0,
	}
	if rew != nil && *rew > 0 {
		err.RewindToTime = *rew - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	if err.StoredBlock != nil {
		return fmt.Sprintf("mismatching %s in database (have block %d, want block %d, rewindto block %d)", err.What, err.StoredBlock, err.NewBlock, err.RewindToBlock)
	}
	return fmt.Sprintf("mismatching %s in database (have timestamp %d, want timestamp %d, rewindto timestamp %d)", err.What, err.StoredTime, err.NewTime, err.RewindToTime)
}

type ConfigValidError struct {
	What string
	A, B interface{}
}

func NewValidErr(what string, a, b interface{}) *ConfigValidError {
	return &ConfigValidError{
		What: what,
		A:    a,
		B:    b,
	}
}

func (err *ConfigValidError) Error() string {
	return fmt.Sprintf("%s, %v/%v", err.What, err.A, err.B)
}

func IsEmpty(anything interface{}) bool {
	if anything == nil {
		return true
	}
	return reflect.DeepEqual(anything, reflect.Zero(reflect.TypeOf(anything)).Interface())
}

func IsValid(conf ctypes.ChainConfigurator, head *uint64) *ConfigValidError {
	// head-agnostic logic
	if conf.GetNetworkID() == nil {
		return NewValidErr("NetworkID cannot be nil", "!=nil", conf.GetNetworkID())
	}
	if head == nil {
		return nil
	}

	// head-full logic
	var bhead = new(big.Int).SetUint64(*head)

	if conf.IsEnabled(conf.GetEIP155Transition, bhead) && conf.GetChainID() == nil {
		return NewValidErr("EIP155 requires ChainID. A:EIP155/B:ChainID", conf.GetEIP155Transition(), conf.GetChainID())
	}

	return nil
}

// Compatible checks whether the two configurations are compatible with each other.
// It returns an error if the configurations are incompatible, or nil if they are.
// If headBlock is nil, it will only check the time-based fork configurations.
// If headBlock is not nil, it will check the block-based fork configurations.
func Compatible(headBlock *big.Int, headTime *uint64, a, b ctypes.ChainConfigurator) *ConfigCompatError {
	// Iterate checkCompatible to find the lowest conflict.
	var lastErr *ConfigCompatError

	for {
		err := compatible(headBlock, headTime, a, b)
		if err == nil || (lastErr != nil && err.RewindToBlock == lastErr.RewindToBlock && err.RewindToTime == lastErr.RewindToTime) {
			break
		}
		lastErr = err

		if err.RewindToTime > 0 {
			headTime = new(uint64)
			*headTime = err.RewindToTime
		} else {
			headBlock = new(big.Int).SetUint64(err.RewindToBlock)
		}
	}

	return lastErr
}

func compatible(headBlock *big.Int, headTime *uint64, a, b ctypes.ChainConfigurator) *ConfigCompatError {
	aFns, aNames := Transitions(a)
	bFns, _ := Transitions(b)
	// Handle forks by block.
	if headBlock != nil {
		for i := range aFns {
			// Skip cross-compatible namespaced transition names, assuming
			// these will not be enforced as hardforks.
			if nameSignalsCompatibility(aNames[i]) {
				continue
			}
			// Skip time-based forks. These are checked separately.
			if !nameSignalsBlockBasedFork(aNames[i]) {
				continue
			}

			// Tolerates nil values.
			check := func(c1, c2, head *big.Int) *ConfigCompatError {
				if isBlockForkIncompatible(c1, c2, head) {
					return newBlockCompatError("incompatible fork value: "+aNames[i], c1, c2)
				}
				return nil
			}

			// Call the respective functions; these return pointers.
			// We need to dereference them to get the actual values for the comparison
			// before converting to big.Int pointers for the actual check.
			av := aFns[i]()
			bv := bFns[i]()
			var aBig, bBig *big.Int
			if av != nil {
				aBig = new(big.Int).SetUint64(*av)
			}
			if bv != nil {
				bBig = new(big.Int).SetUint64(*bv)
			}
			if err := check(aBig, bBig, headBlock); err != nil {
				return err
			}
		}
		if a.IsEnabled(a.GetEIP155Transition, headBlock) {
			if a.GetChainID().Cmp(b.GetChainID()) != 0 {
				ta := a.GetEIP155Transition()
				tb := b.GetEIP155Transition()
				tai := new(big.Int).SetUint64(*ta)
				tbi := new(big.Int).SetUint64(*tb)
				return newBlockCompatError("mismatching chain ids after EIP155 transition", tai, tbi)
			}
		}
	}

	// Handle forks by time.
	if headTime != nil {
		for i, afn := range aFns {
			// Skip cross-compatible namespaced transition names, assuming
			// these will not be enforced as hardforks.
			if nameSignalsCompatibility(aNames[i]) {
				continue
			}
			// Skip time-based forks. These are checked separately.
			if !nameSignalsTimeBasedFork(aNames[i]) {
				continue
			}

			// Check blocks.
			check := func(c1, c2, head *uint64) *ConfigCompatError {
				if isTimeForkIncompatible(c1, c2, head) {
					return newTimestampCompatError("incompatible fork value: "+aNames[i], c1, c2)
				}
				return nil
			}
			if err := check(afn(), bFns[i](), headTime); err != nil {
				return err
			}
		}
	}

	return nil
}

// isBigNilOrMaxed returns true if the given big.Int is nil or has a value of
// any math max value (uint64, int64, int, int32, int16, int8).
func isBigNilOrMaxed(b *big.Int) bool {
	if b == nil {
		return true
	}
	x := b.Uint64()
	return x == math.MaxUint64 ||
		x == 0x7FFFFFFFFFFFFFFF ||
		x == 0x7FFFFFFFFFFFFFF ||
		x == 0x7FFFFFFFFFFFFF
}

func isUint64PNilOrMaxed(p *uint64) bool {
	if p == nil {
		return true
	}
	return *p == math.MaxUint64 ||
		*p == 0x7FFFFFFFFFFFFFFF ||
		*p == 0x7FFFFFFFFFFFFFF ||
		*p == 0x7FFFFFFFFFFFFF
}

func Equivalent(a, b ctypes.ChainConfigurator) error {
	if a.GetConsensusEngineType() != b.GetConsensusEngineType() {
		return fmt.Errorf("mismatch consensus engine types, A: %s, B: %s", a.GetConsensusEngineType(), b.GetConsensusEngineType())
	}

	// Check forks sameness.
	fa, fb := BlockForks(a), BlockForks(b)
	if len(fa) != len(fb) {
		return fmt.Errorf("different block-fork count: %d / %d (%v / %v)", len(fa), len(fb), fa, fb)
	}
	for i := range fa {
		if fa[i] != fb[i] {
			if fa[i] == math.MaxUint64 {
				return fmt.Errorf("fa bigmax: %d", fa[i])
			}
			if fb[i] == math.MaxUint64 {
				return fmt.Errorf("fb bigmax: %d", fb[i])
			}
			return fmt.Errorf("block-fork index %d not same: %d / %d", i, fa[i], fb[i])
		}
	}

	fat, fbt := TimeForks(a, 0), TimeForks(b, 0)
	if len(fat) != len(fbt) {
		return fmt.Errorf("different time-fork count: %d / %d (%v / %v)", len(fat), len(fbt), fat, fbt)
	}
	for i := range fat {
		if fat[i] != fbt[i] {
			if fat[i] == math.MaxUint64 {
				return fmt.Errorf("fat bigmax: %d", fat[i])
			}
			if fbt[i] == math.MaxUint64 {
				return fmt.Errorf("fbt bigmax: %d", fbt[i])
			}
			return fmt.Errorf("time-fork index %d not same: %d / %d", i, fat[i], fbt[i])
		}
	}

	// Check initial, at- and around-fork, and eventual compatibility.
	var blockForks = []uint64{}
	copy(blockForks, fa)
	// Don't care about dupes.
	for _, f := range fa {
		blockForks = append(blockForks, f-1)
	}
	blockForks = append(blockForks, 0, math.MaxUint64)

	var timeForks = []uint64{}
	copy(timeForks, fat)
	for _, f := range fat {
		timeForks = append(timeForks, f-1)
	}
	timeForks = append(timeForks, 0, math.MaxUint64)

	// blockHeightsEssentiallyEquivalent treats nil and bitsize-max numbers as essentially equivalent.
	blockHeightsEssentiallyEquivalent := func(bigX, bigY *big.Int) bool {
		if isBigNilOrMaxed(bigX) && isBigNilOrMaxed(bigY) {
			return true
		}
		if isBigNilOrMaxed(bigX) && !isBigNilOrMaxed(bigY) {
			return false
		}
		if !isBigNilOrMaxed(bigX) && isBigNilOrMaxed(bigY) {
			return false
		}
		return bigX.Cmp(bigY) == 0
	}

	for _, h := range blockForks {
		if err := Compatible(new(big.Int).SetUint64(h), nil, a, b); err != nil { // nolint:gosec
			if !blockHeightsEssentiallyEquivalent(err.StoredBlock, err.NewBlock) {
				return err
			}
		}
	}

	blockTimesEssentiallyEquivalent := func(x, y *uint64) bool {
		if isUint64PNilOrMaxed(x) && isUint64PNilOrMaxed(y) {
			return true
		}
		if isUint64PNilOrMaxed(x) && !isUint64PNilOrMaxed(y) {
			return false
		}
		if !isUint64PNilOrMaxed(x) && isUint64PNilOrMaxed(y) {
			return false
		}
		return *x == *y
	}

	for _, t := range timeForks {
		if err := Compatible(nil, &t, a, b); err != nil {
			if !blockTimesEssentiallyEquivalent(err.StoredTime, err.NewTime) {
				return err
			}
		}
	}

	if a.GetConsensusEngineType() == ctypes.ConsensusEngineT_Ethash {
		for _, f := range fa { // fa and fb are fork-equivalent
			ar := ctypes.EthashBlockReward(a, new(big.Int).SetUint64(f))
			br := ctypes.EthashBlockReward(b, new(big.Int).SetUint64(f))
			if ar.Cmp(br) != 0 {
				return fmt.Errorf("mismatch block reward, fork block: %v, A: %v, B: %v", f, ar, br)
			}
			// TODO: add difficulty comparison
			// Currently tough/complex to do because of necessary overhead (ie build a parent block).
		}
	} else if a.GetConsensusEngineType() == ctypes.ConsensusEngineT_Clique {
		if a.GetCliqueEpoch() != b.GetCliqueEpoch() {
			return fmt.Errorf("mismatch clique epochs: A: %v, B: %v", a.GetCliqueEpoch(), b.GetCliqueEpoch())
		}
		if a.GetCliquePeriod() != b.GetCliquePeriod() {
			return fmt.Errorf("mismatch clique periods: A: %v, B: %v", a.GetCliquePeriod(), b.GetCliquePeriod())
		}
	}
	return nil
}

// Transitions gets all available transition (fork) functions and their names for a ChainConfigurator.
func Transitions(conf ctypes.ChainConfigurator) (fns []func() *uint64, names []string) {
	names = []string{}
	fns = []func() *uint64{}
	k := reflect.TypeOf(conf)
	for i := 0; i < k.NumMethod(); i++ {
		method := k.Method(i)
		if !nameSignalsBlockBasedFork(method.Name) && !nameSignalsTimeBasedFork(method.Name) {
			continue
		}
		m := reflect.ValueOf(conf).MethodByName(method.Name).Interface()
		fns = append(fns, m.(func() *uint64))
		names = append(names, method.Name)
	}
	return fns, names
}

// BlockForks returns non-nil, non <maxUin64>, unique sorted forks defined by block number for a ChainConfigurator.
func BlockForks(conf ctypes.ChainConfigurator) []uint64 {
	var forks []uint64
	var forksM = make(map[uint64]struct{}) // Will key for uniqueness as fork numbers are appended to slice.

	transitions, names := Transitions(conf)
	for i, tr := range transitions {
		// Skip time-based transition names, assuming these are not block-based hardforks.
		if nameSignalsTimeBasedFork(names[i]) {
			continue
		}
		// Skip cross-compatible namespaced transition names, assuming
		// these will not be enforced as hardforks.
		if nameSignalsCompatibility(names[i]) {
			continue
		}
		// Extract the fork rule block number and aggregate it
		response := tr()
		if response == nil ||
			*response == math.MaxUint64 ||
			*response == 0x7fffffffffffff ||
			*response == 0x7FFFFFFFFFFFFFFF {
			continue
		}

		// Only append unique fork numbers, excluding 0 (genesis config is not considered a fork)
		if _, ok := forksM[*response]; !ok && *response != 0 {
			forks = append(forks, *response)
			forksM[*response] = struct{}{}
		}
	}
	sort.Slice(forks, func(i, j int) bool {
		return forks[i] < forks[j]
	})

	// Skip any forks in block 0, that's the genesis ruleset
	if len(forks) > 0 && forks[0] == 0 {
		forks = forks[1:]
	}

	return forks
}

// TimeForks returns non-nil, non <maxUin64>, unique sorted forks defined by block time for a ChainConfigurator.
func TimeForks(conf ctypes.ChainConfigurator, genesis uint64) []uint64 {
	var forks []uint64
	var forksM = make(map[uint64]struct{}) // Will key for uniqueness as fork numbers are appended to slice.

	transitions, names := Transitions(conf)
	for i, tr := range transitions {
		// Skip non-time-based transition names, assuming these are not time-based hardforks.
		if !nameSignalsTimeBasedFork(names[i]) {
			continue
		}
		// Skip cross-compatible namespaced transition names, assuming
		// these will not be enforced as hardforks.
		if nameSignalsCompatibility(names[i]) {
			continue
		}
		// Extract the fork rule block number and aggregate it
		response := tr()
		if response == nil ||
			*response == math.MaxUint64 ||
			*response == 0x7fffffffffffff ||
			*response == 0x7FFFFFFFFFFFFFFF {
			continue
		}

		// Only append unique fork numbers, excluding 0 (genesis config is not considered a fork)
		if _, ok := forksM[*response]; !ok && *response != 0 {
			forks = append(forks, *response)
			forksM[*response] = struct{}{}
		}
	}
	sort.Slice(forks, func(i, j int) bool {
		return forks[i] < forks[j]
	})

	// Skip any forks before genesis.
	for len(forks) > 0 && forks[0] <= genesis {
		forks = forks[1:]
	}

	return forks
}

func isBlockForkIncompatible(a, b, head *big.Int) bool {
	// If the head is nil, then either fork config is ok. Return incompatible = false.
	if head == nil {
		return false
	}
	// If one config is forked and the other not, return incompatible = true.
	if isBlockForked(a, head) != isBlockForked(b, head) {
		return true
	}
	// If they are both forked, then we need to check if they are/were forked at the same block.
	// Note that we could technically only check one (since we know they are both forked or both not forked),
	// but checking them both because it is more explicit.
	if isBlockForked(a, head) && isBlockForked(b, head) {
		return a.Cmp(b) != 0
	}
	return false
}

func isBlockForked(x, head *big.Int) bool {
	if x == nil || head == nil {
		return false
	}
	return x.Cmp(head) <= 0
}

func isTimeForkIncompatible(a, b, head *uint64) bool {
	return (isTimeForked(a, head) || isTimeForked(b, head)) && a != b
}

func isTimeForked(x, head *uint64) bool {
	if x == nil || head == nil {
		return false
	}
	return *x <= *head
}

// Uint64Ptr2Big converts a *uint64 to a *big.Int.
// It returns nil if the input is nil.
func Uint64Ptr2Big(x *uint64) *big.Int {
	if x == nil {
		return nil
	}
	return new(big.Int).SetUint64(*x)
}
