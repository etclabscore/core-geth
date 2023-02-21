package internal

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

type GlobalVarsConfigurator struct {
}

var gc = &GlobalVarsConfigurator{}

func GlobalConfigurator() *GlobalVarsConfigurator {
	return gc
}

func newU64(u uint64) *uint64 {
	return &u
}

func (_ GlobalVarsConfigurator) GetAccountStartNonce() *uint64 {
	return newU64(0)
}

func (_ GlobalVarsConfigurator) SetAccountStartNonce(n *uint64) error {
	if n == nil {
		return nil
	}
	if *n != 0 {
		return ctypes.ErrUnsupportedConfigFatal
	}
	return nil
}

func (_ GlobalVarsConfigurator) GetMaximumExtraDataSize() *uint64 {
	return newU64(vars.MaximumExtraDataSize)
}

func (_ GlobalVarsConfigurator) SetMaximumExtraDataSize(n *uint64) error {
	vars.MaximumExtraDataSize = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetMinGasLimit() *uint64 {
	return newU64(vars.MinGasLimit)
}

func (_ GlobalVarsConfigurator) SetMinGasLimit(n *uint64) error {
	vars.MinGasLimit = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetGasLimitBoundDivisor() *uint64 {
	return newU64(vars.GasLimitBoundDivisor)
}

func (_ GlobalVarsConfigurator) SetGasLimitBoundDivisor(n *uint64) error {
	vars.GasLimitBoundDivisor = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetMaxCodeSize() *uint64 {
	return newU64(vars.MaxCodeSize)
}

func (_ GlobalVarsConfigurator) SetMaxCodeSize(n *uint64) error {
	if n == nil {
		return nil
	}
	vars.MaxCodeSize = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetElasticityMultiplier() uint64 {
	return vars.DefaultElasticityMultiplier
}

func (_ GlobalVarsConfigurator) SetElasticityMultiplier(n uint64) error {
	// Noop.
	return nil
}

func (_ GlobalVarsConfigurator) GetBaseFeeChangeDenominator() uint64 {
	return vars.DefaultBaseFeeChangeDenominator
}

func (_ GlobalVarsConfigurator) SetBaseFeeChangeDenominator(n uint64) error {
	// Noop.
	return nil
}

func (_ GlobalVarsConfigurator) GetEthashMinimumDifficulty() *big.Int {
	return vars.MinimumDifficulty
}
func (_ GlobalVarsConfigurator) SetEthashMinimumDifficulty(i *big.Int) error {
	if i == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	vars.MinimumDifficulty = i
	return nil
}

func (_ GlobalVarsConfigurator) GetEthashDifficultyBoundDivisor() *big.Int {
	return vars.DifficultyBoundDivisor
}

func (_ GlobalVarsConfigurator) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	if i == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	vars.DifficultyBoundDivisor = i
	return nil
}

func (_ GlobalVarsConfigurator) GetEthashDurationLimit() *big.Int {
	return vars.DurationLimit
}

func (_ GlobalVarsConfigurator) SetEthashDurationLimit(i *big.Int) error {
	if i == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	vars.DurationLimit = i
	return nil
}
