package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// farming module sentinel errors
var (
	ErrPlanNotExists         = sdkerrors.Register(ModuleName, 2, "plan not exists")
	ErrPlanTypeNotExists     = sdkerrors.Register(ModuleName, 3, "plan type not exists")
	ErrInvalidPlanType       = sdkerrors.Register(ModuleName, 4, "invalid plan type")
	ErrInvalidPlanEndTime    = sdkerrors.Register(ModuleName, 5, "invalid plan end time")
	ErrInvalidPlanEpochDays  = sdkerrors.Register(ModuleName, 6, "invalid plan epoch days")
	ErrInvalidPlanEpochRatio = sdkerrors.Register(ModuleName, 7, "invalid plan epoch ratio")
)
