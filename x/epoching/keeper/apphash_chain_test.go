package keeper_test

//import (
//	"math/rand"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/cosmos/cosmos-sdk/x/epoching/testepoching/datagen"
//
//	"github.com/cosmos/cosmos-sdk/x/epoching/keeper"
//	"github.com/cosmos/cosmos-sdk/x/epoching/testepoching"
//	"github.com/cosmos/cosmos-sdk/x/epoching/types"
//)
//
//func FuzzAppHashChain(f *testing.F) {
//	datagen.AddRandomSeedsToFuzzer(f, 10)
//
//	f.Fuzz(func(t *testing.T, seed int64) {
//		r := rand.New(rand.NewSource(seed))
//
//		helper := testepoching.NewHelper(t)
//		ctx, k := helper.Ctx, helper.EpochingKeeper
//		// ensure that the epoch info is correct at the genesis
//		epoch := k.GetEpoch(ctx)
//		require.Equal(t, epoch.EpochNumber, uint64(0))
//		require.Equal(t, epoch.FirstBlockHeight, uint64(0))
//
//		// set a random epoch interval
//		epochInterval := r.Uint64()%100 + 2 // the epoch interval should at at least 2
//
//		params := types.Params{
//			EpochInterval: epochInterval,
//		}
//
//		if err := k.SetParams(ctx, params); err != nil {
//			panic(err)
//		}
//
//		// reach the end of the 1st epoch
//		expectedHeight := epochInterval
//		expectedAppHashs := [][]byte{}
//		for i := uint64(0); i < expectedHeight; i++ {
//			ctx = helper.GenAndApplyEmptyBlock(r)
//			expectedAppHashs = append(expectedAppHashs, ctx.BlockHeader().AppHash)
//		}
//		// ensure epoch number is 1
//		epoch = k.GetEpoch(ctx)
//		require.Equal(t, uint64(1), epoch.EpochNumber)
//
//		// ensure appHashs are same as expectedAppHashs
//		appHashs, err := k.GetAllAppHashsForEpoch(ctx, epoch)
//		require.NoError(t, err)
//		require.Equal(t, expectedAppHashs, appHashs)
//
//		// ensure prover and verifier are correct
//		randomHeightInEpoch := uint64(r.Intn(int(expectedHeight)) + 1)
//		randomAppHash, err := k.GetAppHash(ctx, randomHeightInEpoch)
//		require.NoError(t, err)
//		proof, err := k.ProveAppHashInEpoch(ctx, randomHeightInEpoch, epoch.EpochNumber)
//		require.NoError(t, err)
//		err = keeper.VerifyAppHashInclusion(randomAppHash, epoch.AppHashRoot, proof)
//		require.NoError(t, err)
//	})
//}
