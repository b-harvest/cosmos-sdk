package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

const (
	fooDenom     = "foo"
	barDenom     = "bar"
	initialPower = int64(100)
	holder       = "holder"
	multiPerm    = "multiple permissions account"
	randomPerm   = "random permission"
)

var (
	holderAcc     = authtypes.NewEmptyModuleAccount(holder)
	burnerAcc     = authtypes.NewEmptyModuleAccount(authtypes.Burner, authtypes.Burner)
	minterAcc     = authtypes.NewEmptyModuleAccount(authtypes.Minter, authtypes.Minter)
	multiPermAcc  = authtypes.NewEmptyModuleAccount(multiPerm, authtypes.Burner, authtypes.Minter, authtypes.Staking)
	randomPermAcc = authtypes.NewEmptyModuleAccount(randomPerm, "random")

	// The default power validators are initialized to have within tests
	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

func newFooCoin(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(fooDenom, amt)
}

func newBarCoin(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(barDenom, amt)
}

// nolint: interfacer
func getCoinsByName(ctx sdk.Context, bk keeper.Keeper, ak types.AccountKeeper, moduleName string) sdk.Coins {
	moduleAddress := ak.GetModuleAddress(moduleName)
	macc := ak.GetAccount(ctx, moduleAddress)
	if macc == nil {
		return sdk.Coins(nil)
	}

	return bk.GetAllBalances(ctx, macc.GetAddress())
}

type IntegrationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *IntegrationTestSuite) initKeepersWithmAccPerms(blockedAddrs map[string]bool) (authkeeper.AccountKeeper, keeper.BaseKeeper) {
	app := suite.app
	maccPerms := simapp.GetMaccPerms()
	appCodec := simapp.MakeTestEncodingConfig().Marshaler

	maccPerms[holder] = nil
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}
	maccPerms[randomPerm] = []string{"random"}
	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	keeper := keeper.NewBaseKeeper(
		appCodec, app.GetKey(types.StoreKey), authKeeper,
		app.GetSubspace(types.ModuleName), blockedAddrs,
	)

	return authKeeper, keeper
}

func (suite *IntegrationTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.BankKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func (suite *IntegrationTestSuite) TestSupply() {
	ctx := suite.ctx

	require := suite.Require()

	// add module accounts to supply keeper
	authKeeper, keeper := suite.initKeepersWithmAccPerms(make(map[string]bool))

	initialPower := int64(100)
	initTokens := suite.app.StakingKeeper.TokensFromConsensusPower(ctx, initialPower)
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))

	// set burnerAcc balance
	authKeeper.SetModuleAccount(ctx, burnerAcc)
	require.NoError(keeper.MintCoins(ctx, authtypes.Minter, totalSupply))
	require.NoError(keeper.SendCoinsFromModuleToAccount(ctx, authtypes.Minter, burnerAcc.GetAddress(), totalSupply))

	total, _, err := keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	require.NoError(err)
	require.Equal(totalSupply, total)

	// burning all supplied tokens
	err = keeper.BurnCoins(ctx, authtypes.Burner, totalSupply)
	require.NoError(err)

	total, _, err = keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	require.NoError(err)
	require.Equal(total.String(), "")
}

func (suite *IntegrationTestSuite) TestSendCoinsFromModuleToAccount_Blocklist() {
	ctx := suite.ctx

	// add module accounts to supply keeper
	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	_, keeper := suite.initKeepersWithmAccPerms(map[string]bool{addr1.String(): true})

	suite.Require().NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	suite.Require().Error(keeper.SendCoinsFromModuleToAccount(
		ctx, minttypes.ModuleName, addr1, initCoins,
	))
}

func (suite *IntegrationTestSuite) TestSupply_SendCoins() {
	ctx := suite.ctx

	// add module accounts to supply keeper
	authKeeper, keeper := suite.initKeepersWithmAccPerms(make(map[string]bool))

	baseAcc := authKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("baseAcc"))

	// set initial balances
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, minttypes.ModuleName, initCoins))

	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, holderAcc.GetAddress(), initCoins))

	authKeeper.SetModuleAccount(ctx, holderAcc)
	authKeeper.SetModuleAccount(ctx, burnerAcc)
	authKeeper.SetAccount(ctx, baseAcc)

	suite.Require().Panics(func() {
		_ = keeper.SendCoinsFromModuleToModule(ctx, "", holderAcc.GetName(), initCoins) // nolint:errcheck
	})

	suite.Require().Panics(func() {
		_ = keeper.SendCoinsFromModuleToModule(ctx, authtypes.Burner, "", initCoins) // nolint:errcheck
	})

	suite.Require().Panics(func() {
		_ = keeper.SendCoinsFromModuleToAccount(ctx, "", baseAcc.GetAddress(), initCoins) // nolint:errcheck
	})

	suite.Require().Error(
		keeper.SendCoinsFromModuleToAccount(ctx, holderAcc.GetName(), baseAcc.GetAddress(), initCoins.Add(initCoins...)),
	)

	suite.Require().NoError(
		keeper.SendCoinsFromModuleToModule(ctx, holderAcc.GetName(), authtypes.Burner, initCoins),
	)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, holderAcc.GetName()).String())
	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner))

	suite.Require().NoError(
		keeper.SendCoinsFromModuleToAccount(ctx, authtypes.Burner, baseAcc.GetAddress(), initCoins),
	)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner).String())
	suite.Require().Equal(initCoins, keeper.GetAllBalances(ctx, baseAcc.GetAddress()))

	suite.Require().NoError(keeper.SendCoinsFromAccountToModule(ctx, baseAcc.GetAddress(), authtypes.Burner, initCoins))
	suite.Require().Equal(sdk.NewCoins().String(), keeper.GetAllBalances(ctx, baseAcc.GetAddress()).String())
	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner))
}

