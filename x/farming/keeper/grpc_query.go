package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	panic("implement me")
}

func (q Querier) Plans(ctx context.Context, request *types.QueryPlansRequest) (*types.QueryPlansResponse, error) {
	panic("implement me")
}

func (q Querier) Plan(ctx context.Context, request *types.QueryPlanRequest) (*types.QueryPlanResponse, error) {
	panic("implement me")
}

func (q Querier) PlanStakings(ctx context.Context, request *types.QueryPlanStakingsRequest) (*types.QueryPlanStakingsResponse, error) {
	panic("implement me")
}

func (q Querier) FarmerStakings(ctx context.Context, request *types.QueryFarmerStakingsRequest) (*types.QueryFarmerStakingsResponse, error) {
	panic("implement me")
}

func (q Querier) FarmerStaking(ctx context.Context, request *types.QueryFarmerStakingRequest) (*types.QueryFarmerStakingResponse, error) {
	panic("implement me")
}

func (q Querier) PlanRewards(ctx context.Context, request *types.QueryPlanRewardsRequest) (*types.QueryPlanRewardsResponse, error) {
	panic("implement me")
}

func (q Querier) FarmerRewards(ctx context.Context, request *types.QueryFarmerRewardsRequest) (*types.QueryFarmerRewardsResponse, error) {
	panic("implement me")
}

func (q Querier) FarmerReward(ctx context.Context, request *types.QueryFarmerRewardRequest) (*types.QueryFarmerRewardResponse, error) {
	panic("implement me")
}
