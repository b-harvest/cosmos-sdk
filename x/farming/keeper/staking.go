package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

//// NewStaking sets the index to a given staking
//func (k Keeper) NewStaking(ctx sdk.Context, staking types.Staking) types.Staking {
//	k.SetPlanIdByFarmerAddrIndex(ctx, staking.PlanId, staking.GetFarmerAddress())
//	return staking
//}

// GetStaking return a specific staking
func (k Keeper) GetStaking(ctx sdk.Context, planId uint64, farmerAcc sdk.AccAddress) (staking types.Staking, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingIndexKey(planId, farmerAcc))
	if bz == nil {
		return staking, false
	}
	k.cdc.MustUnmarshal(bz, &staking)
	return staking, true
}

// GetAllStakings returns all stakings in the Keeper.
func (k Keeper) GetAllStakings(ctx sdk.Context) (stakings []types.Staking) {
	k.IterateAllStakings(ctx, func(staking types.Staking) (stop bool) {
		stakings = append(stakings, staking)
		return false
	})

	return stakings
}

// SetStaking implements Staking.
func (k Keeper) SetStaking(ctx sdk.Context, staking types.Staking) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&staking)
	store.Set(types.GetStakingIndexKey(staking.PlanId, staking.GetFarmerAddress()), bz)
}

// RemoveStaking removes an staking for the staking mapper store.
func (k Keeper) RemoveStaking(ctx sdk.Context, staking types.Staking) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetStakingIndexKey(staking.PlanId, staking.GetFarmerAddress()))
}

// IterateAllStakings iterates over all the stored stakings and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllStakings(ctx sdk.Context, cb func(staking types.Staking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakingKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var staking types.Staking
		k.cdc.MustUnmarshal(iterator.Value(), &staking)
		if cb(staking) {
			break
		}
	}
}

// GetStakingsByPlanId reads from kvstore and return a specific Staking indexed by given plan id
func (k Keeper) GetStakingsByPlanId(ctx sdk.Context, planId uint64) (stakings []types.Staking) {
	k.IterateStakingsByPlanId(ctx, planId, func(staking types.Staking) bool {
		stakings = append(stakings, staking)
		return false
	})

	return stakings
}

// IterateAllStakings iterates over all the stored stakings and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateStakingsByPlanId(ctx sdk.Context, planId uint64, cb func(staking types.Staking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetStakingPrefix(planId))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var staking types.Staking
		k.cdc.MustUnmarshal(iterator.Value(), &staking)
		if cb(staking) {
			break
		}
	}
}

// TODO: WIP
func (k Keeper) Stake(ctx sdk.Context, msg *types.MsgStake) (types.Staking, error) {
	plan, found := k.GetPlan(ctx, msg.PlanId)
	if !found {
		return types.Staking{}, types.ErrPlanNotExists
	}
	farmerAcc, _ := sdk.AccAddressFromBech32(msg.Farmer)
	staking, found := k.GetStaking(ctx, plan.GetId(), farmerAcc)
	if !found {
		staking = types.Staking{
			PlanId:      plan.GetId(),
			Farmer:      msg.Farmer,
			StakedCoins: nil,
			QueuedCoins: msg.StakingCoins,
		}
		k.SetPlanIdByFarmerAddrIndex(ctx, staking.PlanId, staking.GetFarmerAddress())
	} else {
		staking.QueuedCoins = staking.QueuedCoins.Add(msg.StakingCoins...)
	}
	// TODO: add validation, check balance, stake to plan.StakingReserveAddress
	k.SetStaking(ctx, staking)
	return staking, nil
}

// TODO: WIP
func (k Keeper) Unstake(ctx sdk.Context, msg *types.MsgUnstake) (types.Staking, error) {
	return types.Staking{}, nil
}
