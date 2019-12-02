package common

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"sort"
	"strings"
)

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *uint64
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func NewCompatError(what string, storedblock, newblock *uint64) *ConfigCompatError {
	var rew *uint64
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || *storedblock < *newblock:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && *rew > 0 {
		err.RewindTo = *rew - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	var have, want interface{}
	if err.StoredConfig != nil {
		have = *err.StoredConfig
	}
	if err.NewConfig != nil {
		want = *err.NewConfig
	}
	if have == nil {
		have = "nil"
	}
	if want == nil {
		want = "nil"
	}
	return fmt.Sprintf("mismatching %s in database (have %v, want %v, rewindto %d)", err.What, have, want, err.RewindTo)
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

func IsValid(conf ChainConfigurator, head *uint64) *ConfigValidError {

	// head-agnostic logic
	if conf.GetNetworkID() == nil || *conf.GetNetworkID() == 0 {
		return NewValidErr("NetworkID cannot be empty nor zero", ">=0", conf.GetNetworkID())
	}
	if head == nil {
		return nil
	}

	// head-full logic
	var bhead = new(big.Int).SetUint64(*head)

	if conf.IsForked(conf.GetEIP155Transition, bhead) && conf.GetChainID() == nil {
		return NewValidErr("EIP155 requires ChainID. A:EIP155/B:ChainID", conf.GetEIP155Transition(), conf.GetChainID())
	}

	return nil
}

func Compatible(head *uint64, a, b ChainConfigurator) *ConfigCompatError {
	// Iterate checkCompatible to find the lowest conflict.
	var lastErr *ConfigCompatError
	for {
		err := compatible(head, a, b)
		if err == nil || (lastErr != nil && err.RewindTo == lastErr.RewindTo) {
			break
		}
		lastErr = err
		head = &err.RewindTo
	}
	return lastErr
}

// FIXME(meows): Incomplete.
func Equivalent(a, b ChainConfigurator) error {
	var m uint64 = math.MaxUint64
	if err := Compatible(&m, a, b); err != nil {
		return err
	}

	fa, fb := Forks(a), Forks(b)
	if len(fa) != len(fb) {
		return fmt.Errorf("different fork count: %d / %d (%v / %v)", len(fa), len(fb), fa, fb)
	}
	for i := range fa {
		if fa[i] != fb[i] {
			if fa[i] == math.MaxUint64 {
				return fmt.Errorf("fa bigmax: %d", fa[i])
			}
			if fb[i] == math.MaxUint64 {
				return fmt.Errorf("fb bigmax: %d", fb[i])
			}
			return fmt.Errorf("fork index %d not same: %d / %d", i, fa[i], fb[i])
		}
	}
	return nil
}

// Transitions gets all available transition (fork) functions and their names for a ChainConfigurator.
func Transitions(conf ChainConfigurator) (fns []func() *uint64, names []string) {
	names = []string{}
	fns = []func() *uint64{}
	k := reflect.TypeOf(conf)
	for i := 0; i < k.NumMethod(); i++ {
		method := k.Method(i)
		if !strings.HasPrefix(method.Name, "Get") || !strings.HasSuffix(method.Name, "Transition") {
			continue
		}
		m := reflect.ValueOf(conf).MethodByName(method.Name).Interface()
		fns = append(fns, m.(func() *uint64))
		names = append(names, method.Name)
	}
	return fns, names
}

// Forks returns non-nil, non <maxUin64>, unique sorted forks for a ChainConfigurator.
func Forks(conf ChainConfigurator) []uint64 {
	var forks []uint64
	var forksM = make(map[uint64]struct{}) // Will key for uniqueness as fork numbers are appended to slice.

	transitions, _ := Transitions(conf)
	for _, tr := range transitions {
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

	return forks
}

func compatible(head *uint64, a, b ChainConfigurator) *ConfigCompatError {
	aFns, aNames := Transitions(a)
	bFns, _ := Transitions(b)
	for i, afn := range aFns {
		if err := func(c1, c2, head *uint64) *ConfigCompatError {
			if isForkIncompatible(c1, c2, head) {
				return NewCompatError("incompatible fork value: "+aNames[i], c1, c2)
			}
			return nil
		}(afn(), bFns[i](), head); err != nil {
			return err
		}
	}
	if head == nil {
		return nil
	}
	if a.IsForked(a.GetEIP155Transition, new(big.Int).SetUint64(*head)) {
		if a.GetChainID().Cmp(b.GetChainID()) != 0 {
			return NewCompatError("mismatching chain ids after EIP155 transition", a.GetEIP155Transition(), b.GetEIP155Transition())
		}
	}

	adao, bdao := a.GetEthashEIP779Transition(), b.GetEthashEIP779Transition()
	if adao == nil && bdao == nil {
		return nil
	}

	headB := new(big.Int).SetUint64(*head)
	if !a.IsForked(a.GetEthashEIP779Transition, headB) && !b.IsForked(b.GetEthashEIP779Transition, headB) {
		return nil
	}

	if *adao != *bdao {
		return NewCompatError("mismatching DAO fork", adao, bdao)

	}

	return nil
}

func isForkIncompatible(a, b, head *uint64) bool {
	return (isForked(a, head) || isForked(b, head)) && !u2Equal(a, b)
}

func isForked(x, head *uint64) bool {
	if x == nil || head == nil {
		return false
	}
	return *x <= *head
}

func u2Equal(x, y *uint64) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return *x == *y
}
