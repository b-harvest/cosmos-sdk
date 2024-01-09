package types_test

import (
	"math/rand"
	"testing"

	"cosmossdk.io/x/epoching/testepoching/datagen"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/x/epoching/types"
)

func FuzzEpoch(f *testing.F) {
	datagen.AddRandomSeedsToFuzzer(f, 10)
	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))

		// generate a random epoch
		epochNumber := uint64(r.Int63()) + 1
		curEpochInterval := r.Uint64()%100 + 2
		firstBlockHeight := r.Uint64() + 1

		e := types.Epoch{
			EpochNumber:          epochNumber,
			CurrentEpochInterval: curEpochInterval,
			FirstBlockHeight:     firstBlockHeight,
		}

		lastBlockHeight := firstBlockHeight + curEpochInterval - 1
		require.Equal(t, lastBlockHeight, e.GetLastBlockHeight())
		secondBlockheight := firstBlockHeight + 1
		require.Equal(t, secondBlockheight, e.GetSecondBlockHeight())
	})
}
