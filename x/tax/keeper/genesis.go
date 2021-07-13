package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// InitGenesis initializes the tax module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := k.ValidateGenesis(ctx, genState); err != nil {
		panic(err)
	}

	k.SetParams(ctx, genState.Params)

	// TODO: unimplemented
	//for _, record := range genState.TaxRecords {
	//	k.SetTaxRecord(ctx, record)
	//}
	//for _, staking := range genState.Stakings {
	//	k.SetStaking(ctx, staking)
	//}
	//for _, reward := range genState.Rewards {
	//	k.SetReward(ctx, reword)
	//}
}

// ExportGenesis returns the tax module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)


	//taxes := k.GetAllTaxes(ctx)
	//stakings := k.GetAllStakings(ctx)
	//rewards := k.GetAllRewards(ctx)

	//for _, tax := range taxes {
	//	record, found := k.GetTaxRecord(ctx, tax)
	//	if found {
	//		taxRecords = append(taxRecords, record)
	//	}
	//}
	//
	//if len(taxRecords) == 0 {
	//	taxRecords = []types.TaxRecord{}
	//}

	return types.NewGenesisState(params, taxRecords, nil, nil)
}

// ValidateGenesis validates the tax module's genesis state.
func (k Keeper) ValidateGenesis(ctx sdk.Context, genState types.GenesisState) error {
	if err := genState.Params.Validate(); err != nil {
		return err
	}

	cc, _ := ctx.CacheContext()
	k.SetParams(cc, genState.Params)

	// TODO: unimplemented
	//for _, record := range genState.TaxRecords {
	//	record = k.SetTaxRecord(cc, record)
	//	if err := k.ValidateTaxRecord(cc, record); err != nil {
	//		return err
	//	}
	//}

	return nil
}
