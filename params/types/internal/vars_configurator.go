package internal

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

type GlobalVarsConfigurator struct {
}

func One() *GlobalVarsConfigurator {
	return &GlobalVarsConfigurator{}
}

func newU64(u uint64) *uint64 {
	return &u
}

func (g GlobalVarsConfigurator) GetAccountStartNonce() *uint64 {
	return newU64(0)
}

func (g GlobalVarsConfigurator) SetAccountStartNonce(n *uint64) error {
	if n == nil {
		return nil
	}
	if *n != 0 {
		return ctypes.ErrUnsupportedConfigFatal
	}
	return nil
}

func (g GlobalVarsConfigurator) GetMaximumExtraDataSize() *uint64 {
	return newU64(vars.MaximumExtraDataSize)
}

func (g GlobalVarsConfigurator) SetMaximumExtraDataSize(n *uint64) error {
	vars.MaximumExtraDataSize = *n
	return nil
}

func (g GlobalVarsConfigurator) GetMinGasLimit() *uint64 {
	return newU64(vars.MinGasLimit)
}

func (g GlobalVarsConfigurator) SetMinGasLimit(n *uint64) error {
	vars.MinGasLimit = *n
	return nil
}

func (g GlobalVarsConfigurator) GetGasLimitBoundDivisor() *uint64 {
	return newU64(vars.GasLimitBoundDivisor)
}

func (g GlobalVarsConfigurator) SetGasLimitBoundDivisor(n *uint64) error {
	vars.GasLimitBoundDivisor = *n
	return nil
}

func (g GlobalVarsConfigurator) GetMaxCodeSize() *uint64 {
	return newU64(vars.MaxCodeSize)
}

func (g GlobalVarsConfigurator) SetMaxCodeSize(n *uint64) error {
	if n == nil {
		return nil
	}
	vars.MaxCodeSize = *n
	return nil
}

func (c GlobalVarsConfigurator) GetEthashMinimumDifficulty() *big.Int {
	return vars.MinimumDifficulty
}
func (c GlobalVarsConfigurator) SetEthashMinimumDifficulty(i *big.Int) error {
	if i == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	vars.MinimumDifficulty = i
	return nil
}

func (c GlobalVarsConfigurator) GetEthashDifficultyBoundDivisor() *big.Int {
	return vars.DifficultyBoundDivisor
}

func (c GlobalVarsConfigurator) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	if i == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	vars.DifficultyBoundDivisor = i
	return nil
}

func (c GlobalVarsConfigurator) GetEthashDurationLimit() *big.Int {
	return vars.DurationLimit
}

func (c GlobalVarsConfigurator) SetEthashDurationLimit(i *big.Int) error {
	if i == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	vars.DurationLimit = i
	return nil
}
