package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// HandlePublicTaxProposal is a handler for executing a fixed amount tax creation proposal.
func HandleSetTaxesProposal(ctx sdk.Context, k Keeper, taxes []types.Tax) error {
	//taxes, err := types.UnpackTaxes(taxes)
	//if err != nil {
	//	return err
	//}

	for _, tax := range taxes {
		fmt.Println(tax)
		//switch p := tax.(type) {
		//case *types.FixedAmountTax:
		//	msg := types.NewMsgCreateFixedAmountTax(
		//		p.GetTaxPoolAddress(),
		//		p.GetStakingCoinWeights(),
		//		p.GetStartTime(),
		//		p.GetEndTime(),
		//		p.GetEpochDays(),
		//		p.EpochAmount,
		//	)
		//
		//	fixedTax := k.CreateFixedAmountTax(ctx, msg, types.TaxTypePublic)
		//
		//	logger := k.Logger(ctx)
		//	logger.Info("created public fixed amount tax", "fixed_amount_tax", fixedTax)
		//
		//case *types.RatioTax:
		//	msg := types.NewMsgCreateRatioTax(
		//		p.GetTaxPoolAddress(),
		//		p.GetStakingCoinWeights(),
		//		p.GetStartTime(),
		//		p.GetEndTime(),
		//		p.GetEpochDays(),
		//		p.EpochRatio,
		//	)
		//
		//	ratioTax := k.CreateRatioTax(ctx, msg, types.TaxTypePublic)
		//
		//	logger := k.Logger(ctx)
		//	logger.Info("created public fixed amount tax", "ratio_tax", ratioTax)

		//default:
		//	return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tax proposal tax type: %T", p)
		//}
	}

	return nil
}