func (suite *IntegrationTestSuite) TestSupply_MintCoins() {
	ctx := suite.ctx

	// add module accounts to supply keeper
	authKeeper, keeper := suite.initKeepersWithmAccPerms(make(map[string]bool))

	authKeeper.SetModuleAccount(ctx, burnerAcc)
	authKeeper.SetModuleAccount(ctx, minterAcc)
	authKeeper.SetModuleAccount(ctx, multiPermAcc)
	authKeeper.SetModuleAccount(ctx, randomPermAcc)

	initialSupply, _, err := keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	suite.Require().NoError(err)

	suite.Require().Panics(func() { keeper.MintCoins(ctx, "", initCoins) }, "no module account")                // nolint:errcheck
	suite.Require().Panics(func() { keeper.MintCoins(ctx, authtypes.Burner, initCoins) }, "invalid permission") // nolint:errcheck

	err = keeper.MintCoins(ctx, authtypes.Minter, sdk.Coins{sdk.Coin{Denom: "denom", Amount: sdk.NewInt(-10)}})
	suite.Require().Error(err, "insufficient coins")

	suite.Require().Panics(func() { keeper.MintCoins(ctx, randomPerm, initCoins) }) // nolint:errcheck

	err = keeper.MintCoins(ctx, authtypes.Minter, initCoins)
	suite.Require().NoError(err)

	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, authtypes.Minter))
	totalSupply, _, err := keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	suite.Require().NoError(err)

	suite.Require().Equal(initialSupply.Add(initCoins...), totalSupply)

	// test same functionality on module account with multiple permissions
	initialSupply, _, err = keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	suite.Require().NoError(err)

	err = keeper.MintCoins(ctx, multiPermAcc.GetName(), initCoins)
	suite.Require().NoError(err)

	totalSupply, _, err = keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(initCoins, getCoinsByName(ctx, keeper, authKeeper, multiPermAcc.GetName()))
	suite.Require().Equal(initialSupply.Add(initCoins...), totalSupply)
	suite.Require().Panics(func() { keeper.MintCoins(ctx, authtypes.Burner, initCoins) }) // nolint:errcheck
}

func (suite *IntegrationTestSuite) TestSupply_BurnCoins() {
	ctx := suite.ctx
	// add module accounts to supply keeper
	authKeeper, keeper := suite.initKeepersWithmAccPerms(make(map[string]bool))

	// set burnerAcc balance
	authKeeper.SetModuleAccount(ctx, burnerAcc)
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, authtypes.Minter, initCoins))
	suite.
		Require().
		NoError(keeper.SendCoinsFromModuleToAccount(ctx, authtypes.Minter, burnerAcc.GetAddress(), initCoins))

	// inflate supply
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, authtypes.Minter, initCoins))
	supplyAfterInflation, _, err := keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})

	suite.Require().Panics(func() { keeper.BurnCoins(ctx, "", initCoins) }, "no module account")                    // nolint:errcheck
	suite.Require().Panics(func() { keeper.BurnCoins(ctx, authtypes.Minter, initCoins) }, "invalid permission")     // nolint:errcheck
	suite.Require().Panics(func() { keeper.BurnCoins(ctx, randomPerm, supplyAfterInflation) }, "random permission") // nolint:errcheck
	err = keeper.BurnCoins(ctx, authtypes.Burner, supplyAfterInflation)
	suite.Require().Error(err, "insufficient coins")

	err = keeper.BurnCoins(ctx, authtypes.Burner, initCoins)
	suite.Require().NoError(err)
	supplyAfterBurn, _, err := keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, authtypes.Burner).String())
	suite.Require().Equal(supplyAfterInflation.Sub(initCoins), supplyAfterBurn)

	// test same functionality on module account with multiple permissions
	suite.
		Require().
		NoError(keeper.MintCoins(ctx, authtypes.Minter, initCoins))

	supplyAfterInflation, _, err = keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	suite.Require().NoError(err)
	suite.Require().NoError(keeper.SendCoins(ctx, authtypes.NewModuleAddress(authtypes.Minter), multiPermAcc.GetAddress(), initCoins))
	authKeeper.SetModuleAccount(ctx, multiPermAcc)

	err = keeper.BurnCoins(ctx, multiPermAcc.GetName(), initCoins)
	supplyAfterBurn, _, err = keeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{})
	suite.Require().NoError(err)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins().String(), getCoinsByName(ctx, keeper, authKeeper, multiPermAcc.GetName()).String())
	suite.Require().Equal(supplyAfterInflation.Sub(initCoins), supplyAfterBurn)
}

func (suite *IntegrationTestSuite) TestSendCoinsNewAccount() {
	app, ctx := suite.app, suite.ctx
	balances := sdk.NewCoins(newFooCoin(100), newBarCoin(50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, balances))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(balances, acc1Balances)

	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	suite.Require().Nil(app.AccountKeeper.GetAccount(ctx, addr2))
	app.BankKeeper.GetAllBalances(ctx, addr2)
	suite.Require().Empty(app.BankKeeper.GetAllBalances(ctx, addr2))

	sendAmt := sdk.NewCoins(newFooCoin(50), newBarCoin(50))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendAmt))

	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	acc1Balances = app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(sendAmt, acc2Balances)
	updatedAcc1Bal := balances.Sub(sendAmt)
	suite.Require().Len(acc1Balances, len(updatedAcc1Bal))
	suite.Require().Equal(acc1Balances, updatedAcc1Bal)
	suite.Require().NotNil(app.AccountKeeper.GetAccount(ctx, addr2))
}

func (suite *IntegrationTestSuite) TestInputOutputNewAccount() {
	app, ctx := suite.app, suite.ctx

	balances := sdk.NewCoins(newFooCoin(100), newBarCoin(50))
	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, balances))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(balances, acc1Balances)

	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	suite.Require().Nil(app.AccountKeeper.GetAccount(ctx, addr2))
	suite.Require().Empty(app.BankKeeper.GetAllBalances(ctx, addr2))

	inputs := []types.Input{
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}
	outputs := []types.Output{
		{Address: addr2.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}

	suite.Require().NoError(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	expected := sdk.NewCoins(newFooCoin(30), newBarCoin(10))
	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	suite.Require().Equal(expected, acc2Balances)
	suite.Require().NotNil(app.AccountKeeper.GetAccount(ctx, addr2))
}

func (suite *IntegrationTestSuite) TestInputOutputCoins() {
	app, ctx := suite.app, suite.ctx
	balances := sdk.NewCoins(newFooCoin(90), newBarCoin(30))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)

	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)
	app.AccountKeeper.SetAccount(ctx, acc2)

	addr3 := sdk.AccAddress([]byte("addr3_______________"))
	acc3 := app.AccountKeeper.NewAccountWithAddress(ctx, addr3)
	app.AccountKeeper.SetAccount(ctx, acc3)

	inputs := []types.Input{
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}
	outputs := []types.Output{
		{Address: addr2.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
		{Address: addr3.String(), Coins: sdk.NewCoins(newFooCoin(30), newBarCoin(10))},
	}

	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, []types.Output{}))
	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, balances))

	insufficientInputs := []types.Input{
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
		{Address: addr1.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
	}
	insufficientOutputs := []types.Output{
		{Address: addr2.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
		{Address: addr3.String(), Coins: sdk.NewCoins(newFooCoin(300), newBarCoin(100))},
	}
	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, insufficientInputs, insufficientOutputs))
	suite.Require().NoError(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	expected := sdk.NewCoins(newFooCoin(30), newBarCoin(10))
	suite.Require().Equal(expected, acc1Balances)

	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	suite.Require().Equal(expected, acc2Balances)

	acc3Balances := app.BankKeeper.GetAllBalances(ctx, addr3)
	suite.Require().Equal(expected, acc3Balances)
}

