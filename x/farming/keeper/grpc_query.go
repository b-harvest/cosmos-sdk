package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

//// FarmingPool queries a farming pool with the given pool id.
//func (k Querier) FarmingPool(c context.Context, req *types.QueryFarmingPoolRequest) (*types.QueryFarmingPoolResponse, error) {
//	empty := &types.QueryFarmingPoolRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	pool, found := k.GetPool(ctx, req.PoolId)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "farming pool %d doesn't exist", req.PoolId)
//	}
//
//	return k.MakeQueryFarmingPoolResponse(pool)
//}
//
//// FarmingPool queries a farming pool with the given pool coin denom.
//func (k Querier) FarmingPoolByPoolCoinDenom(c context.Context, req *types.QueryFarmingPoolByPoolCoinDenomRequest) (*types.QueryFarmingPoolResponse, error) {
//	empty := &types.QueryFarmingPoolByPoolCoinDenomRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//	ctx := sdk.UnwrapSDKContext(c)
//	reserveAcc, err := types.GetReserveAcc(req.PoolCoinDenom)
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "farming pool with pool coin denom %s doesn't exist", req.PoolCoinDenom)
//	}
//	pool, found := k.GetPoolByReserveAccIndex(ctx, reserveAcc)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "farming pool with pool coin denom %s doesn't exist", req.PoolCoinDenom)
//	}
//	return k.MakeQueryFarmingPoolResponse(pool)
//}
//
//// FarmingPool queries a farming pool with the given reserve account address.
//func (k Querier) FarmingPoolByReserveAcc(c context.Context, req *types.QueryFarmingPoolByReserveAccRequest) (*types.QueryFarmingPoolResponse, error) {
//	empty := &types.QueryFarmingPoolByReserveAccRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//	ctx := sdk.UnwrapSDKContext(c)
//	reserveAcc, err := sdk.AccAddressFromBech32(req.ReserveAcc)
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "the reserve account address %s is not valid", req.ReserveAcc)
//	}
//	pool, found := k.GetPoolByReserveAccIndex(ctx, reserveAcc)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "farming pool with pool reserve account %s doesn't exist", req.ReserveAcc)
//	}
//	return k.MakeQueryFarmingPoolResponse(pool)
//}
//
//// FarmingPoolBatch queries a farming pool batch with the given pool id.
//func (k Querier) FarmingPoolBatch(c context.Context, req *types.QueryFarmingPoolBatchRequest) (*types.QueryFarmingPoolBatchResponse, error) {
//	empty := &types.QueryFarmingPoolBatchRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	batch, found := k.GetPoolBatch(ctx, req.PoolId)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "farming pool batch %d doesn't exist", req.PoolId)
//	}
//
//	return &types.QueryFarmingPoolBatchResponse{
//		Batch: batch,
//	}, nil
//}
//
//// Pools queries all farming pools currently existed with each farming pool with batch and metadata.
//func (k Querier) FarmingPools(c context.Context, req *types.QueryFarmingPoolsRequest) (*types.QueryFarmingPoolsResponse, error) {
//	ctx := sdk.UnwrapSDKContext(c)
//
//	store := ctx.KVStore(k.storeKey)
//	poolStore := prefix.NewStore(store, types.PoolKeyPrefix)
//
//	var pools types.Pools
//
//	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
//		pool, err := types.UnmarshalPool(k.cdc, value)
//		if err != nil {
//			return err
//		}
//		pools = append(pools, pool)
//		return nil
//	})
//
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	if len(pools) == 0 {
//		return nil, status.Error(codes.NotFound, "There are no pools present.")
//	}
//
//	return &types.QueryFarmingPoolsResponse{
//		Pools:      pools,
//		Pagination: pageRes,
//	}, nil
//}
//
//// PoolBatchSwapMsg queries the pool batch swap message with the message index of the farming pool.
//func (k Querier) PoolBatchSwapMsg(c context.Context, req *types.QueryPoolBatchSwapMsgRequest) (*types.QueryPoolBatchSwapMsgResponse, error) {
//	empty := &types.QueryPoolBatchSwapMsgRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	msg, found := k.GetPoolBatchSwapMsgState(ctx, req.PoolId, req.MsgIndex)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "the msg given msg_index %d doesn't exist or deleted", req.MsgIndex)
//	}
//
//	return &types.QueryPoolBatchSwapMsgResponse{
//		Swap: msg,
//	}, nil
//}
//
//// PoolBatchSwapMsgs queries all pool batch swap messages of the farming pool.
//func (k Querier) PoolBatchSwapMsgs(c context.Context, req *types.QueryPoolBatchSwapMsgsRequest) (*types.QueryPoolBatchSwapMsgsResponse, error) {
//	empty := &types.QueryPoolBatchSwapMsgsRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	_, found := k.GetPool(ctx, req.PoolId)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "farming pool %d doesn't exist", req.PoolId)
//	}
//
//	store := ctx.KVStore(k.storeKey)
//	msgStore := prefix.NewStore(store, types.GetPoolBatchSwapMsgStatesPrefix(req.PoolId))
//
//	var msgs []types.SwapMsgState
//
//	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
//		msg, err := types.UnmarshalSwapMsgState(k.cdc, value)
//		if err != nil {
//			return err
//		}
//
//		msgs = append(msgs, msg)
//
//		return nil
//	})
//
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &types.QueryPoolBatchSwapMsgsResponse{
//		Swaps:      msgs,
//		Pagination: pageRes,
//	}, nil
//}
//
//// PoolBatchDepositMsg queries the pool batch deposit message with the msg_index of the farming pool.
//func (k Querier) PoolBatchDepositMsg(c context.Context, req *types.QueryPoolBatchDepositMsgRequest) (*types.QueryPoolBatchDepositMsgResponse, error) {
//	empty := &types.QueryPoolBatchDepositMsgRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	msg, found := k.GetPoolBatchDepositMsgState(ctx, req.PoolId, req.MsgIndex)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "the msg given msg_index %d doesn't exist or deleted", req.MsgIndex)
//	}
//
//	return &types.QueryPoolBatchDepositMsgResponse{
//		Deposit: msg,
//	}, nil
//}
//
//// PoolBatchDepositMsgs queries all pool batch deposit messages of the farming pool.
//func (k Querier) PoolBatchDepositMsgs(c context.Context, req *types.QueryPoolBatchDepositMsgsRequest) (*types.QueryPoolBatchDepositMsgsResponse, error) {
//	empty := &types.QueryPoolBatchDepositMsgsRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	_, found := k.GetPool(ctx, req.PoolId)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "farming pool %d doesn't exist", req.PoolId)
//	}
//
//	store := ctx.KVStore(k.storeKey)
//	msgStore := prefix.NewStore(store, types.GetPoolBatchDepositMsgStatesPrefix(req.PoolId))
//	var msgs []types.DepositMsgState
//
//	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
//		msg, err := types.UnmarshalDepositMsgState(k.cdc, value)
//		if err != nil {
//			return err
//		}
//
//		msgs = append(msgs, msg)
//
//		return nil
//	})
//
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &types.QueryPoolBatchDepositMsgsResponse{
//		Deposits:   msgs,
//		Pagination: pageRes,
//	}, nil
//}
//
//// PoolBatchWithdrawMsg queries the pool batch withdraw message with the msg_index of the farming pool.
//func (k Querier) PoolBatchWithdrawMsg(c context.Context, req *types.QueryPoolBatchWithdrawMsgRequest) (*types.QueryPoolBatchWithdrawMsgResponse, error) {
//	empty := &types.QueryPoolBatchWithdrawMsgRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	msg, found := k.GetPoolBatchWithdrawMsgState(ctx, req.PoolId, req.MsgIndex)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "the msg given msg_index %d doesn't exist or deleted", req.MsgIndex)
//	}
//
//	return &types.QueryPoolBatchWithdrawMsgResponse{
//		Withdraw: msg,
//	}, nil
//}
//
//// PoolBatchWithdrawMsgs queries all pool batch withdraw messages of the farming pool.
//func (k Querier) PoolBatchWithdrawMsgs(c context.Context, req *types.QueryPoolBatchWithdrawMsgsRequest) (*types.QueryPoolBatchWithdrawMsgsResponse, error) {
//	empty := &types.QueryPoolBatchWithdrawMsgsRequest{}
//	if req == nil || *req == *empty {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	_, found := k.GetPool(ctx, req.PoolId)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "farming pool %d doesn't exist", req.PoolId)
//	}
//
//	store := ctx.KVStore(k.storeKey)
//	msgStore := prefix.NewStore(store, types.GetPoolBatchWithdrawMsgsPrefix(req.PoolId))
//	var msgs []types.WithdrawMsgState
//
//	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
//		msg, err := types.UnmarshalWithdrawMsgState(k.cdc, value)
//		if err != nil {
//			return err
//		}
//
//		msgs = append(msgs, msg)
//
//		return nil
//	})
//
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &types.QueryPoolBatchWithdrawMsgsResponse{
//		Withdraws:  msgs,
//		Pagination: pageRes,
//	}, nil
//}
//
//// Params queries params of farming module.
//func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
//	ctx := sdk.UnwrapSDKContext(c)
//	params := k.GetParams(ctx)
//
//	return &types.QueryParamsResponse{
//		Params: params,
//	}, nil
//}
//
//// MakeQueryFarmingPoolResponse wraps MakeQueryFarmingPoolResponse.
//func (k Querier) MakeQueryFarmingPoolResponse(pool types.Pool) (*types.QueryFarmingPoolResponse, error) {
//	return &types.QueryFarmingPoolResponse{
//		Pool: pool,
//	}, nil
//}
//
//// MakeQueryFarmingPoolsResponse wraps a list of QueryFarmingPoolResponses.
//func (k Querier) MakeQueryFarmingPoolsResponse(pools types.Pools) (*[]types.QueryFarmingPoolResponse, error) {
//	resp := make([]types.QueryFarmingPoolResponse, len(pools))
//	for i, pool := range pools {
//		res := types.QueryFarmingPoolResponse{
//			Pool: pool,
//		}
//		resp[i] = res
//	}
//	return &resp, nil
//}
