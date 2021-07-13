package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

func TestGetPoolInformation(t *testing.T) {
	commonTerminationAcc := sdk.AccAddress([]byte("terminationAddr"))
	commonStartTime := time.Now().UTC()
	commonEndTime := commonStartTime.AddDate(1, 0, 0)
	commonEpochDays := uint32(1)
	commonCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)

	testCases := []struct {
		taxId          uint64
		taxType        types.TaxType
		taxPoolAddr string
		rewardPoolAddr  string
		terminationAddr string
		reserveAddr     string
		coinWeights     sdk.DecCoins
	}{
		{
			taxId:          uint64(1),
			taxType:        types.TaxTypePublic,
			taxPoolAddr: sdk.AccAddress([]byte("taxPoolAddr1")).String(),
			rewardPoolAddr:  "cosmos1yqurgw7xa94psk95ctje76ferlddg8vykflaln6xsgarj5w6jkrsuvh9dj",
			reserveAddr:     "cosmos18f2zl0q0gpexruasqzav2vfwdthl4779gtmdxgqdpdl03sq9eygq42ff0u",
		},
	}

	for _, tc := range testCases {
		taxName := types.TaxName(tc.taxId, tc.taxType, tc.taxPoolAddr)
		rewardPoolAcc := types.GenerateRewardPoolAcc(taxName)
		stakingReserveAcc := types.GenerateStakingReserveAcc(taxName)
		baseTax := types.NewBaseTax(tc.taxId, tc.taxType, tc.taxPoolAddr, commonTerminationAcc.String(), commonCoinWeights, commonStartTime, commonEndTime, commonEpochDays)
		require.Equal(t, baseTax.RewardPoolAddress, rewardPoolAcc.String())
		require.Equal(t, baseTax.StakingReserveAddress, stakingReserveAcc.String())
	}
}
