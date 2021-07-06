package keeper

// TODO: unimplemented

//// GetPlan reads from kvstore and returns a specific plan
//func (k Keeper) GetPlan(ctx sdk.Context, planID uint64) (plan types.PlanI, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPlanKey(planID)
//
//	value := store.Get(key)
//	if value == nil {
//		return plan, false
//	}
//
//	plan = types.MustUnmarshalPlan(k.cdc, value)
//
//	return plan, true
//}
//
//// SetPlan sets to kvstore a specific plan
//func (k Keeper) SetPlan(ctx sdk.Context, plan types.PlanI) {
//	store := ctx.KVStore(k.storeKey)
//	b := types.MustMarshalPlan(k.cdc, plan)
//	store.Set(types.GetPlanKey(plan.Id), b)
//}
//
//// delete from kvstore a specific farmingPlan
//func (k Keeper) DeletePlan(ctx sdk.Context, plan types.PlanI) {
//	store := ctx.KVStore(k.storeKey)
//	Key := types.GetPlanKey(plan.Id)
//	store.Delete(Key)
//}
//
//// IterateAllPlans iterate through all of the farmingPlans
//func (k Keeper) IterateAllPlans(ctx sdk.Context, cb func(plan types.PlanI) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	iterator := sdk.KVStorePrefixIterator(store, types.PlanKeyPrefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		plan := types.MustUnmarshalPlan(k.cdc, iterator.Value())
//		if cb(plan) {
//			break
//		}
//	}
//}
//
//// GetAllPlans returns all farmingPlans used during genesis dump
//func (k Keeper) GetAllPlans(ctx sdk.Context) (plans []types.PlanI) {
//	k.IterateAllPlans(ctx, func(farmingPlan types.PlanI) bool {
//		plans = append(plans, farmingPlan)
//		return false
//	})
//
//	return plans
//}

//// GetNextPlanIDWithUpdate returns and increments the global Plan ID counter.
//// If the global account number is not set, it initializes it with value 0.
//func (k Keeper) GetNextPlanIDWithUpdate(ctx sdk.Context) uint64 {
//	store := ctx.KVStore(k.storeKey)
//	planID := k.GetNextPlanID(ctx)
//	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: planID + 1})
//	store.Set(types.GlobalFarmingPlanIDKey, bz)
//	return planID
//}

