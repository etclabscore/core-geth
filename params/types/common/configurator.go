package common

import (
	"fmt"
	"math/big"
	"reflect"
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
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
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

func Valid(conf ChainConfigurator, head *uint64) *ConfigValidError {
	var bhead *big.Int
	if head != nil {
		bhead = new(big.Int).SetUint64(*head)
	}
	// head-full logic
	if conf.IsForked(conf.GetEIP155Transition, bhead) && conf.GetChainID() == nil {
		return NewValidErr("EIP155 requires ChainID. A:EIP155/B:ChainID", conf.GetEIP155Transition(), conf.GetChainID())
	}
	if !u2Equal(conf.GetEthashEIP100BTransition(), conf.GetEthashEIP649TransitionV()) {}
}

func Compatible(head uint64, a, b ChainConfigurator) *ConfigCompatError {
	bhead := &head
	// Iterate checkCompatible to find the lowest conflict.
	var lastErr *ConfigCompatError
	for {
		err := compatible(bhead, a, b)
		if err == nil || (lastErr != nil && err.RewindTo == lastErr.RewindTo) {
			break
		}
		lastErr = err
		*bhead = err.RewindTo
	}
	return lastErr
}

// transitions gets transition functions and their names for a ChainConfigurator.
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