func (suite *IntegrationTestSuite) TestSendCoins() {
	app, ctx := suite.app, suite.ctx
	balances := sdk.NewCoins(newFooCoin(100), newBarCoin(50))

	addr1 := sdk.AccAddress("addr1_______________")
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc1)

	addr2 := sdk.AccAddress("addr2_______________")
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)
	app.AccountKeeper.SetAccount(ctx, acc2)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, balances))

	sendAmt := sdk.NewCoins(newFooCoin(50), newBarCoin(25))
	suite.Require().Error(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendAmt))

	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, balances))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendAmt))

	acc1Balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	expected := sdk.NewCoins(newFooCoin(50), newBarCoin(25))
	suite.Require().Equal(expected, acc1Balances)

	acc2Balances := app.BankKeeper.GetAllBalances(ctx, addr2)
	expected = sdk.NewCoins(newFooCoin(150), newBarCoin(75))
	suite.Require().Equal(expected, acc2Balances)

	// we sent all foo coins to acc2, so foo balance should be deleted for acc1 and bar should be still there
	var coins []sdk.Coin
	app.BankKeeper.IterateAccountBalances(ctx, addr1, func(c sdk.Coin) (stop bool) {
		coins = append(coins, c)
		return true
	})
	suite.Require().Len(coins, 1)
	suite.Require().Equal(newBarCoin(25), coins[0], "expected only bar coins in the account balance, got: %v", coins)
}

func (suite *IntegrationTestSuite) TestValidateBalance() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	suite.Require().Error(app.BankKeeper.ValidateBalance(ctx, addr1))

	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	app.AccountKeeper.SetAccount(ctx, acc)

	balances := sdk.NewCoins(newFooCoin(100))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, balances))
	suite.Require().NoError(app.BankKeeper.ValidateBalance(ctx, addr1))

	bacc := authtypes.NewBaseAccountWithAddress(addr2)
	vacc := vesting.NewContinuousVestingAccount(bacc, balances.Add(balances...), now.Unix(), endTime.Unix())

	app.AccountKeeper.SetAccount(ctx, vacc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, balances))
	suite.Require().Error(app.BankKeeper.ValidateBalance(ctx, addr2))
}

func (suite *IntegrationTestSuite) TestSendEnabled() {
	app, ctx := suite.app, suite.ctx
	enabled := true
	params := types.DefaultParams()
	suite.Require().Equal(enabled, params.DefaultSendEnabled)

	app.BankKeeper.SetParams(ctx, params)

	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())
	fooCoin := sdk.NewCoin("foocoin", sdk.OneInt())
	barCoin := sdk.NewCoin("barcoin", sdk.OneInt())

	// assert with default (all denom) send enabled both Bar and Bond Denom are enabled
	suite.Require().Equal(enabled, app.BankKeeper.IsSendEnabledCoin(ctx, barCoin))
	suite.Require().Equal(enabled, app.BankKeeper.IsSendEnabledCoin(ctx, bondCoin))

	// Both coins should be send enabled.
	err := app.BankKeeper.IsSendEnabledCoins(ctx, fooCoin, bondCoin)
	suite.Require().NoError(err)

	// Set default send_enabled to !enabled, add a foodenom that overrides default as enabled
	params.DefaultSendEnabled = !enabled
	params = params.SetSendEnabledParam(fooCoin.Denom, enabled)
	app.BankKeeper.SetParams(ctx, params)

	// Expect our specific override to be enabled, others to be !enabled.
	suite.Require().Equal(enabled, app.BankKeeper.IsSendEnabledCoin(ctx, fooCoin))
	suite.Require().Equal(!enabled, app.BankKeeper.IsSendEnabledCoin(ctx, barCoin))
	suite.Require().Equal(!enabled, app.BankKeeper.IsSendEnabledCoin(ctx, bondCoin))

	// Foo coin should be send enabled.
	err = app.BankKeeper.IsSendEnabledCoins(ctx, fooCoin)
	suite.Require().NoError(err)

	// Expect an error when one coin is not send enabled.
	err = app.BankKeeper.IsSendEnabledCoins(ctx, fooCoin, bondCoin)
	suite.Require().Error(err)

	// Expect an error when all coins are not send enabled.
	err = app.BankKeeper.IsSendEnabledCoins(ctx, bondCoin, barCoin)
	suite.Require().Error(err)
}

func (suite *IntegrationTestSuite) TestHasBalance() {
	app, ctx := suite.app, suite.ctx
	addr := sdk.AccAddress([]byte("addr1_______________"))

	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)

	balances := sdk.NewCoins(newFooCoin(100))
	suite.Require().False(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(99)))

	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr, balances))
	suite.Require().False(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(101)))
	suite.Require().True(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(100)))
	suite.Require().True(app.BankKeeper.HasBalance(ctx, addr, newFooCoin(1)))
}

func (suite *IntegrationTestSuite) TestMsgSendEvents() {
	app, ctx := suite.app, suite.ctx
	addr := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, acc)
	newCoins := sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr, newCoins))

	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr, addr2, newCoins))
	event1 := sdk.Event{
		Type:       types.EventTypeTransfer,
		Attributes: []abci.EventAttribute{},
	}
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeyRecipient), Value: []byte(addr2.String())},
	)
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr.String())},
	)
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(sdk.AttributeKeyAmount), Value: []byte(newCoins.String())},
	)

	event2 := sdk.Event{
		Type:       sdk.EventTypeMessage,
		Attributes: []abci.EventAttribute{},
	}
	event2.Attributes = append(
		event2.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr.String())},
	)

	// events are shifted due to the funding account events
	events := ctx.EventManager().ABCIEvents()
	suite.Require().Equal(10, len(events))
	suite.Require().Equal(abci.Event(event1), events[8])
	suite.Require().Equal(abci.Event(event2), events[9])
}

