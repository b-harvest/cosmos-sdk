package farming_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/farming"
	"github.com/cosmos/cosmos-sdk/x/farming/keeper"
	"github.com/cosmos/cosmos-sdk/x/farming/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// createTestInput Returns a simapp with custom StakingKeeper
// to avoid messing with the hooks.
func createTestInput() (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.FarmingKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		map[string]bool{},
	)

	return app, ctx
}

func TestMsgCreateFixedAmountPlan(t *testing.T) {
	app, ctx := createTestInput()

	farmingPoolAddr := sdk.AccAddress([]byte("farmingPoolAddr"))
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime, err := time.Parse(time.RFC3339, "2021-11-01T22:08:41+00:00")
	require.NoError(t, err)
	epochDays := uint32(1)
	epochAmount := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1)))

	msg := types.NewMsgCreateFixedAmountPlan(
		farmingPoolAddr,
		stakingCoinWeights,
		&startTime,
		&endTime,
		epochDays,
		epochAmount,
	)

	handler := farming.NewHandler(app.FarmingKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	plans := app.FarmingKeeper.GetAllPlans(ctx)
	require.Equal(t, 1, len(plans))
	require.Equal(t, farmingPoolAddr.String(), plans[0].GetFarmingPoolAddress().String())
}
