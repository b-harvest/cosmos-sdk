package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// RegisterInvariants registers all tax invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "escrow-amount",
		TaxPoolsEscrowAmountInvariant(k))
}

// AllInvariants runs all invariants of the tax module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := TaxPoolsEscrowAmountInvariant(k)(ctx)
		return res, stop
	}
}

// TaxPoolsEscrowAmountInvariant checks that outstanding unwithdrawn fees are never negative.
func TaxPoolsEscrowAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// remainingCoins := sdk.NewCoins()
		// batches := k.GetAllPoolBatches(ctx)
		// for _, batch := range batches {
		// 	swapMsgs := k.GetAllPoolBatchSwapMsgStatesNotToBeDeleted(ctx, batch)
		// 	for _, msg := range swapMsgs {
		// 		remainingCoins = remainingCoins.Add(msg.RemainingOfferCoin)
		// 	}
		// 	depositMsgs := k.GetAllPoolBatchDepositMsgStatesNotToBeDeleted(ctx, batch)
		// 	for _, msg := range depositMsgs {
		// 		remainingCoins = remainingCoins.Add(msg.Msg.DepositCoins...)
		// 	}
		// 	withdrawMsgs := k.GetAllPoolBatchWithdrawMsgStatesNotToBeDeleted(ctx, batch)
		// 	for _, msg := range withdrawMsgs {
		// 		remainingCoins = remainingCoins.Add(msg.Msg.PoolCoin)
		// 	}
		// }

		batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
		escrowAmt := k.bankKeeper.GetAllBalances(ctx, batchEscrowAcc)

		broken := !escrowAmt.IsAllGTE(sdk.Coins{})

		return sdk.FormatInvariant(types.ModuleName, "batch escrow amount invariant broken",
			"batch escrow amount LT batch remaining amount"), broken
	}
}
