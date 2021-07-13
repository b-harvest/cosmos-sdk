package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// Simulation parameter constants
const (
	PrivateTaxCreationFee = "private_tax_creation_fee"
)

// GenPrivateTaxCreationFee return default PrivateTaxCreationFee
func GenPrivateTaxCreationFee(r *rand.Rand) sdk.Coins {
	// TODO: randomize private tax creation fee
	return types.DefaultPrivateTaxCreationFee
}

// RandomizedGenState generates a random GenesisState for tax
func RandomizedGenState(simState *module.SimulationState) {
	var privateTaxCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, PrivateTaxCreationFee, &privateTaxCreationFee, simState.Rand,
		func(r *rand.Rand) { privateTaxCreationFee = GenPrivateTaxCreationFee(r) },
	)

	taxGenesis := types.GenesisState{
		Params: types.Params{
			PrivateTaxCreationFee: privateTaxCreationFee,
		},
	}

	bz, _ := json.MarshalIndent(&taxGenesis, "", " ")
	fmt.Printf("Selected randomly generated tax parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&taxGenesis)
}
