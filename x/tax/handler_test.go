package tax_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/tax"
	"github.com/cosmos/cosmos-sdk/x/tax/keeper"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// createTestInput returns a simapp with custom TaxKeeper
// to avoid messing with the hooks.
func createTestInput() (*simapp.SimApp, sdk.Context, []sdk.AccAddress) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.TaxKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		map[string]bool{},
	)

	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(100000))

	return app, ctx, addrs
}

func TestMsgCreateFixedAmountTax(t *testing.T) {
	app, ctx, addrs := createTestInput()

	taxPoolAddr := addrs[0]
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	epochDays := uint32(1)
	epochAmount := sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1)))

	msg := types.NewMsgCreateFixedAmountTax(
		taxPoolAddr,
		stakingCoinWeights,
		startTime,
		endTime,
		epochDays,
		epochAmount,
	)

	handler := tax.NewHandler(app.TaxKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	taxes := app.TaxKeeper.GetAllTaxes(ctx)
	require.Equal(t, 1, len(taxes))
	require.Equal(t, taxPoolAddr.String(), taxes[0].GetTaxPoolAddress().String())
}

func TestMsgCreateRatioTax(t *testing.T) {
	app, ctx, _ := createTestInput()

	taxPoolAddr := sdk.AccAddress([]byte("taxPoolAddr"))
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)
	epochDays := uint32(1)
	epochAmount := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1)))

	msg := types.NewMsgCreateFixedAmountTax(
		taxPoolAddr,
		stakingCoinWeights,
		startTime,
		endTime,
		epochDays,
		epochAmount,
	)

	handler := tax.NewHandler(app.TaxKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	taxes := app.TaxKeeper.GetAllTaxes(ctx)
	require.Equal(t, 1, len(taxes))
	require.Equal(t, taxPoolAddr.String(), taxes[0].GetTaxPoolAddress().String())
}

func TestMsgStake(t *testing.T) {
	// TODO: not implemented yet
}

func TestMsgUnstake(t *testing.T) {
	// TODO: not implemented yet
}

func TestMsgClaim(t *testing.T) {
	// TODO: not implemented yet
}
