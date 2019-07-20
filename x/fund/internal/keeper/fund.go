package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/fund/internal/types"
)

// SetProposal set a proposal to store
//func (keeper Keeper) SetProposal(ctx sdk.Context, proposal Proposal) {
//	store := ctx.KVStore(keeper.storeKey)
//	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal)
//	store.Set(ProposalKey(proposal.ProposalID), bz)
//}

// set a fund
func (k Keeper) SetFund(ctx sdk.Context, fund types.Fund) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalFund(k.cdc, fund)
	store.Set(types.GetFundKey(fund.FundID), b)
}

// return a specific fund
func (k Keeper) GetFund(ctx sdk.Context, fundID uint64) (
	fund types.Fund, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetFundKey(fundID)
	value := store.Get(key)
	if value == nil {
		return fund, false
	}

	fund = types.MustUnmarshalFund(k.cdc, value)
	return fund, true
}

// IterateAllFunds iterate through all of the delegations
func (k Keeper) IterateAllFunds(ctx sdk.Context, cb func(fund types.Fund) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FundKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		fund := types.MustUnmarshalFund(k.cdc, iterator.Value())
		if cb(fund) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump
func (k Keeper) GetAllFunds(ctx sdk.Context) (funds []types.Fund) {
	k.IterateAllFunds(ctx, func(fund types.Fund) bool {
		funds = append(funds, fund)
		return false // TODO: need to check logic
	})
	return funds
}

//TODO

//// GetNextAccountNumber Returns and increments the global account number counter
//func (ak AccountKeeper) GetNextFundId(ctx sdk.Context) uint64 {
//	var accNumber uint64
//	store := ctx.KVStore(ak.key)
//	bz := store.Get(types.GlobalAccountNumberKey)
//	if bz == nil {
//		accNumber = 0
//	} else {
//		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
//	store.Set(types.GlobalAccountNumberKey, bz)
//
//	return accNumber
//}
