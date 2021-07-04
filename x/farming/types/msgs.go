package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = (*MsgCreateFixedAmountPlan)(nil)

	// TODO: implemented interfaces for each msgs
	//_ sdk.Msg = (*MsgCreateRatioPlan)(nil)
	//_ sdk.Msg = (*MsgStake)(nil)
	//_ sdk.Msg = (*MsgUnstake)(nil)
	//_ sdk.Msg = (*MsgClaim)(nil)
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
func NewMsgCreateFixedAmountPlan() *MsgCreateFixedAmountPlan {
	return &MsgCreateFixedAmountPlan{
		// TODO: unimplemented
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
