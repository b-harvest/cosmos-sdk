package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Querier) Taxes(c context.Context, req *types.QueryTaxesRequest) (*types.QueryTaxesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.TaxSourceAddress != "" {
		if _, err := sdk.AccAddressFromBech32(req.TaxSourceAddress); err != nil {
			return nil, err
		}
	}

	if req.PoolAddress != "" {
		if _, err := sdk.AccAddressFromBech32(req.PoolAddress); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	taxStore := prefix.NewStore(store, types.TaxKeyPrefix)
	params := k.GetParams(ctx)

	var taxes []*types.Tax

	pageRes, err := query.FilteredPaginate(taxStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		tax, err := k.UnmarshalTax(value)
		if err != nil {
			return false, err
		}
		if err != nil {
			return false, err
		}

		if accumulate {
			taxes = append(taxes, tax)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTaxesResponse{Taxes: taxes, Pagination: pageRes}, nil
}