func (suite *IntegrationTestSuite) TestMsgMultiSendEvents() {
	app, ctx := suite.app, suite.ctx

	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	addr := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addr3 := sdk.AccAddress([]byte("addr3_______________"))
	addr4 := sdk.AccAddress([]byte("addr4_______________"))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, acc)
	app.AccountKeeper.SetAccount(ctx, acc2)

	newCoins := sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))
	newCoins2 := sdk.NewCoins(sdk.NewInt64Coin(barDenom, 100))
	inputs := []types.Input{
		{Address: addr.String(), Coins: newCoins},
		{Address: addr2.String(), Coins: newCoins2},
	}
	outputs := []types.Output{
		{Address: addr3.String(), Coins: newCoins},
		{Address: addr4.String(), Coins: newCoins2},
	}

	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	events := ctx.EventManager().ABCIEvents()
	suite.Require().Equal(0, len(events))

	// Set addr's coins but not addr2's coins
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr, sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))))
	suite.Require().Error(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	events = ctx.EventManager().ABCIEvents()
	suite.Require().Equal(8, len(events)) // 7 events because account funding causes extra minting + coin_spent + coin_recv events

	event1 := sdk.Event{
		Type:       sdk.EventTypeMessage,
		Attributes: []abci.EventAttribute{},
	}
	event1.Attributes = append(
		event1.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr.String())},
	)
	suite.Require().Equal(abci.Event(event1), events[7])

	// Set addr's coins and addr2's coins
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr, sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))))
	newCoins = sdk.NewCoins(sdk.NewInt64Coin(fooDenom, 50))

	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, sdk.NewCoins(sdk.NewInt64Coin(barDenom, 100))))
	newCoins2 = sdk.NewCoins(sdk.NewInt64Coin(barDenom, 100))

	suite.Require().NoError(app.BankKeeper.InputOutputCoins(ctx, inputs, outputs))

	events = ctx.EventManager().ABCIEvents()
	suite.Require().Equal(28, len(events)) // 25 due to account funding + coin_spent + coin_recv events

	event2 := sdk.Event{
		Type:       sdk.EventTypeMessage,
		Attributes: []abci.EventAttribute{},
	}
	event2.Attributes = append(
		event2.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeySender), Value: []byte(addr2.String())},
	)
	event3 := sdk.Event{
		Type:       types.EventTypeTransfer,
		Attributes: []abci.EventAttribute{},
	}
	event3.Attributes = append(
		event3.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeyRecipient), Value: []byte(addr3.String())},
	)
	event3.Attributes = append(
		event3.Attributes,
		abci.EventAttribute{Key: []byte(sdk.AttributeKeyAmount), Value: []byte(newCoins.String())})
	event4 := sdk.Event{
		Type:       types.EventTypeTransfer,
		Attributes: []abci.EventAttribute{},
	}
	event4.Attributes = append(
		event4.Attributes,
		abci.EventAttribute{Key: []byte(types.AttributeKeyRecipient), Value: []byte(addr4.String())},
	)
	event4.Attributes = append(
		event4.Attributes,
		abci.EventAttribute{Key: []byte(sdk.AttributeKeyAmount), Value: []byte(newCoins2.String())},
	)
	// events are shifted due to the funding account events
	suite.Require().Equal(abci.Event(event1), events[21])
	suite.Require().Equal(abci.Event(event2), events[23])
	suite.Require().Equal(abci.Event(event3), events[25])
	suite.Require().Equal(abci.Event(event4), events[27])
}

func (suite *IntegrationTestSuite) TestSpendableCoins() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	delCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))

	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule)
	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, macc)
	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, origCoins))

	suite.Require().Equal(origCoins, app.BankKeeper.SpendableCoins(ctx, addr2))

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins))
	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.SpendableCoins(ctx, addr1))
}

func (suite *IntegrationTestSuite) TestVestingAccountSend() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, now.Unix(), endTime.Unix())

	app.AccountKeeper.SetAccount(ctx, vacc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))

	// require that no coins be sendable at the beginning of the vesting schedule
	suite.Require().Error(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))

	// receive some coins
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, sendCoins))
	// require that all vested coins are spendable plus any received
	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))
	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr1))
}

func (suite *IntegrationTestSuite) TestPeriodicVestingAccountSend() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	periods := vesting.Periods{
		vesting.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 50)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
	}

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewPeriodicVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), periods)

	app.AccountKeeper.SetAccount(ctx, vacc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))

	// require that no coins be sendable at the beginning of the vesting schedule
	suite.Require().Error(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))

	// receive some coins
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, sendCoins))

	// require that all vested coins are spendable plus any received
	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoins))
	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr1))
}

func (suite *IntegrationTestSuite) TestVestingAccountReceive() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, origCoins))

	// send some coins to the vesting account
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr2, addr1, sendCoins))

	// require the coins are spendable
	vacc = app.AccountKeeper.GetAccount(ctx, addr1).(*vesting.ContinuousVestingAccount)
	balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(origCoins.Add(sendCoins...), balances)
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now)), sendCoins)

	// require coins are spendable plus any that have vested
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now.Add(12*time.Hour))), origCoins)
}

func (suite *IntegrationTestSuite) TestPeriodicVestingAccountReceive() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	periods := vesting.Periods{
		vesting.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 50)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
		vesting.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin("stake", 25)}},
	}

	vacc := vesting.NewPeriodicVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), periods)
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, origCoins))

	// send some coins to the vesting account
	suite.Require().NoError(app.BankKeeper.SendCoins(ctx, addr2, addr1, sendCoins))

	// require the coins are spendable
	vacc = app.AccountKeeper.GetAccount(ctx, addr1).(*vesting.PeriodicVestingAccount)
	balances := app.BankKeeper.GetAllBalances(ctx, addr1)
	suite.Require().Equal(origCoins.Add(sendCoins...), balances)
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now)), sendCoins)

	// require coins are spendable plus any that have vested
	suite.Require().Equal(balances.Sub(vacc.LockedCoins(now.Add(12*time.Hour))), origCoins)
}

func (suite *IntegrationTestSuite) TestDelegateCoins() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	delCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))

	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)
	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, origCoins))

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))

	// require the ability for a non-vesting account to delegate
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins))
	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.GetAllBalances(ctx, addr2))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addrModule))

	// require the ability for a vesting account to delegate
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addr1))

	// require that delegated vesting amount is equal to what was delegated with DelegateCoins
	acc = app.AccountKeeper.GetAccount(ctx, addr1)
	vestingAcc, ok := acc.(exported.VestingAccount)
	suite.Require().True(ok)
	suite.Require().Equal(delCoins, vestingAcc.GetDelegatedVesting())
}

func (suite *IntegrationTestSuite) TestDelegateCoins_Invalid() {
	app, ctx := suite.app, suite.ctx

	origCoins := sdk.NewCoins(newFooCoin(100))
	delCoins := sdk.NewCoins(newFooCoin(50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))
	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)

	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))
	invalidCoins := sdk.Coins{sdk.Coin{Denom: "fooDenom", Amount: sdk.NewInt(-50)}}
	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, invalidCoins))

	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))
	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().Error(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, origCoins.Add(origCoins...)))
}

