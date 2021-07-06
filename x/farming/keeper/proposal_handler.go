package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

// HandleFixedAmountPlanProposal is a handler for executing a fixed amount plan creation proposal.
func HandleFixedAmountPlanProposal(ctx sdk.Context, k Keeper, plan types.FixedAmountPlan) error {
	// TODO: not implemented yet

	// checks needed in Validate:
	// 1. type check
	// 2. address check
	// 3. staking coins weight
	// 4. start time, end time
	// 5. epoch days - integer
	if err := plan.Validate(); err != nil {
		return err
	}

	// get plan
	// plan id + 1
	// set plan

	logger := k.Logger(ctx)
	logger.Info("HandleFixedAmountPlanProposal")
	// logger.Info("transferred from the community pool to recipient", "amount", p.Amount.String(), "recipient", p.Recipient)

	return nil
}

// HandleRatioPlanProposal is a handler for executing a ratio plan creation proposal.
func HandleRatioPlanProposal(ctx sdk.Context, k Keeper, plan types.RatioPlan) error {
	// TODO: not implemented yet

	logger := k.Logger(ctx)
	logger.Info("HandleRatioPlanProposal")
	// logger.Info("transferred from the community pool to recipient", "amount", p.Amount.String(), "recipient", p.Recipient)

	return nil
}

// ModifyFixedAmountPlanProposal is a handler for executing proposal modification.
func ModifyFixedAmountPlanProposal(ctx sdk.Context, k Keeper, plan types.FixedAmountPlan) error {
	// TODO: not implemented yet
	return nil
}

// ModifyRatioPlanProposal is a handler for executing proposal modification.
func ModifyRatioPlanProposal(ctx sdk.Context, k Keeper, plan types.RatioPlan) error {
	// TODO: not implemented yet
	return nil
}