//// GetPlanByReserveAccIndex reads from kvstore and return a specific farmingPlan indexed by given reserve account
//func (k Keeper) GetPlanByReserveAccIndex(ctx sdk.Context, reserveAcc sdk.AccAddress) (plan types.Plan, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPlanByReserveAccIndexKey(reserveAcc)
//
//	value := store.Get(key)
//	if value == nil {
//		return plan, false
//	}
//
//	val := gogotypes.UInt64Value{}
//	err := k.cdc.UnmarshalBinaryBare(value, &val)
//	if err != nil {
//		return plan, false
//	}
//	planID := val.GetValue()
//	return k.GetPlan(ctx, planID)
//}
//
//// SetPlanByReserveAccIndex sets Index by ReserveAcc for plan duplication check
//func (k Keeper) SetPlanByReserveAccIndex(ctx sdk.Context, plan types.Plan) {
//	store := ctx.KVStore(k.storeKey)
//	b := k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: plan.Id})
//	store.Set(types.GetPlanByReserveAccIndexKey(plan.GetReserveAccount()), b)
//}
//
//// SetPlanAtomic sets plan with set global plan id index +1 and index by reserveAcc
//func (k Keeper) SetPlanAtomic(ctx sdk.Context, plan types.Plan) types.Plan {
//	plan.Id = k.GetNextPlanIDWithUpdate(ctx)
//	k.SetPlan(ctx, plan)
//	k.SetPlanByReserveAccIndex(ctx, plan)
//	return plan
//}
//
//// GetPlanBatchIndex returns the plan's latest batch index
//func (k Keeper) GetPlanBatchIndex(ctx sdk.Context, planID uint64) uint64 {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPlanBatchIndexKey(planID)
//
//	bz := store.Get(key)
//	if bz == nil {
//		return 0
//	}
//	return sdk.BigEndianToUint64(bz)
//}
//
//// SetPlanBatchIndex sets index for plan batch, it should be increase after batch executed
//func (k Keeper) SetPlanBatchIndex(ctx sdk.Context, planID, batchIndex uint64) {
//	store := ctx.KVStore(k.storeKey)
//	b := sdk.Uint64ToBigEndian(batchIndex)
//	store.Set(types.GetPlanBatchIndexKey(planID), b)
//}
//
//// GetPlanBatch returns a specific plan batch
//func (k Keeper) GetPlanBatch(ctx sdk.Context, planID uint64) (planBatch types.PlanBatch, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPlanBatchKey(planID)
//
//	value := store.Get(key)
//	if value == nil {
//		return planBatch, false
//	}
//
//	planBatch = types.MustUnmarshalPlanBatch(k.cdc, value)
//
//	return planBatch, true
//}
//
//// GetNextPlanBatchIndexWithUpdate returns next batch index, with set index increased
//func (k Keeper) GetNextPlanBatchIndexWithUpdate(ctx sdk.Context, planID uint64) (batchIndex uint64) {
//	batchIndex = k.GetPlanBatchIndex(ctx, planID)
//	batchIndex++
//	k.SetPlanBatchIndex(ctx, planID, batchIndex)
//	return
//}
//
//// GetAllPlanBatches returns all batches of the all existed farming plans
//func (k Keeper) GetAllPlanBatches(ctx sdk.Context) (planBatches []types.PlanBatch) {
//	k.IterateAllPlanBatches(ctx, func(planBatch types.PlanBatch) bool {
//		planBatches = append(planBatches, planBatch)
//		return false
//	})
//
//	return planBatches
//}
//
//// IterateAllPlanBatches iterate through all of the plan batches
//func (k Keeper) IterateAllPlanBatches(ctx sdk.Context, cb func(planBatch types.PlanBatch) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	iterator := sdk.KVStorePrefixIterator(store, types.PlanBatchKeyPrefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		planBatch := types.MustUnmarshalPlanBatch(k.cdc, iterator.Value())
//		if cb(planBatch) {
//			break
//		}
//	}
//}
//
//// DeletePlanBatch deletes batch of the plan, it used for test case
//func (k Keeper) DeletePlanBatch(ctx sdk.Context, planBatch types.PlanBatch) {
//	store := ctx.KVStore(k.storeKey)
//	batchKey := types.GetPlanBatchKey(planBatch.PlanId)
//	store.Delete(batchKey)
//}
//
//// SetPlanBatch sets batch of the plan, with current state
//func (k Keeper) SetPlanBatch(ctx sdk.Context, planBatch types.PlanBatch) {
//	store := ctx.KVStore(k.storeKey)
//	b := types.MustMarshalPlanBatch(k.cdc, planBatch)
//	store.Set(types.GetPlanBatchKey(planBatch.PlanId), b)
//}
//
//// GetPlanBatchDepositMsgState returns a specific DepositMsgState
//func (k Keeper) GetPlanBatchDepositMsgState(ctx sdk.Context, planID, msgIndex uint64) (state types.DepositMsgState, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPlanBatchDepositMsgStateIndexKey(planID, msgIndex)
//
//	value := store.Get(key)
//	if value == nil {
//		return state, false
//	}
//
//	state = types.MustUnmarshalDepositMsgState(k.cdc, value)
//	return state, true
//}
//
//// SetPlanBatchDepositMsgState sets deposit msg state of the plan batch, with current state
//func (k Keeper) SetPlanBatchDepositMsgState(ctx sdk.Context, planID uint64, state types.DepositMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	b := types.MustMarshalDepositMsgState(k.cdc, state)
//	store.Set(types.GetPlanBatchDepositMsgStateIndexKey(planID, state.MsgIndex), b)
//}
//
//// SetPlanBatchDepositMsgStatesByPointer sets deposit batch msgs of the plan batch, with current state using pointers
//func (k Keeper) SetPlanBatchDepositMsgStatesByPointer(ctx sdk.Context, planID uint64, states []*types.DepositMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	for _, state := range states {
//		if planID != state.Msg.PlanId {
//			continue
//		}
//		b := types.MustMarshalDepositMsgState(k.cdc, *state)
//		store.Set(types.GetPlanBatchDepositMsgStateIndexKey(planID, state.MsgIndex), b)
//	}
//}
//
//// SetPlanBatchDepositMsgStates sets deposit batch msgs of the plan batch, with current state
//func (k Keeper) SetPlanBatchDepositMsgStates(ctx sdk.Context, planID uint64, states []types.DepositMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	for _, state := range states {
//		if planID != state.Msg.PlanId {
//			continue
//		}
//		b := types.MustMarshalDepositMsgState(k.cdc, state)
//		store.Set(types.GetPlanBatchDepositMsgStateIndexKey(planID, state.MsgIndex), b)
//	}
//}
//
//// IterateAllPlanBatchDepositMsgStates iterate through all of the DepositMsgStates in the batch
//func (k Keeper) IterateAllPlanBatchDepositMsgStates(ctx sdk.Context, planBatch types.PlanBatch, cb func(state types.DepositMsgState) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	prefix := types.GetPlanBatchDepositMsgStatesPrefix(planBatch.PlanId)
//	iterator := sdk.KVStorePrefixIterator(store, prefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalDepositMsgState(k.cdc, iterator.Value())
//		if cb(state) {
//			break
//		}
//	}
//}
//
//// IterateAllDepositMsgStates iterate through all of the DepositMsgState of all batches
//func (k Keeper) IterateAllDepositMsgStates(ctx sdk.Context, cb func(state types.DepositMsgState) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	prefix := types.PlanBatchDepositMsgStateIndexKeyPrefix
//	iterator := sdk.KVStorePrefixIterator(store, prefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalDepositMsgState(k.cdc, iterator.Value())
//		if cb(state) {
//			break
//		}
//	}
//}
//
//// GetAllDepositMsgStates returns all BatchDepositMsgs for all batches.
//func (k Keeper) GetAllDepositMsgStates(ctx sdk.Context) (states []types.DepositMsgState) {
//	k.IterateAllDepositMsgStates(ctx, func(state types.DepositMsgState) bool {
//		states = append(states, state)
//		return false
//	})
//	return states
//}
//
//// GetAllPlanBatchDepositMsgs returns all BatchDepositMsgs indexed by the plan batch
//func (k Keeper) GetAllPlanBatchDepositMsgs(ctx sdk.Context, planBatch types.PlanBatch) (states []types.DepositMsgState) {
//	k.IterateAllPlanBatchDepositMsgStates(ctx, planBatch, func(state types.DepositMsgState) bool {
//		states = append(states, state)
//		return false
//	})
//	return states
//}
//
//// GetAllPlanBatchDepositMsgStatesNotToBeDeleted returns all Not toDelete BatchDepositMsgs indexed by the farmingPlanBatch
//func (k Keeper) GetAllPlanBatchDepositMsgStatesNotToBeDeleted(ctx sdk.Context, planBatch types.PlanBatch) (states []types.DepositMsgState) {
//	k.IterateAllPlanBatchDepositMsgStates(ctx, planBatch, func(state types.DepositMsgState) bool {
//		if !state.ToBeDeleted {
//			states = append(states, state)
//		}
//		return false
//	})
//	return states
//}
//
//// GetAllRemainingPlanBatchDepositMsgStates returns all remaining DepositMsgStates after endblock,
//// which are executed but not to be deleted
//func (k Keeper) GetAllRemainingPlanBatchDepositMsgStates(ctx sdk.Context, planBatch types.PlanBatch) (states []*types.DepositMsgState) {
//	k.IterateAllPlanBatchDepositMsgStates(ctx, planBatch, func(state types.DepositMsgState) bool {
//		if state.Executed && !state.ToBeDeleted {
//			states = append(states, &state)
//		}
//		return false
//	})
//	return states
//}
//
//// delete deposit batch msgs of the farming plan batch which has state ToBeDeleted
//func (k Keeper) DeleteAllReadyPlanBatchDepositMsgStates(ctx sdk.Context, planBatch types.PlanBatch) {
//	store := ctx.KVStore(k.storeKey)
//	iterator := sdk.KVStorePrefixIterator(store, types.GetPlanBatchDepositMsgStatesPrefix(planBatch.PlanId))
//	defer iterator.Close()
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalDepositMsgState(k.cdc, iterator.Value())
//		if state.ToBeDeleted {
//			store.Delete(iterator.Key())
//		}
//	}
//}
//
//// return a specific farmingPlanBatchWithdrawMsg
//func (k Keeper) GetPlanBatchWithdrawMsgState(ctx sdk.Context, planID, msgIndex uint64) (state types.WithdrawMsgState, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPlanBatchWithdrawMsgStateIndexKey(planID, msgIndex)
//
//	value := store.Get(key)
//	if value == nil {
//		return state, false
//	}
//
//	state = types.MustUnmarshalWithdrawMsgState(k.cdc, value)
//	return state, true
//}
//
//// set withdraw batch msg of the farming plan batch, with current state
//func (k Keeper) SetPlanBatchWithdrawMsgState(ctx sdk.Context, planID uint64, state types.WithdrawMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	b := types.MustMarshalWithdrawMsgState(k.cdc, state)
//	store.Set(types.GetPlanBatchWithdrawMsgStateIndexKey(planID, state.MsgIndex), b)
//}
//
//// set withdraw batch msgs of the farming plan batch, with current state using pointers
//func (k Keeper) SetPlanBatchWithdrawMsgStatesByPointer(ctx sdk.Context, planID uint64, states []*types.WithdrawMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	for _, state := range states {
//		if planID != state.Msg.PlanId {
//			continue
//		}
//		b := types.MustMarshalWithdrawMsgState(k.cdc, *state)
//		store.Set(types.GetPlanBatchWithdrawMsgStateIndexKey(planID, state.MsgIndex), b)
//	}
//}
//
//// set withdraw batch msgs of the plan batch, with current state
//func (k Keeper) SetPlanBatchWithdrawMsgStates(ctx sdk.Context, planID uint64, states []types.WithdrawMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	for _, state := range states {
//		if planID != state.Msg.PlanId {
//			continue
//		}
//		b := types.MustMarshalWithdrawMsgState(k.cdc, state)
//		store.Set(types.GetPlanBatchWithdrawMsgStateIndexKey(planID, state.MsgIndex), b)
//	}
//}
//
//// IterateAllPlanBatchWithdrawMsgStates iterate through all of the FarmingPlanBatchWithdrawMsgs
//func (k Keeper) IterateAllPlanBatchWithdrawMsgStates(ctx sdk.Context, planBatch types.PlanBatch, cb func(state types.WithdrawMsgState) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	prefix := types.GetPlanBatchWithdrawMsgsPrefix(planBatch.PlanId)
//	iterator := sdk.KVStorePrefixIterator(store, prefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalWithdrawMsgState(k.cdc, iterator.Value())
//		if cb(state) {
//			break
//		}
//	}
//}
//
//// IterateAllWithdrawMsgStates iterate through all of the WithdrawMsgState of all batches
//func (k Keeper) IterateAllWithdrawMsgStates(ctx sdk.Context, cb func(state types.WithdrawMsgState) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	prefix := types.PlanBatchWithdrawMsgStateIndexKeyPrefix
//	iterator := sdk.KVStorePrefixIterator(store, prefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalWithdrawMsgState(k.cdc, iterator.Value())
//		if cb(state) {
//			break
//		}
//	}
//}
//
//// GetAllWithdrawMsgStates returns all BatchWithdrawMsgs for all batches
//func (k Keeper) GetAllWithdrawMsgStates(ctx sdk.Context) (states []types.WithdrawMsgState) {
//	k.IterateAllWithdrawMsgStates(ctx, func(state types.WithdrawMsgState) bool {
//		states = append(states, state)
//		return false
//	})
//	return states
//}
//
//// GetAllPlanBatchWithdrawMsgStates returns all BatchWithdrawMsgs indexed by the farmingPlanBatch
//func (k Keeper) GetAllPlanBatchWithdrawMsgStates(ctx sdk.Context, planBatch types.PlanBatch) (states []types.WithdrawMsgState) {
//	k.IterateAllPlanBatchWithdrawMsgStates(ctx, planBatch, func(state types.WithdrawMsgState) bool {
//		states = append(states, state)
//		return false
//	})
//	return states
//}
//
//// GetAllPlanBatchWithdrawMsgStatesNotToBeDeleted returns all Not to delete BatchWithdrawMsgs indexed by the farmingPlanBatch
//func (k Keeper) GetAllPlanBatchWithdrawMsgStatesNotToBeDeleted(ctx sdk.Context, planBatch types.PlanBatch) (states []types.WithdrawMsgState) {
//	k.IterateAllPlanBatchWithdrawMsgStates(ctx, planBatch, func(state types.WithdrawMsgState) bool {
//		if !state.ToBeDeleted {
//			states = append(states, state)
//		}
//		return false
//	})
//	return states
//}
//
//// GetAllRemainingPlanBatchWithdrawMsgStates returns All only remaining BatchWithdrawMsgs after endblock, executed but not toDelete
//func (k Keeper) GetAllRemainingPlanBatchWithdrawMsgStates(ctx sdk.Context, planBatch types.PlanBatch) (states []*types.WithdrawMsgState) {
//	k.IterateAllPlanBatchWithdrawMsgStates(ctx, planBatch, func(state types.WithdrawMsgState) bool {
//		if state.Executed && !state.ToBeDeleted {
//			states = append(states, &state)
//		}
//		return false
//	})
//	return states
//}
//
//// delete withdraw batch msgs of the farming plan batch which has state ToBeDeleted
//func (k Keeper) DeleteAllReadyPlanBatchWithdrawMsgStates(ctx sdk.Context, planBatch types.PlanBatch) {
//	store := ctx.KVStore(k.storeKey)
//	iterator := sdk.KVStorePrefixIterator(store, types.GetPlanBatchWithdrawMsgsPrefix(planBatch.PlanId))
//	defer iterator.Close()
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalWithdrawMsgState(k.cdc, iterator.Value())
//		if state.ToBeDeleted {
//			store.Delete(iterator.Key())
//		}
//	}
//}
//
//// return a specific SwapMsgState given the plan_id with the msg_index
//func (k Keeper) GetPlanBatchSwapMsgState(ctx sdk.Context, planID, msgIndex uint64) (state types.SwapMsgState, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPlanBatchSwapMsgStateIndexKey(planID, msgIndex)
//
//	value := store.Get(key)
//	if value == nil {
//		return state, false
//	}
//
//	state = types.MustUnmarshalSwapMsgState(k.cdc, value)
//	return state, true
//}
//
//// set swap batch msg of the farming plan batch, with current state
//func (k Keeper) SetPlanBatchSwapMsgState(ctx sdk.Context, planID uint64, state types.SwapMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	b := types.MustMarshalSwapMsgState(k.cdc, state)
//	store.Set(types.GetPlanBatchSwapMsgStateIndexKey(planID, state.MsgIndex), b)
//}
//
//// Delete swap batch msg of the farming plan batch, it used for test case
//func (k Keeper) DeletePlanBatchSwapMsgState(ctx sdk.Context, planID uint64, msgIndex uint64) {
//	store := ctx.KVStore(k.storeKey)
//	batchKey := types.GetPlanBatchSwapMsgStateIndexKey(planID, msgIndex)
//	store.Delete(batchKey)
//}
//
//// IterateAllPlanBatchSwapMsgStates iterate through all of the FarmingPlanBatchSwapMsgs
//func (k Keeper) IterateAllPlanBatchSwapMsgStates(ctx sdk.Context, planBatch types.PlanBatch, cb func(state types.SwapMsgState) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	prefix := types.GetPlanBatchSwapMsgStatesPrefix(planBatch.PlanId)
//	iterator := sdk.KVStorePrefixIterator(store, prefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalSwapMsgState(k.cdc, iterator.Value())
//		if cb(state) {
//			break
//		}
//	}
//}
//
//// IterateAllSwapMsgStates iterate through all of the SwapMsgState of all batches
//func (k Keeper) IterateAllSwapMsgStates(ctx sdk.Context, cb func(state types.SwapMsgState) (stop bool)) {
//	store := ctx.KVStore(k.storeKey)
//
//	prefix := types.PlanBatchSwapMsgStateIndexKeyPrefix
//	iterator := sdk.KVStorePrefixIterator(store, prefix)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalSwapMsgState(k.cdc, iterator.Value())
//		if cb(state) {
//			break
//		}
//	}
//}
//
//// GetAllSwapMsgStates returns all BatchSwapMsgs of all batches
//func (k Keeper) GetAllSwapMsgStates(ctx sdk.Context) (states []types.SwapMsgState) {
//	k.IterateAllSwapMsgStates(ctx, func(state types.SwapMsgState) bool {
//		states = append(states, state)
//		return false
//	})
//	return states
//}
//
//// delete swap batch msgs of the farming plan batch which has state ToBeDeleted
//func (k Keeper) DeleteAllReadyPlanBatchSwapMsgStates(ctx sdk.Context, planBatch types.PlanBatch) {
//	store := ctx.KVStore(k.storeKey)
//	iterator := sdk.KVStorePrefixIterator(store, types.GetPlanBatchSwapMsgStatesPrefix(planBatch.PlanId))
//	defer iterator.Close()
//	for ; iterator.Valid(); iterator.Next() {
//		state := types.MustUnmarshalSwapMsgState(k.cdc, iterator.Value())
//		if state.ToBeDeleted {
//			store.Delete(iterator.Key())
//		}
//	}
//}
//
//// GetAllPlanBatchSwapMsgStatesAsPointer returns all BatchSwapMsgs pointer indexed by the farmingPlanBatch
//func (k Keeper) GetAllPlanBatchSwapMsgStatesAsPointer(ctx sdk.Context, planBatch types.PlanBatch) (states []*types.SwapMsgState) {
//	k.IterateAllPlanBatchSwapMsgStates(ctx, planBatch, func(state types.SwapMsgState) bool {
//		states = append(states, &state)
//		return false
//	})
//	return states
//}
//
//// GetAllPlanBatchSwapMsgStates returns all BatchSwapMsgs indexed by the farmingPlanBatch
//func (k Keeper) GetAllPlanBatchSwapMsgStates(ctx sdk.Context, planBatch types.PlanBatch) (states []types.SwapMsgState) {
//	k.IterateAllPlanBatchSwapMsgStates(ctx, planBatch, func(state types.SwapMsgState) bool {
//		states = append(states, state)
//		return false
//	})
//	return states
//}
//
//// GetAllNotProcessedPlanBatchSwapMsgStates returns All only not processed swap msgs, not executed with not succeed and not toDelete BatchSwapMsgs indexed by the farmingPlanBatch
//func (k Keeper) GetAllNotProcessedPlanBatchSwapMsgStates(ctx sdk.Context, planBatch types.PlanBatch) (states []*types.SwapMsgState) {
//	k.IterateAllPlanBatchSwapMsgStates(ctx, planBatch, func(state types.SwapMsgState) bool {
//		if !state.Executed && !state.Succeeded && !state.ToBeDeleted {
//			states = append(states, &state)
//		}
//		return false
//	})
//	return states
//}
//
//// GetAllRemainingPlanBatchSwapMsgStates returns All only remaining after endblock swap msgs, executed but not toDelete
//func (k Keeper) GetAllRemainingPlanBatchSwapMsgStates(ctx sdk.Context, planBatch types.PlanBatch) (states []*types.SwapMsgState) {
//	k.IterateAllPlanBatchSwapMsgStates(ctx, planBatch, func(state types.SwapMsgState) bool {
//		if state.Executed && !state.ToBeDeleted {
//			states = append(states, &state)
//		}
//		return false
//	})
//	return states
//}
//
//// GetAllPlanBatchSwapMsgStatesNotToBeDeleted returns All only not to delete swap msgs
//func (k Keeper) GetAllPlanBatchSwapMsgStatesNotToBeDeleted(ctx sdk.Context, planBatch types.PlanBatch) (states []*types.SwapMsgState) {
//	k.IterateAllPlanBatchSwapMsgStates(ctx, planBatch, func(state types.SwapMsgState) bool {
//		if !state.ToBeDeleted {
//			states = append(states, &state)
//		}
//		return false
//	})
//	return states
//}
//
//// set swap batch msgs of the farming plan batch, with current state using pointers
//func (k Keeper) SetPlanBatchSwapMsgStatesByPointer(ctx sdk.Context, planID uint64, states []*types.SwapMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	for _, state := range states {
//		if planID != state.Msg.PlanId {
//			continue
//		}
//		b := types.MustMarshalSwapMsgState(k.cdc, *state)
//		store.Set(types.GetPlanBatchSwapMsgStateIndexKey(planID, state.MsgIndex), b)
//	}
//}
//
//// set swap batch msgs of the farming plan batch, with current state
//func (k Keeper) SetPlanBatchSwapMsgStates(ctx sdk.Context, planID uint64, states []types.SwapMsgState) {
//	store := ctx.KVStore(k.storeKey)
//	for _, state := range states {
//		if planID != state.Msg.PlanId {
//			continue
//		}
//		b := types.MustMarshalSwapMsgState(k.cdc, state)
//		store.Set(types.GetPlanBatchSwapMsgStateIndexKey(planID, state.MsgIndex), b)
//	}
//}