func (suite *IntegrationTestSuite) TestUndelegateCoins() {
	app, ctx := suite.app, suite.ctx
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	delCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))

	bacc := authtypes.NewBaseAccountWithAddress(addr1)
	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing

	vacc := vesting.NewContinuousVestingAccount(bacc, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, acc)
	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, origCoins))

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))

	// require the ability for a non-vesting account to delegate
	err := app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins)
	suite.Require().NoError(err)

	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.GetAllBalances(ctx, addr2))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addrModule))

	// require the ability for a non-vesting account to undelegate
	suite.Require().NoError(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr2, delCoins))

	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr2))
	suite.Require().True(app.BankKeeper.GetAllBalances(ctx, addrModule).Empty())

	// require the ability for a vesting account to delegate
	suite.Require().NoError(app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins))

	suite.Require().Equal(origCoins.Sub(delCoins), app.BankKeeper.GetAllBalances(ctx, addr1))
	suite.Require().Equal(delCoins, app.BankKeeper.GetAllBalances(ctx, addrModule))

	// require the ability for a vesting account to undelegate
	suite.Require().NoError(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))

	suite.Require().Equal(origCoins, app.BankKeeper.GetAllBalances(ctx, addr1))
	suite.Require().True(app.BankKeeper.GetAllBalances(ctx, addrModule).Empty())

	// require that delegated vesting amount is completely empty, since they were completely undelegated
	acc = app.AccountKeeper.GetAccount(ctx, addr1)
	vestingAcc, ok := acc.(exported.VestingAccount)
	suite.Require().True(ok)
	suite.Require().Empty(vestingAcc.GetDelegatedVesting())
}

func (suite *IntegrationTestSuite) TestUndelegateCoins_Invalid() {
	app, ctx := suite.app, suite.ctx

	origCoins := sdk.NewCoins(newFooCoin(100))
	delCoins := sdk.NewCoins(newFooCoin(50))

	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addrModule := sdk.AccAddress([]byte("moduleAcc___________"))
	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)

	suite.Require().Error(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))

	app.AccountKeeper.SetAccount(ctx, macc)
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))

	suite.Require().Error(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))
	app.AccountKeeper.SetAccount(ctx, acc)

	suite.Require().Error(app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins))
}

func (suite *IntegrationTestSuite) TestSetDenomMetaData() {
	app, ctx := suite.app, suite.ctx

	metadata := suite.getTestMetadata()

	for i := range []int{1, 2} {
		app.BankKeeper.SetDenomMetaData(ctx, metadata[i])
	}

	actualMetadata, found := app.BankKeeper.GetDenomMetaData(ctx, metadata[1].Base)
	suite.Require().True(found)
	suite.Require().Equal(metadata[1].GetBase(), actualMetadata.GetBase())
	suite.Require().Equal(metadata[1].GetDisplay(), actualMetadata.GetDisplay())
	suite.Require().Equal(metadata[1].GetDescription(), actualMetadata.GetDescription())
	suite.Require().Equal(metadata[1].GetDenomUnits()[1].GetDenom(), actualMetadata.GetDenomUnits()[1].GetDenom())
	suite.Require().Equal(metadata[1].GetDenomUnits()[1].GetExponent(), actualMetadata.GetDenomUnits()[1].GetExponent())
	suite.Require().Equal(metadata[1].GetDenomUnits()[1].GetAliases(), actualMetadata.GetDenomUnits()[1].GetAliases())
}

func (suite *IntegrationTestSuite) TestIterateAllDenomMetaData() {
	app, ctx := suite.app, suite.ctx

	expectedMetadata := suite.getTestMetadata()
	// set metadata
	for i := range []int{1, 2} {
		app.BankKeeper.SetDenomMetaData(ctx, expectedMetadata[i])
	}
	// retrieve metadata
	actualMetadata := make([]types.Metadata, 0)
	app.BankKeeper.IterateAllDenomMetaData(ctx, func(metadata types.Metadata) bool {
		actualMetadata = append(actualMetadata, metadata)
		return false
	})
	// execute checks
	for i := range []int{1, 2} {
		suite.Require().Equal(expectedMetadata[i].GetBase(), actualMetadata[i].GetBase())
		suite.Require().Equal(expectedMetadata[i].GetDisplay(), actualMetadata[i].GetDisplay())
		suite.Require().Equal(expectedMetadata[i].GetDescription(), actualMetadata[i].GetDescription())
		suite.Require().Equal(expectedMetadata[i].GetDenomUnits()[1].GetDenom(), actualMetadata[i].GetDenomUnits()[1].GetDenom())
		suite.Require().Equal(expectedMetadata[i].GetDenomUnits()[1].GetExponent(), actualMetadata[i].GetDenomUnits()[1].GetExponent())
		suite.Require().Equal(expectedMetadata[i].GetDenomUnits()[1].GetAliases(), actualMetadata[i].GetDenomUnits()[1].GetAliases())
	}
}

func (suite *IntegrationTestSuite) TestBalanceTrackingEvents() {
	// replace account keeper and bank keeper otherwise the account keeper won't be aware of the
	// existence of the new module account because GetModuleAccount checks for the existence via
	// permissions map and not via state... weird
	maccPerms := simapp.GetMaccPerms()
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}

	suite.app.AccountKeeper = authkeeper.NewAccountKeeper(
		suite.app.AppCodec(), suite.app.GetKey(authtypes.StoreKey), suite.app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)

	suite.app.BankKeeper = keeper.NewBaseKeeper(suite.app.AppCodec(), suite.app.GetKey(types.StoreKey),
		suite.app.AccountKeeper, suite.app.GetSubspace(types.ModuleName), nil)

	// set account with multiple permissions
	suite.app.AccountKeeper.SetModuleAccount(suite.ctx, multiPermAcc)
	// mint coins
	suite.Require().NoError(
		suite.app.BankKeeper.MintCoins(
			suite.ctx,
			multiPermAcc.Name,
			sdk.NewCoins(sdk.NewCoin("utxo", sdk.NewInt(100000)))),
	)
	// send coins to address
	addr1 := sdk.AccAddress("addr1_______________")
	suite.Require().NoError(
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(
			suite.ctx,
			multiPermAcc.Name,
			addr1,
			sdk.NewCoins(sdk.NewCoin("utxo", sdk.NewInt(50000))),
		),
	)

	// burn coins from module account
	suite.Require().NoError(
		suite.app.BankKeeper.BurnCoins(
			suite.ctx,
			multiPermAcc.Name,
			sdk.NewCoins(sdk.NewInt64Coin("utxo", 1000)),
		),
	)

	// process balances and supply from events
	supply := sdk.NewCoins()

	balances := make(map[string]sdk.Coins)

	for _, e := range suite.ctx.EventManager().ABCIEvents() {
		switch e.Type {
		case types.EventTypeCoinBurn:
			burnedCoins, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			supply = supply.Sub(burnedCoins)

		case types.EventTypeCoinMint:
			mintedCoins, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			supply = supply.Add(mintedCoins...)

		case types.EventTypeCoinSpent:
			coinsSpent, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			spender, err := sdk.AccAddressFromBech32((string)(e.Attributes[0].Value))
			suite.Require().NoError(err)
			balances[spender.String()] = balances[spender.String()].Sub(coinsSpent)

		case types.EventTypeCoinReceived:
			coinsRecv, err := sdk.ParseCoinsNormalized((string)(e.Attributes[1].Value))
			suite.Require().NoError(err)
			receiver, err := sdk.AccAddressFromBech32((string)(e.Attributes[0].Value))
			suite.Require().NoError(err)
			balances[receiver.String()] = balances[receiver.String()].Add(coinsRecv...)
		}
	}

	// check balance and supply tracking
	suite.Require().True(suite.app.BankKeeper.HasSupply(suite.ctx, "utxo"))
	savedSupply := suite.app.BankKeeper.GetSupply(suite.ctx, "utxo")
	utxoSupply := savedSupply
	suite.Require().Equal(utxoSupply.Amount, supply.AmountOf("utxo"))
	// iterate accounts and check balances
	suite.app.BankKeeper.IterateAllBalances(suite.ctx, func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {
		// if it's not utxo coin then skip
		if coin.Denom != "utxo" {
			return false
		}

		balance, exists := balances[address.String()]
		suite.Require().True(exists)

		expectedUtxo := sdk.NewCoin("utxo", balance.AmountOf(coin.Denom))
		suite.Require().Equal(expectedUtxo.String(), coin.String())
		return false
	})
}

