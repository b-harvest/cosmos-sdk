package farming

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/farming/keeper"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateFixedAmountPlan:
			res, err := msgServer.CreateFixedAmountPlan(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCreateRatioPlan:
			res, err := msgServer.CreateRatioPlan(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgStake:
			res, err := msgServer.Stake(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgUnstake:
			res, err := msgServer.Unstake(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgClaim:
			res, err := msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func NewCreatePublicFarmingPlanProposal(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.FixedAmountPlanProposal:
			return keeper.HandleFixedAmountPlanProposal(ctx, k, c.Plan)

		case *types.RatioPlanProposal:
			return keeper.HandleRatioPlanProposal(ctx, k, c.Plan)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized farming proposal content type: %T", c)
		}
	}
}

func NewModifyPublicFarmingPlanProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.FixedAmountPlanProposal:
			return keeper.ModifyFixedAmountPlanProposal(ctx, k, c.Plan)

		case *types.RatioPlanProposal:
			return keeper.ModifyRatioPlanProposal(ctx, k, c.Plan)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized farming proposal content type: %T", c)
		}
	}
}
