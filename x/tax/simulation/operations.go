package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/tax/keeper"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateFixedAmountTax = "op_weight_msg_create_fixed_amount_tax"
	OpWeightMsgCreateRatioTax       = "op_weight_msg_create_ratio_tax"
	OpWeightMsgStake                 = "op_weight_msg_stake"
	OpWeightMsgUnstake               = "op_weight_msg_unstake"
	OpWeightMsgClaim                 = "op_weight_msg_claim"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreateFixedAmountTax int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateFixedAmountTax, &weightMsgCreateFixedAmountTax, nil,
		func(_ *rand.Rand) {
			weightMsgCreateFixedAmountTax = params.DefaultWeightMsgCreateFixedAmountTax
		},
	)

	var weightMsgCreateRatioTax int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateRatioTax, &weightMsgCreateRatioTax, nil,
		func(_ *rand.Rand) {
			weightMsgCreateRatioTax = params.DefaultWeightMsgCreateRatioTax
		},
	)

	var weightMsgStake int
	appParams.GetOrGenerate(cdc, OpWeightMsgStake, &weightMsgStake, nil,
		func(_ *rand.Rand) {
			weightMsgStake = params.DefaultWeightMsgStake
		},
	)

	var weightMsgUnstake int
	appParams.GetOrGenerate(cdc, OpWeightMsgUnstake, &weightMsgUnstake, nil,
		func(_ *rand.Rand) {
			weightMsgUnstake = params.DefaultWeightMsgUnstake
		},
	)

	var weightMsgClaim int
	appParams.GetOrGenerate(cdc, OpWeightMsgClaim, &weightMsgClaim, nil,
		func(_ *rand.Rand) {
			weightMsgClaim = params.DefaultWeightMsgClaim
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateFixedAmountTax,
			SimulateMsgCreateFixedAmountTax(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateRatioTax,
			SimulateMsgCreateRatioTax(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgStake,
			SimulateMsgStake(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUnstake,
			SimulateMsgUnstake(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgClaim,
			SimulateMsgClaim(ak, bk, k),
		),
	}
}

// SimulateMsgCreateFixedAmountTax generates a MsgCreateFixedAmountTax with random values
// nolint: interfacer
func SimulateMsgCreateFixedAmountTax(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}

// SimulateMsgCreateRatioTax generates a MsgCreateRatioTax with random values
// nolint: interfacer
func SimulateMsgCreateRatioTax(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}

// SimulateMsgStake generates a MsgCreateFixedAmountTax with random values
// nolint: interfacer
func SimulateMsgStake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}

// SimulateMsgUnstake generates a SimulateMsgUnstake with random values
// nolint: interfacer
func SimulateMsgUnstake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}

// SimulateMsgClaim generates a MsgClaim with random values
// nolint: interfacer
func SimulateMsgClaim(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: not implemented yet
		return simtypes.OperationMsg{}, nil, nil
	}
}