func (suite *IntegrationTestSuite) getTestMetadata() []types.Metadata {
	return []types.Metadata{
		{
			Name:        "Cosmos Hub Atom",
			Symbol:      "ATOM",
			Description: "The native staking token of the Cosmos Hub.",
			DenomUnits: []*types.DenomUnit{
				{"uatom", uint32(0), []string{"microatom"}},
				{"matom", uint32(3), []string{"milliatom"}},
				{"atom", uint32(6), nil},
			},
			Base:    "uatom",
			Display: "atom",
		},
		{
			Name:        "Token",
			Symbol:      "TOKEN",
			Description: "The native staking token of the Token Hub.",
			DenomUnits: []*types.DenomUnit{
				{"1token", uint32(5), []string{"decitoken"}},
				{"2token", uint32(4), []string{"centitoken"}},
				{"3token", uint32(7), []string{"dekatoken"}},
			},
			Base:    "utoken",
			Display: "token",
		},
	}
}

func (suite *IntegrationTestSuite) TestMintCoinRestrictions() {
	type BankMintingRestrictionFn func(ctx sdk.Context, coins sdk.Coins) error

	maccPerms := simapp.GetMaccPerms()
	maccPerms[multiPerm] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}

	suite.app.AccountKeeper = authkeeper.NewAccountKeeper(
		suite.app.AppCodec(), suite.app.GetKey(authtypes.StoreKey), suite.app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	suite.app.AccountKeeper.SetModuleAccount(suite.ctx, multiPermAcc)

	type testCase struct {
		coinsToTry sdk.Coin
		expectPass bool
	}

	tests := []struct {
		name          string
		restrictionFn BankMintingRestrictionFn
		testCases     []testCase
	}{
		{
			"restriction",
			func(ctx sdk.Context, coins sdk.Coins) error {
				for _, coin := range coins {
					if coin.Denom != fooDenom {
						return fmt.Errorf("Module %s only has perms for minting %s coins, tried minting %s coins", types.ModuleName, fooDenom, coin.Denom)
					}
				}
				return nil
			},
			[]testCase{
				{
					coinsToTry: newFooCoin(100),
					expectPass: true,
				},
				{
					coinsToTry: newBarCoin(100),
					expectPass: false,
				},
			},
		},
	}

	for _, test := range tests {
		suite.app.BankKeeper = keeper.NewBaseKeeper(suite.app.AppCodec(), suite.app.GetKey(types.StoreKey),
			suite.app.AccountKeeper, suite.app.GetSubspace(types.ModuleName), nil).WithMintCoinsRestriction(keeper.MintingRestrictionFn(test.restrictionFn))
		for _, testCase := range test.testCases {
			if testCase.expectPass {
				suite.Require().NoError(
					suite.app.BankKeeper.MintCoins(
						suite.ctx,
						multiPermAcc.Name,
						sdk.NewCoins(testCase.coinsToTry),
					),
				)
			} else {
				suite.Require().Error(
					suite.app.BankKeeper.MintCoins(
						suite.ctx,
						multiPermAcc.Name,
						sdk.NewCoins(testCase.coinsToTry),
					),
				)
			}
		}
	}
}

