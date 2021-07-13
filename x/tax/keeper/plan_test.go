package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

func TestGetSetNewTax(t *testing.T) {
	app, ctx := createTestApp(true)

	taxPoolAddr := sdk.AccAddress([]byte("taxPoolAddr"))
	terminationAddr := sdk.AccAddress([]byte("terminationAddr"))
	farmerAddr := sdk.AccAddress([]byte("farmer"))
	stakingCoins := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000000)))
	coinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	baseTax := types.NewBaseTax(1, 1, taxPoolAddr.String(), terminationAddr.String(), coinWeights, startTime, endTime, 1)
	fixedTax := types.NewFixedAmountTax(baseTax, sdk.NewCoins(sdk.NewCoin("testFarmCoinDenom", sdk.NewInt(1000000))))
	app.TaxKeeper.SetTax(ctx, fixedTax)

	taxGet, found := app.TaxKeeper.GetTax(ctx, 1)
	require.True(t, found)
	require.Equal(t, fixedTax, taxGet)

	taxes := app.TaxKeeper.GetAllTaxes(ctx)
	require.Len(t, taxes, 1)
	require.Equal(t, fixedTax, taxes[0])

	// TODO: tmp test codes for testing functionality, need to separated
	msgStake := types.NewMsgStake(fixedTax.Id, farmerAddr, stakingCoins)
	app.TaxKeeper.Stake(ctx, msgStake)

	stakings := app.TaxKeeper.GetAllStakings(ctx)
	stakingsByTax := app.TaxKeeper.GetStakingsByTaxID(ctx, fixedTax.Id)
	require.Equal(t, stakings, stakingsByTax)
	taxesByFarmer := app.TaxKeeper.GetTaxesByFarmerAddrIndex(ctx, farmerAddr)

	require.Equal(t, taxes, taxesByFarmer)
}
