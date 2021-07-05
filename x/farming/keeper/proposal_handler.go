package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

// HandleFixedAmountPlanProposal is a handler for executing a fixed amount plan creation  proposal
func HandleFixedAmountPlanProposal(ctx sdk.Context, k Keeper, p *types.CreateFixedAmountPlanProposal) error {
	logger := k.Logger(ctx)
	logger.Info("HandleFixedAmountPlanProposal")
	// logger.Info("transferred from the community pool to recipient", "amount", p.Amount.String(), "recipient", p.Recipient)

	return nil
}

// HandleRatioPlanProposal is a handler for executing a ratio plan creation proposal
func HandleRatioPlanProposal(ctx sdk.Context, k Keeper, p *types.CreateRatioPlanProposal) error {

	logger := k.Logger(ctx)
	logger.Info("HandleRatioPlanProposal")
	// logger.Info("transferred from the community pool to recipient", "amount", p.Amount.String(), "recipient", p.Recipient)

	return nil
}
