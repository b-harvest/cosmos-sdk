package tax

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/tax/keeper"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	//msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		//case *types.MsgCreateFixedAmountTax:
		//	res, err := msgServer.CreateFixedAmountTax(sdk.WrapSDKContext(ctx), msg)
		//	return sdk.WrapServiceResult(ctx, res, err)
		//
		//case *types.MsgCreateRatioTax:
		//	res, err := msgServer.CreateRatioTax(sdk.WrapSDKContext(ctx), msg)
		//	return sdk.WrapServiceResult(ctx, res, err)
		//
		//case *types.MsgStake:
		//	res, err := msgServer.Stake(sdk.WrapSDKContext(ctx), msg)
		//	return sdk.WrapServiceResult(ctx, res, err)
		//
		//case *types.MsgUnstake:
		//	res, err := msgServer.Unstake(sdk.WrapSDKContext(ctx), msg)
		//	return sdk.WrapServiceResult(ctx, res, err)
		//
		//case *types.MsgClaim:
		//	res, err := msgServer.Claim(sdk.WrapSDKContext(ctx), msg)
		//	return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func NewSetTaxesProposal(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.SetTaxesProposal:
			return keeper.HandleSetTaxesProposal(ctx, k, c.Taxes)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tax proposal content type: %T", c)
		}
	}
}