func (suite *IntegrationTestSuite) TestSendRestrictions() {
	type BankSendRestrictionFn func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error)

	// Define a new module account for the test
	newModuleAcc := authtypes.NewEmptyModuleAccount("newbankmodule", authtypes.Minter, authtypes.Burner)

	// Setup the permissions for the module account
	maccPerms := simapp.GetMaccPerms()
	maccPerms[newModuleAcc.Name] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}

	// Update the AccountKeeper with the new module account
	suite.app.AccountKeeper = authkeeper.NewAccountKeeper(
		suite.app.AppCodec(), suite.app.GetKey(authtypes.StoreKey), suite.app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	suite.app.AccountKeeper.SetModuleAccount(suite.ctx, newModuleAcc)

	// Initialize the BankKeeper
	suite.app.BankKeeper = keeper.NewBaseKeeper(
		suite.app.AppCodec(), suite.app.GetKey(types.StoreKey),
		suite.app.AccountKeeper, suite.app.GetSubspace(types.ModuleName), nil,
	)

	// Define normal addresses
	normalAddr1 := sdk.AccAddress("addr1---------------")
	normalAddr2 := sdk.AccAddress("addr2---------------")
	blockedAddr := sdk.AccAddress("blockedaddr---------")
	coin := sdk.NewCoin("testcoin", sdk.NewInt(1000))

	acc1 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, normalAddr1)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc1)

	acc2 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, normalAddr2)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc2)

	suite.Require().NoError(simapp.FundAccount(suite.app.BankKeeper, suite.ctx, normalAddr1, sdk.NewCoins(coin)))
	suite.Require().NoError(simapp.FundAccount(suite.app.BankKeeper, suite.ctx, normalAddr2, sdk.NewCoins(coin)))

	// Test cases
	type testCase struct {
		fromAddr      sdk.AccAddress
		toAddr        sdk.AccAddress
		coinsToTry    sdk.Coins
		expectedError error
	}

	tests := []struct {
		name          string
		setup         func(ctx sdk.Context)
		restrictionFn BankSendRestrictionFn
		testCases     []testCase
	}{
		{
			name: "normal transfer",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				return toAddr, nil
			},
			testCases: []testCase{
				{
					fromAddr:      normalAddr1,
					toAddr:        normalAddr2,
					coinsToTry:    sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(100))),
					expectedError: nil,
				},
			},
		},
		{
			name: "blocked sender",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if fromAddr.Equals(blockedAddr) {
					return toAddr, fmt.Errorf("%s is blocked from sending %s", fromAddr, coin.Denom)
				}
				return toAddr, nil
			},
			testCases: []testCase{
				{
					fromAddr:      blockedAddr,
					toAddr:        normalAddr2,
					coinsToTry:    sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(100))),
					expectedError: fmt.Errorf("%s is blocked from sending %s", blockedAddr, coin.Denom),
				},
			},
		},
		{
			name: "blocked receiver",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if toAddr.Equals(blockedAddr) {
					return toAddr, fmt.Errorf("%s is blocked from receiving %s", toAddr, coin.Denom)
				}
				return toAddr, nil
			},
			testCases: []testCase{
				{
					fromAddr:      normalAddr1,
					toAddr:        blockedAddr,
					coinsToTry:    sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(100))),
					expectedError: fmt.Errorf("%s is blocked from receiving %s", blockedAddr, coin.Denom),
				},
			},
		},
		{
			name: "send coins to module account",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if !toAddr.Equals(newModuleAcc.GetAddress()) {
					return toAddr, fmt.Errorf("only module account can receive coins for burning")
				}
				return toAddr, nil
			},
			testCases: []testCase{
				{
					fromAddr:      normalAddr1,
					toAddr:        newModuleAcc.GetAddress(),
					coinsToTry:    sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(100))),
					expectedError: nil,
				},
			},
		},
		{
			name: "receive coins from module account",
			setup: func(ctx sdk.Context) {
				suite.Require().NoError(
					suite.app.BankKeeper.MintCoins(
						ctx,
						"newbankmodule",
						sdk.NewCoins(coin)),
				)
			},
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if !fromAddr.Equals(newModuleAcc.GetAddress()) {
					return toAddr, fmt.Errorf("only module account can mint coins")
				}
				return toAddr, nil
			},
			testCases: []testCase{
				{
					fromAddr:      newModuleAcc.GetAddress(),
					toAddr:        normalAddr1,
					coinsToTry:    sdk.NewCoins(coin),
					expectedError: nil,
				},
			},
		},
	}

	for _, test := range tests {
		// Set the restriction function for this particular test
		suite.app.BankKeeper = keeper.NewBaseKeeper(
			suite.app.AppCodec(), suite.app.GetKey(types.StoreKey),
			suite.app.AccountKeeper, suite.app.GetSubspace(types.ModuleName), nil,
		).WithSendCoinsRestriction(keeper.SendRestrictionFn(test.restrictionFn))

		// Execute each test case within the current test
		for _, testCase := range test.testCases {
			if test.setup != nil {
				test.setup(suite.ctx)
			}
			if testCase.expectedError == nil {
				err := suite.app.BankKeeper.SendCoins(suite.ctx, testCase.fromAddr, testCase.toAddr, testCase.coinsToTry)
				suite.Require().NoError(err, "Test case failed: %v", testCase)
			} else {
				err := suite.app.BankKeeper.SendCoins(suite.ctx, testCase.fromAddr, testCase.toAddr, testCase.coinsToTry)
				suite.Require().EqualError(err, testCase.expectedError.Error(), "expected error: %v but got: %v", testCase.expectedError, err)
			}
		}
	}
}

func (suite *IntegrationTestSuite) TestNestedSendRestrictions() {
	type BankSendRestrictionFn func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error)

	// Define a new module account for the test
	newModuleAcc := authtypes.NewEmptyModuleAccount("newbankmodule", authtypes.Minter, authtypes.Burner)

	// Setup the permissions for the module account
	maccPerms := simapp.GetMaccPerms()
	maccPerms[newModuleAcc.Name] = []string{authtypes.Burner, authtypes.Minter, authtypes.Staking}

	// Update the AccountKeeper with the new module account
	suite.app.AccountKeeper = authkeeper.NewAccountKeeper(
		suite.app.AppCodec(), suite.app.GetKey(authtypes.StoreKey), suite.app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	suite.app.AccountKeeper.SetModuleAccount(suite.ctx, newModuleAcc)

	// Initialize the BankKeeper
	suite.app.BankKeeper = keeper.NewBaseKeeper(
		suite.app.AppCodec(), suite.app.GetKey(types.StoreKey),
		suite.app.AccountKeeper, suite.app.GetSubspace(types.ModuleName), nil,
	)

	// Define normal addresses
	normalAddr1 := sdk.AccAddress("addr1---------------")
	normalAddr2 := sdk.AccAddress("addr2---------------")
	blockedAddr := sdk.AccAddress("blockedaddr---------")
	coin := sdk.NewCoin("testcoin", sdk.NewInt(1000))

	acc1 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, normalAddr1)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc1)

	acc2 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, normalAddr2)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc2)

	suite.Require().NoError(simapp.FundAccount(suite.app.BankKeeper, suite.ctx, normalAddr1, sdk.NewCoins(coin)))
	suite.Require().NoError(simapp.FundAccount(suite.app.BankKeeper, suite.ctx, normalAddr2, sdk.NewCoins(coin)))

	// Test cases
	type testCase struct {
		fromAddr      sdk.AccAddress
		toAddr        sdk.AccAddress
		coinsToTry    sdk.Coins
		expectedError error
	}

	tests := []struct {
		name           string
		restrictionFns []BankSendRestrictionFn // List of multiple restrictions to apply
		testCases      []testCase
	}{
		{
			"nested restrictions - sender blocked then amount restricted",
			[]BankSendRestrictionFn{
				// First restriction: Block the sender if it's blockedAddr
				func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
					if fromAddr.Equals(blockedAddr) {
						return toAddr, fmt.Errorf("%s is blocked from sending %s", fromAddr, coin.Denom)
					}
					return toAddr, nil
				},
				// Second restriction: Restrict the amount from being sent if it's greater than 500
				func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
					if amt.AmountOf(coin.Denom).GT(sdk.NewInt(500)) {
						return toAddr, fmt.Errorf("cannot send more than 500 %s", coin.Denom)
					}
					return toAddr, nil
				},
			},
			[]testCase{
				{
					fromAddr:      blockedAddr,
					toAddr:        normalAddr1,
					coinsToTry:    sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(100))),
					expectedError: fmt.Errorf("%s is blocked from sending %s", blockedAddr, coin.Denom),
				},
				{
					fromAddr:      normalAddr1,
					toAddr:        normalAddr2,
					coinsToTry:    sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(600))),
					expectedError: fmt.Errorf("cannot send more than 500 %s", coin.Denom),
				},
				{
					fromAddr:      normalAddr1,
					toAddr:        normalAddr2,
					coinsToTry:    sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(400))),
					expectedError: nil,
				},
			},
		},
	}

	for _, test := range tests {
		// Set up the nested restrictions by applying each restriction function sequentially
		suite.app.BankKeeper = keeper.NewBaseKeeper(
			suite.app.AppCodec(), suite.app.GetKey(types.StoreKey),
			suite.app.AccountKeeper, suite.app.GetSubspace(types.ModuleName), nil,
		).WithSendCoinsRestriction(keeper.SendRestrictionFn(test.restrictionFns[0])).WithSendCoinsRestriction(keeper.SendRestrictionFn(test.restrictionFns[1]))

		// Execute each test case within the current test
		for _, testCase := range test.testCases {
			if testCase.expectedError == nil {
				err := suite.app.BankKeeper.SendCoins(suite.ctx, testCase.fromAddr, testCase.toAddr, testCase.coinsToTry)
				suite.Require().NoError(err, "Test case failed: %v", testCase)
			} else {
				err := suite.app.BankKeeper.SendCoins(suite.ctx, testCase.fromAddr, testCase.toAddr, testCase.coinsToTry)
				suite.Require().EqualError(err, testCase.expectedError.Error(), "expected error: %v but got: %v", testCase.expectedError, err)
			}
		}
	}
}

