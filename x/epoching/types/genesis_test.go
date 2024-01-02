package types_test

//import (
//	"testing"
//
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//	"github.com/cosmos/cosmos-sdk/app"
//	"github.com/cosmos/cosmos-sdk/testutil/nullify"
//	"github.com/stretchr/testify/require"
//
//	"github.com/cosmos/cosmos-sdk/x/epoching"
//	"github.com/cosmos/cosmos-sdk/x/epoching/types"
//)
//
//func TestGenesis(t *testing.T) {
//	// This test requires setting up the staking module
//	// Otherwise the epoching module cannot initialise the genesis validator set
//	app := app.Setup(t, false)
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
//	keeper := app.EpochingKeeper
//
//	genesisState := types.GenesisState{
//		Params: types.DefaultParams(),
//	}
//
//	epoching.InitGenesis(ctx, keeper, genesisState)
//	got := epoching.ExportGenesis(ctx, keeper)
//	require.NotNil(t, got)
//
//	nullify.Fill(&genesisState)
//	nullify.Fill(got)
//}
//
//func TestGenesisState_Validate(t *testing.T) {
//	for _, tc := range []struct {
//		desc     string
//		genState *types.GenesisState
//		valid    bool
//	}{
//		{
//			desc:     "default is valid",
//			genState: types.DefaultGenesis(),
//			valid:    true,
//		},
//		{
//			desc: "valid genesis state",
//			genState: &types.GenesisState{
//				Params: types.Params{
//					EpochInterval: 100,
//				},
//			},
//			valid: true,
//		},
//	} {
//		t.Run(tc.desc, func(t *testing.T) {
//			err := tc.genState.Validate()
//			if tc.valid {
//				require.NoError(t, err)
//			} else {
//				require.Error(t, err)
//			}
//		})
//	}
//}
