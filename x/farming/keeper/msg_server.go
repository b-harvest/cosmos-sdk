package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the farming MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateFixedAmountPlan defines a method for creating fixed amount farming plan.
func (k msgServer) CreateFixedAmountPlan(goCtx context.Context, msg *types.MsgCreateFixedAmountPlan) (*types.MsgCreateFixedAmountPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	nextId := k.GetNextPlanID(ctx)
	farmingPoolAddr := msg.GetFarmingPoolAddress()
	terminationAddr := farmingPoolAddr

	// TODO: consider having CreateRatioPlan keeper
	basePlan := types.NewBasePlan(
		nextId,
		types.PlanTypePrivate,
		farmingPoolAddr,
		terminationAddr,
		msg.GetStakingCoinWeights(),
		msg.StartTime,
		msg.EndTime,
		msg.GetEpochDays(),
	)

	fixedPlan := types.NewFixedAmountPlan(basePlan, msg.EpochAmount)

	k.SetPlan(ctx, fixedPlan)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedAmountPlan,
			sdk.NewAttribute(types.AttributeKeyFarmingPoolAddress, msg.GetFarmingPoolAddress()),
			sdk.NewAttribute(types.AttributeKeyRewardPoolAddress, fixedPlan.RewardPoolAddress),
			sdk.NewAttribute(types.AttributeKeyStakingReserveAddress, fixedPlan.StakingReserveAddress),
			sdk.NewAttribute(types.AttributeKeyStartTime, msg.StartTime.String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyEpochDays, fmt.Sprint(msg.GetEpochDays())),
			sdk.NewAttribute(types.AttributeKeyEpochAmount, fmt.Sprint(msg.GetEpochAmount())),
		),
	})

	return &types.MsgCreateFixedAmountPlanResponse{}, nil
}

// CreateRatioPlan defines a method for creating ratio farming plan.
func (k msgServer) CreateRatioPlan(goCtx context.Context, msg *types.MsgCreateRatioPlan) (*types.MsgCreateRatioPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	nextId := k.GetNextPlanID(ctx)
	farmingPoolAddr := msg.GetFarmingPoolAddress()
	terminationAddr := farmingPoolAddr

	// TODO: consider having CreateRatioPlan keeper
	basePlan := types.NewBasePlan(
		nextId,
		types.PlanTypePrivate,
		farmingPoolAddr,
		terminationAddr,
		msg.GetStakingCoinWeights(),
		msg.StartTime,
		msg.EndTime,
		msg.GetEpochDays(),
	)

	ratioPlan := types.NewRatioPlan(basePlan, msg.EpochRatio)

	k.SetPlan(ctx, ratioPlan)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateRatioPlan,
			sdk.NewAttribute(types.AttributeKeyFarmingPoolAddress, msg.GetFarmingPoolAddress()),
			sdk.NewAttribute(types.AttributeKeyRewardPoolAddress, ratioPlan.RewardPoolAddress),
			sdk.NewAttribute(types.AttributeKeyStakingReserveAddress, ratioPlan.StakingReserveAddress),
			sdk.NewAttribute(types.AttributeKeyStartTime, msg.StartTime.String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyEpochDays, fmt.Sprint(msg.GetEpochDays())),
			sdk.NewAttribute(types.AttributeKeyEpochRatio, fmt.Sprint(msg.EpochRatio)),
		),
	})

	return &types.MsgCreateRatioPlanResponse{}, nil
}

// Stake defines a method for staking coins to the farming plan.
func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	plan := k.GetPlan(ctx, msg.PlanId)
	if plan == nil {
		return nil, types.ErrPlanNotExists
	}

	// SetStake
	// SetUnstake
	// SetClaim
	// k.Stake(msg.PlanId, msg.Farmer, msg.StakingCoins)

	return &types.MsgStakeResponse{}, nil
}

// Unstake defines a method for unstaking coins from the farming plan.
func (k msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	plan := k.GetPlan(ctx, msg.PlanId)
	if plan == nil {
		return nil, types.ErrPlanNotExists
	}

	return &types.MsgUnstakeResponse{}, nil
}

// Claim defines a method for claiming farming rewards from the farming plan.
func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fmt.Println("ctx: ", ctx)
	return &types.MsgClaimResponse{}, nil
}