func (suite *IntegrationTestSuite) TestDelegateUndelegateRestrictions() {
	// Initial setup
	app, ctx := suite.app, suite.ctx

	// Set block time
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	endTime := now.Add(24 * time.Hour)

	// Set coin amounts
	origCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	delCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))

	// Set account addresses
	addr1 := sdk.AccAddress("addr1_______________") // vesting account
	addr2 := sdk.AccAddress("addr2_______________") // non-vesting account
	addrModule := sdk.AccAddress("moduleAcc___________")

	// Create base and module accounts
	bacc1 := authtypes.NewBaseAccountWithAddress(addr1)
	bacc2 := authtypes.NewBaseAccountWithAddress(addr2)
	macc := app.AccountKeeper.NewAccountWithAddress(ctx, addrModule)

	// Create continuous vesting account for addr1 (vesting account)
	vacc := vesting.NewContinuousVestingAccount(bacc1, origCoins, ctx.BlockHeader().Time.Unix(), endTime.Unix())

	// Set accounts in the AccountKeeper
	app.AccountKeeper.SetAccount(ctx, vacc)
	app.AccountKeeper.SetAccount(ctx, bacc2)
	app.AccountKeeper.SetAccount(ctx, macc)

	// Fund both accounts
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr1, origCoins))
	suite.Require().NoError(simapp.FundAccount(app.BankKeeper, ctx, addr2, origCoins))

	// Define test cases
	testCases := []struct {
		name          string
		restrictionFn func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error)
		delegatorAddr sdk.AccAddress
		moduleAccAddr sdk.AccAddress
		amt           sdk.Coins
		expectError   error
		action        string // "delegate" or "undelegate"
		accountType   string // "vesting" or "non-vesting"
	}{
		// Vesting account test cases
		{
			name: "Vesting account: Delegate blocked when delegator is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if fromAddr.Equals(addr1) {
					return toAddr, fmt.Errorf("vesting delegatorAddr is restricted")
				}
				return toAddr, nil
			},
			delegatorAddr: addr1,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   fmt.Errorf("vesting delegatorAddr is restricted"),
			action:        "delegate",
			accountType:   "vesting",
		},
		{
			name: "Vesting account: Delegate succeeds when neither is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				return toAddr, nil
			},
			delegatorAddr: addr1,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   nil,
			action:        "delegate",
			accountType:   "vesting",
		},
		{
			name: "Vesting account: Undelegate blocked when delegator is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if fromAddr.Equals(addrModule) && toAddr.Equals(addr1) {
					return toAddr, fmt.Errorf("vesting delegatorAddr is restricted")
				}
				return toAddr, nil
			},
			delegatorAddr: addr1,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   fmt.Errorf("vesting delegatorAddr is restricted"),
			action:        "undelegate",
			accountType:   "vesting",
		},
		{
			name: "Vesting account: Undelegate succeeds when neither is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				return toAddr, nil
			},
			delegatorAddr: addr1,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   nil,
			action:        "undelegate",
			accountType:   "vesting",
		},
		// Non-vesting account test cases
		{
			name: "Non-vesting account: Delegate blocked when delegator is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if fromAddr.Equals(addr2) {
					return toAddr, fmt.Errorf("non-vesting delegatorAddr is restricted")
				}
				return toAddr, nil
			},
			delegatorAddr: addr2,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   fmt.Errorf("non-vesting delegatorAddr is restricted"),
			action:        "delegate",
			accountType:   "non-vesting",
		},
		{
			name: "Non-vesting account: Delegate succeeds when neither is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				return toAddr, nil
			},
			delegatorAddr: addr2,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   nil,
			action:        "delegate",
			accountType:   "non-vesting",
		},
		{
			name: "Non-vesting account: Undelegate blocked when delegator is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				if fromAddr.Equals(addrModule) && toAddr.Equals(addr2) {
					return toAddr, fmt.Errorf("non-vesting delegatorAddr is restricted")
				}
				return toAddr, nil
			},
			delegatorAddr: addr2,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   fmt.Errorf("non-vesting delegatorAddr is restricted"),
			action:        "undelegate",
			accountType:   "non-vesting",
		},
		{
			name: "Non-vesting account: Undelegate succeeds when neither is restricted",
			restrictionFn: func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
				return toAddr, nil
			},
			delegatorAddr: addr2,
			moduleAccAddr: addrModule,
			amt:           delCoins,
			expectError:   nil,
			action:        "undelegate",
			accountType:   "non-vesting",
		},
	}

	// Iterate over each test case
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Set up the BankKeeper with the restriction function
			suite.app.BankKeeper = keeper.NewBaseKeeper(
				suite.app.AppCodec(),
				suite.app.GetKey(types.StoreKey),
				suite.app.AccountKeeper,
				suite.app.GetSubspace(types.ModuleName),
				nil, // Module account is not needed for this test
			).WithSendCoinsRestriction(keeper.SendRestrictionFn(tc.restrictionFn))

			var err error

			// Execute the appropriate action (delegate or undelegate)
			switch tc.action {
			case "delegate":
				err = suite.app.BankKeeper.DelegateCoins(suite.ctx, tc.delegatorAddr, tc.moduleAccAddr, tc.amt)
			case "undelegate":
				err = suite.app.BankKeeper.UndelegateCoins(suite.ctx, tc.moduleAccAddr, tc.delegatorAddr, tc.amt)
			}

			// Check for expected error or success
			if tc.expectError == nil {
				suite.Require().NoError(err, "expected no error but got: %v", err)
			} else {
				suite.Require().EqualError(err, tc.expectError.Error(), "expected error: %v but got: %v", tc.expectError, err)
			}
		})
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
