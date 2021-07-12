package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

//// NewStaking sets the index to a given staking
//func (k Keeper) NewStaking(ctx sdk.Context, staking types.Staking) types.Staking {
//	k.SetPlanIDByFarmerAddrIndex(ctx, staking.PlanId, staking.GetFarmerAddress())
//	return staking
//}

// GetStaking returns a staking owned by the farmer for a given plan.
func (k Keeper) GetStaking(ctx sdk.Context, planID uint64, farmerAcc sdk.AccAddress) (staking types.Staking, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakingIndexKey(planID, farmerAcc))
	if bz == nil {
		return staking, false
	}
	k.cdc.MustUnmarshal(bz, &staking)
	return staking, true
}

// GetAllStakings returns all the stakings.
func (k Keeper) GetAllStakings(ctx sdk.Context) (stakings []types.Staking) {
	k.IterateAllStakings(ctx, func(staking types.Staking) (stop bool) {
		stakings = append(stakings, staking)
		return false
	})

	return stakings
}

// GetStakingsByPlanID returns all the stakings for a given plan.
func (k Keeper) GetStakingsByPlanID(ctx sdk.Context, planID uint64) (stakings []types.Staking) {
	k.IterateStakingsByPlanID(ctx, planID, func(staking types.Staking) bool {
		stakings = append(stakings, staking)
		return false
	})

	return stakings
}

// SetStaking sets a staking.
func (k Keeper) SetStaking(ctx sdk.Context, staking types.Staking) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&staking)
	store.Set(types.GetStakingIndexKey(staking.PlanId, staking.GetFarmerAddress()), bz)
}

// RemoveStaking removes a staking.
func (k Keeper) RemoveStaking(ctx sdk.Context, staking types.Staking) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetStakingIndexKey(staking.PlanId, staking.GetFarmerAddress()))
}

// IterateAllStakings iterates over all the stakings.
// Once the callback returns true, it stops the iteration.
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

// IterateStakingsByPlanID iterates over all the stakings associated with a given plan.
// Once the callback returns true, it stops the iteration.
func (k Keeper) IterateStakingsByPlanID(ctx sdk.Context, planID uint64, cb func(staking types.Staking) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetStakingPrefix(planID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var staking types.Staking
		k.cdc.MustUnmarshal(iterator.Value(), &staking)
		if cb(staking) {
			break
		}
	}
}

// UnmarshalStaking unmarshals a staking from bytes.
func (k Keeper) UnmarshalStaking(bz []byte) (types.Staking, error) {
	var staking types.Staking
	return staking, k.cdc.Unmarshal(bz, &staking)
}

// ReserveStakingCoins sends staking coins to the staking reserve account.
func (k Keeper) ReserveStakingCoins(ctx sdk.Context, farmer, reserveAcc sdk.AccAddress, stakingCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoins(ctx, farmer, reserveAcc, stakingCoins); err != nil {
		return err
	}
	return nil
}

// ReleaseStakingCoins sends staking coins back to the farmer.
func (k Keeper) ReleaseStakingCoins(ctx sdk.Context, reserveAcc, farmer sdk.AccAddress, unstakingCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoins(ctx, reserveAcc, farmer, unstakingCoins); err != nil {
		return err
	}
	return nil
}

// Stake stores staking coins to queued coins and it will be processed in the next epoch.
func (k Keeper) Stake(ctx sdk.Context, msg *types.MsgStake) (types.Staking, error) {
	plan, found := k.GetPlan(ctx, msg.PlanId)
	if !found {
		return types.Staking{}, types.ErrPlanNotExists
	}

	farmerAcc, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		return types.Staking{}, err
	}

	staking, found := k.GetStaking(ctx, plan.GetId(), farmerAcc)
	if !found {
		staking = types.Staking{
			PlanId:      plan.GetId(),
			Farmer:      msg.Farmer,
			StakedCoins: nil,
			QueuedCoins: msg.StakingCoins,
		}
		k.SetPlanIDByFarmerAddrIndex(ctx, staking.PlanId, staking.GetFarmerAddress())
	} else {
		staking.QueuedCoins = staking.QueuedCoins.Add(msg.StakingCoins...)
	}

	k.SetStaking(ctx, staking)
	k.ReserveStakingCoins(ctx, farmerAcc, plan.GetStakingReserveAddress(), staking.QueuedCoins)

	return staking, nil
}

// Unstake unstakes an amount of staking coins from the staking reserve account.
func (k Keeper) Unstake(ctx sdk.Context, msg *types.MsgUnstake) (types.Staking, error) {
	plan, found := k.GetPlan(ctx, msg.PlanId)
	if !found {
		return types.Staking{}, types.ErrPlanNotExists
	}

	farmerAcc, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		return types.Staking{}, err
	}

	staking, found := k.GetStaking(ctx, plan.GetId(), farmerAcc)
	if !found {
		return types.Staking{}, types.ErrStakingNotExists
	}

	// TODO: double check with this logic
	stakedDiff, hasNeg := staking.StakedCoins.SafeSub(msg.UnstakingCoins)
	if hasNeg {
		diff := stakedDiff.Add(staking.QueuedCoins...)
		if diff.IsAnyNegative() {
			return types.Staking{}, types.ErrInsufficientStakingAmount
		}
		staking.StakedCoins = sdk.Coins{}
		staking.QueuedCoins = diff
	}
	staking.StakedCoins = stakedDiff

	k.SetStaking(ctx, staking)
	k.ReleaseStakingCoins(ctx, plan.GetStakingReserveAddress(), farmerAcc, msg.UnstakingCoins)

	return types.Staking{}, nil
}
