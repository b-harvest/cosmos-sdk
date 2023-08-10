package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Simulation parameter constants
const (
	unbondingTime     = "unbonding_time"
	maxValidators     = "max_validators"
	historicalEntries = "historical_entries"
)

// genUnbondingTime returns randomized UnbondingTime
func genUnbondingTime(r *rand.Rand) (ubdTime time.Duration) {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

// genMaxValidators returns randomized MaxValidators
func genMaxValidators(r *rand.Rand) (maxValidators uint32) {
	return uint32(r.Intn(250) + 1)
}

// getHistEntries returns randomized HistoricalEntries between 0-100.
func getHistEntries(r *rand.Rand) uint32 {
	return uint32(r.Intn(int(types.DefaultHistoricalEntries + 1)))
}

// RandomizedGenState generates a random GenesisState for staking
func RandomizedGenState(simState *module.SimulationState) {
	params := types.NewParams(
		types.DefaultUnbondingTime, 100, 7, 10000, sdk.DefaultBondDenom,
	)

	// validators & delegations
	var (
		validators  []types.Validator
		delegations []types.Delegation
	)

	valAddrs := make([]sdk.ValAddress, simState.NumBonded)

	for i := 0; i < int(simState.NumBonded); i++ {
		valAddr := sdk.ValAddress(simState.Accounts[i].Address)
		valAddrs[i] = valAddr

		maxCommission := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(simState.Rand, 1, 100)), 2)
		commission := types.NewCommission(
			simulation.RandomDecAmount(simState.Rand, maxCommission),
			maxCommission,
			simulation.RandomDecAmount(simState.Rand, maxCommission),
		)

		validator, err := types.NewValidator(valAddr, simState.Accounts[i].ConsKey.PubKey(), types.Description{})
		if err != nil {
			panic(err)
		}
		validator.Tokens = simState.InitialStake
		validator.DelegatorShares = simState.InitialStake.ToDec()
		validator.Commission = commission

		delegation := types.NewDelegation(simState.Accounts[i].Address, valAddr, simState.InitialStake.ToDec())

		validators = append(validators, validator)
		delegations = append(delegations, delegation)
	}

	stakingGenesis := types.NewGenesisState(params, validators, delegations)

	bz, err := json.MarshalIndent(&stakingGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated staking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(stakingGenesis)
}
