package bls

import (
	"github.com/cosmos/cosmos-sdk/x/bls/keeper"
	"github.com/cosmos/cosmos-sdk/x/bls/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetGenBlsKeys(ctx, genState.GenesisKeys)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	return genesis
}
