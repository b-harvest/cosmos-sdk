package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// tax module sentinel errors
var (
	ErrTaxNotExists             = sdkerrors.Register(ModuleName, 2, "tax not exists")
	ErrTaxTypeNotExists         = sdkerrors.Register(ModuleName, 3, "tax type not exists")
	ErrInvalidTaxType           = sdkerrors.Register(ModuleName, 4, "invalid tax type")
	ErrInvalidTaxEndTime        = sdkerrors.Register(ModuleName, 5, "invalid tax end time")
	ErrInvalidTaxEpochDays      = sdkerrors.Register(ModuleName, 6, "invalid tax epoch days")
	ErrInvalidTaxEpochRatio     = sdkerrors.Register(ModuleName, 7, "invalid tax epoch ratio")
	ErrEmptyEpochAmount          = sdkerrors.Register(ModuleName, 8, "epoch amount must not be empty")
	ErrEmptyStakingCoinWeights   = sdkerrors.Register(ModuleName, 9, "staking coin weights must not be empty")
	ErrStakingNotExists          = sdkerrors.Register(ModuleName, 10, "staking not exists")
	ErrRewardNotExists           = sdkerrors.Register(ModuleName, 10, "reward not exists")
	ErrInsufficientStakingAmount = sdkerrors.Register(ModuleName, 11, "insufficient staking amount")
)
