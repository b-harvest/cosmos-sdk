package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

// HandlePublicPlanProposal is a handler for executing a fixed amount plan creation proposal.
func HandlePublicPlanProposal(ctx sdk.Context, k Keeper, plansAny []*codectypes.Any) error {
	plans, err := types.UnpackPlans(plansAny)
	if err != nil {
		return err
	}

	for _, plan := range plans {
		switch p := plan.(type) {
		case *types.FixedAmountPlan:
			// TODO: not implemented
		case *types.RatioPlan:
			// TODO: not implemented
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized farming proposal plan type: %T", p)
		}
	}

	// get plan
	// plan id + 1
	// set plan

	logger := k.Logger(ctx)
	logger.Info("HandleFixedAmountPlanProposal")
	// logger.Info("transferred from the community pool to recipient", "amount", p.Amount.String(), "recipient", p.Recipient)

	return nil
}
