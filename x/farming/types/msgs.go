package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = (*MsgCreateFixedAmountPlan)(nil)
	_ sdk.Msg = (*MsgCreateRatioPlan)(nil)
	_ sdk.Msg = (*MsgStake)(nil)
	_ sdk.Msg = (*MsgUnstake)(nil)
	_ sdk.Msg = (*MsgClaim)(nil)
)

// Message types for the farming module
const (
	TypeMsgCreateFixedAmountPlan = "create_fixed_amount_plan"
	TypeMsgCreateRatioPlan       = "create_ratio_plan"
	TypeMsgStake                 = "stake"
	TypeMsgUnstake               = "unstake"
	TypeMsgClaim                 = "claim"
)

// NewMsgCreateFixedAmountPlan creates a new MsgCreateFixedAmountPlan.
func NewMsgCreateFixedAmountPlan(
	farming_pool_address sdk.AccAddress,
	staking_coins_weight sdk.DecCoins,
	start_time *time.Time,
	end_time *time.Time,
	epoch_days uint32,
	epoch_amount sdk.Coins,
) *MsgCreateFixedAmountPlan {
	return &MsgCreateFixedAmountPlan{
		FarmingPoolAddress: farming_pool_address.String(),
		StakingCoinsWeight: staking_coins_weight,
		StartTime:          start_time,
		EndTime:            end_time,
		EpochDays:          epoch_days,
		EpochAmount:        epoch_amount,
	}
}

func (msg MsgCreateFixedAmountPlan) Route() string { return RouterKey }

func (msg MsgCreateFixedAmountPlan) Type() string { return TypeMsgCreateFixedAmountPlan }

func (msg MsgCreateFixedAmountPlan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.FarmingPoolAddress); err != nil {
		return ErrInvalidFarmingPoolAddress
	}
	if err := msg.EpochAmount.Validate(); err != nil {
		return err
	}
	// TODO: more details for each field
	return nil
}

func (msg MsgCreateFixedAmountPlan) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateFixedAmountPlan) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FarmingPoolAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateFixedAmountPlan) GetPlanCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FarmingPoolAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCreateRatioPlan creates a new MsgCreateRatioPlan.
func NewMsgCreateRatioPlan(
	farming_pool_address sdk.AccAddress,
	staking_coins_weight sdk.DecCoins,
	start_time *time.Time,
	end_time *time.Time,
	epoch_days uint32,
	epoch_ratio sdk.Dec,
) *MsgCreateRatioPlan {
	return &MsgCreateRatioPlan{
		FarmingPoolAddress: farming_pool_address.String(),
		StakingCoinsWeight: staking_coins_weight,
		StartTime:          start_time,
		EndTime:            end_time,
		EpochDays:          epoch_days,
		EpochRatio:         epoch_ratio,
	}
}

func (msg MsgCreateRatioPlan) Route() string { return RouterKey }

func (msg MsgCreateRatioPlan) Type() string { return TypeMsgCreateRatioPlan }

func (msg MsgCreateRatioPlan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.FarmingPoolAddress); err != nil {
		return ErrInvalidFarmingPoolAddress
	}
	// TODO: more details for each field
	return nil
}

func (msg MsgCreateRatioPlan) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateRatioPlan) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FarmingPoolAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateRatioPlan) GetPlanCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FarmingPoolAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgStake creates a new MsgStake.
func NewMsgStake(
	plan_id uint64,
	farmer sdk.AccAddress,
	staking_coins sdk.Coins,
) *MsgStake {
	return &MsgStake{
		PlanId:       plan_id,
		Farmer:       farmer.String(),
		StakingCoins: staking_coins,
	}
}

func (msg MsgStake) Route() string { return RouterKey }

func (msg MsgStake) Type() string { return TypeMsgStake }

func (msg MsgStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return ErrInvalidFarmerAddress
	}
	if err := msg.StakingCoins.Validate(); err != nil {
		return err
	}
	// TODO: more details for each field
	return nil
}

func (msg MsgStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgStake) GetStaker() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgUnstake creates a new MsgUnstake.
func NewMsgUnstake(
	plan_id uint64,
	farmer sdk.AccAddress,
	unstaking_coins sdk.Coins,
) *MsgUnstake {
	return &MsgUnstake{
		PlanId:         plan_id,
		Farmer:         farmer.String(),
		UnstakingCoins: unstaking_coins,
	}
}

func (msg MsgUnstake) Route() string { return RouterKey }

func (msg MsgUnstake) Type() string { return TypeMsgUnstake }

func (msg MsgUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return ErrInvalidFarmerAddress
	}
	if err := msg.UnstakingCoins.Validate(); err != nil {
		return err
	}
	// TODO: more details for each field
	return nil
}

func (msg MsgUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUnstake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgUnstake) GetUnstaker() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgClaim creates a new MsgClaim.
func NewMsgClaim(
	plan_id uint64,
	farmer sdk.AccAddress,
) *MsgClaim {
	return &MsgClaim{
		PlanId: plan_id,
		Farmer: farmer.String(),
	}
}

func (msg MsgClaim) Route() string { return RouterKey }

func (msg MsgClaim) Type() string { return TypeMsgClaim }

func (msg MsgClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return ErrInvalidFarmerAddress
	}
	// TODO: more details for each field
	return nil
}

func (msg MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgClaim) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgClaim) GetClaimer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}
