package keeper_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/babylonchain/babylon/testutil/datagen"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/epoching/testepoching"
	"github.com/cosmos/cosmos-sdk/x/epoching/types"
)

func FuzzSlashedValSet(f *testing.F) {
	datagen.AddRandomSeedsToFuzzer(f, 10)
	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))

		helper := testepoching.NewHelperWithValSet(t)
		ctx, keeper, stakingKeeper := helper.Ctx, helper.EpochingKeeper, helper.StakingKeeper
		getValSet := keeper.GetValidatorSet(ctx, 0)

		// slash a random subset of validators
		numSlashed := r.Intn(len(getValSet))
		excpectedSlashedVals := []sdk.ValAddress{}
		for i := 0; i < numSlashed; i++ {
			idx := r.Intn(len(getValSet))
			slashedVal := getValSet[idx]
			stakingKeeper.Slash(ctx, sdk.ConsAddress(slashedVal.Addr), 0, slashedVal.Power, sdk.OneDec())
			// add the slashed validator to the slashed validator set
			excpectedSlashedVals = append(excpectedSlashedVals, slashedVal.Addr)
			// remove the slashed validator from the validator set in order to avoid slashing a validator more than once
			getValSet = append(getValSet[:idx], getValSet[idx+1:]...)
		}

		// check whether the slashed validator set in DB is consistent or not
		actualSlashedVals := keeper.GetSlashedValidators(ctx, 0)
		require.Equal(t, len(excpectedSlashedVals), len(actualSlashedVals))
		sortVals(excpectedSlashedVals)
		actualSlashedVals = types.NewSortedValidatorSet(actualSlashedVals)
		for i := range actualSlashedVals {
			require.Equal(t, excpectedSlashedVals[i], actualSlashedVals[i].GetValAddress())
		}

		// go to the 1st block and thus epoch 1
		ctx = helper.GenAndApplyEmptyBlock(r)
		epochNumber := keeper.GetEpoch(ctx).EpochNumber
		require.Equal(t, uint64(1), epochNumber)
		// no validator is slashed in epoch 1
		require.Empty(t, keeper.GetSlashedValidators(ctx, 1))
	})
}

func sortVals(vals []sdk.ValAddress) {
	sort.Slice(vals, func(i, j int) bool {
		return sdk.BigEndianToUint64(vals[i]) < sdk.BigEndianToUint64(vals[j])
	})
}
