package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

// NewPlan sets the next plan number to a given plan interface
func (k Keeper) NewPlan(ctx sdk.Context, plan types.PlanI) types.PlanI {
	if err := plan.SetId(k.GetNextPlanIDWithUpdate(ctx)); err != nil {
		panic(err)
	}

	return plan
}

// GetPlan returns the plan for a given id.
func (k Keeper) GetPlan(ctx sdk.Context, id uint64) (plan types.PlanI, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPlanKey(id))
	if bz == nil {
		return plan, false
	}

	return k.decodePlan(bz), true
}

// GetAllPlans returns all the plans.
func (k Keeper) GetAllPlans(ctx sdk.Context) (plans []types.PlanI) {
	k.IterateAllPlans(ctx, func(plan types.PlanI) (stop bool) {
		plans = append(plans, plan)
		return false
	})

	return plans
}

// SetPlan sets a plan.
func (k Keeper) SetPlan(ctx sdk.Context, plan types.PlanI) {
	id := plan.GetId()
	store := ctx.KVStore(k.storeKey)

	bz, err := k.MarshalPlan(plan)
	if err != nil {
		panic(err)
	}

	store.Set(types.GetPlanKey(id), bz)
}

// RemovePlan removes a plan.
// NOTE: this will cause supply invariant violation if called
func (k Keeper) RemovePlan(ctx sdk.Context, plan types.PlanI) {
	id := plan.GetId()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPlanKey(id))
}

// IterateAllPlans iterates over all the plans.
// Once the callback returns true, it stops the iteration.
func (k Keeper) IterateAllPlans(ctx sdk.Context, cb func(plan types.PlanI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PlanKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		plan := k.decodePlan(iterator.Value())

		if cb(plan) {
			break
		}
	}
}

// GetPlansByFarmerAddrIndex returns all the plans the farmer is staking to.
func (k Keeper) GetPlansByFarmerAddrIndex(ctx sdk.Context, farmerAcc sdk.AccAddress) (plans []types.PlanI) {
	k.IteratePlansByFarmerAddr(ctx, farmerAcc, func(plan types.PlanI) bool {
		plans = append(plans, plan)
		return false
	})

	return plans
}

// IteratePlansByFarmerAddr iterates over all the plans.
// Once the callback returns true, it stops the iteration.
func (k Keeper) IteratePlansByFarmerAddr(ctx sdk.Context, farmerAcc sdk.AccAddress, cb func(plan types.PlanI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPlansByFarmerAddrIndexKey(farmerAcc))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		planID := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(iterator.Value(), &planID)
		if err != nil {
			panic(err)
		}
		plan, _ := k.GetPlan(ctx, planID.GetValue())
		if cb(plan) {
			break
		}
	}
}

// SetPlanIDByFarmerAddrIndex sets association within a plan and a farmer.
// TODO: need to gas cost check for existing check or update everytime
func (k Keeper) SetPlanIDByFarmerAddrIndex(ctx sdk.Context, planID uint64, farmerAcc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: planID})
	store.Set(types.GetPlanByFarmerAddrIndexKey(farmerAcc, planID), b)
}

// CreateFixedAmountPlan sets a fixed amount plan from the msg and a given plan type.
func (k Keeper) CreateFixedAmountPlan(ctx sdk.Context, msg *types.MsgCreateFixedAmountPlan, typ types.PlanType) *types.FixedAmountPlan {
	nextId := k.GetNextPlanIDWithUpdate(ctx)
	farmingPoolAddr := msg.GetFarmingPoolAddress()
	terminationAddr := farmingPoolAddr

	basePlan := types.NewBasePlan(
		nextId,
		typ,
		farmingPoolAddr,
		terminationAddr,
		msg.GetStakingCoinWeights(),
		msg.StartTime,
		msg.EndTime,
		msg.GetEpochDays(),
	)

	fixedPlan := types.NewFixedAmountPlan(basePlan, msg.EpochAmount)

	k.SetPlan(ctx, fixedPlan)

	return fixedPlan
}

// CreateRatioPlan sets a ratio plan from the msg and a given plan type.
func (k Keeper) CreateRatioPlan(ctx sdk.Context, msg *types.MsgCreateRatioPlan, typ types.PlanType) *types.RatioPlan {
	nextId := k.GetNextPlanIDWithUpdate(ctx)
	farmingPoolAddr := msg.GetFarmingPoolAddress()
	terminationAddr := farmingPoolAddr

	basePlan := types.NewBasePlan(
		nextId,
		typ,
		farmingPoolAddr,
		terminationAddr,
		msg.GetStakingCoinWeights(),
		msg.StartTime,
		msg.EndTime,
		msg.GetEpochDays(),
	)

	ratioPlan := types.NewRatioPlan(basePlan, msg.EpochRatio)

	k.SetPlan(ctx, ratioPlan)

	return ratioPlan
}
