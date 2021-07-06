package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
)

func TestGetSetNewPlan(t *testing.T) {
	app, ctx := createTestApp(true)

	farmingPoolAddr := sdk.AccAddress([]byte("farmingPoolAddr"))
	distPoolAddr := sdk.AccAddress([]byte("distPoolAddr"))
	terminationAddr := sdk.AccAddress([]byte("terminationAddr"))
	reserveAddr := sdk.AccAddress([]byte("reserveAddr"))
	coinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	basePlan := types.NewBasePlan(1, 1, farmingPoolAddr.String(), distPoolAddr.String(), terminationAddr.String(), reserveAddr.String(), coinWeights, startTime, endTime, 1)
	fixedPlan := types.NewFixedAmountPlan(basePlan, sdk.NewCoins(sdk.NewCoin("testFarmCoinDenom", sdk.NewInt(1000000))))
	app.FarmingKeeper.SetPlan(ctx, fixedPlan)

	planGet := app.FarmingKeeper.GetPlan(ctx, 1)
	require.Equal(t, fixedPlan, planGet)

	plans := app.FarmingKeeper.GetAllPlans(ctx)
	require.Len(t, plans, 1)
	require.Equal(t, fixedPlan, plans[0])
}
