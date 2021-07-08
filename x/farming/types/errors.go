package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// farming module sentinel errors
var (
	ErrInvalidPlanType      = sdkerrors.Register(ModuleName, 2, "invalid plan type")
	ErrInvalidPlanEndTime   = sdkerrors.Register(ModuleName, 3, "invalid plan end time")
	ErrInvalidPlanEpochDays = sdkerrors.Register(ModuleName, 4, "invalid plan epoch days")

	// TODO: TBD more err types for farming module
)
