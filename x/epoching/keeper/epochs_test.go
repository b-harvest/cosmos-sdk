package keeper_test

import (
	"math/rand"
	"testing"

	"github.com/babylonchain/babylon/testutil/datagen"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/x/epoching/testepoching"
	"github.com/cosmos/cosmos-sdk/x/epoching/types"
)

func FuzzEpochs(f *testing.F) {
	datagen.AddRandomSeedsToFuzzer(f, 10)

	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))

		helper := testepoching.NewHelper(t)
		ctx, keeper := helper.Ctx, helper.EpochingKeeper
		// ensure that the epoch info is correct at the genesis
		epoch := keeper.GetEpoch(ctx)
		require.Equal(t, epoch.EpochNumber, uint64(0))
		require.Equal(t, epoch.FirstBlockHeight, uint64(0))

		// set a random epoch interval
		epochInterval := r.Uint64()%100 + 2 // the epoch interval should at at least 2

		params := types.Params{
			EpochInterval: epochInterval,
		}

		if err := keeper.SetParams(ctx, params); err != nil {
			panic(err)
		}

		// increment a random number of new blocks
		numIncBlocks := r.Uint64()%1000 + 1
		for i := uint64(0); i < numIncBlocks; i++ {
			ctx = helper.GenAndApplyEmptyBlock(r)
		}

		// ensure that the epoch info is still correct
		expectedEpochNumber := numIncBlocks / epochInterval
		if numIncBlocks%epochInterval > 0 {
			expectedEpochNumber += 1
		}
		actualNewEpoch := keeper.GetEpoch(ctx)
		require.Equal(t, expectedEpochNumber, actualNewEpoch.EpochNumber)
		require.Equal(t, epochInterval, actualNewEpoch.CurrentEpochInterval)
		require.Equal(t, (expectedEpochNumber-1)*epochInterval+1, actualNewEpoch.FirstBlockHeight)
	})
}
