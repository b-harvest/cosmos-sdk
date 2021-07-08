package types_test

import (
	fmt "fmt"
	"testing"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
)

const (
	DefaultEpochDays = uint32(1)
)

func TestMsgCreatePool(t *testing.T) {
	farmingPoolAddr := sdk.AccAddress(crypto.AddressHash([]byte("farmingPoolAddr")))
	stakingCoinWeights := sdk.NewDecCoins(
		sdk.DecCoin{Denom: "testFarmStakingCoinDenom", Amount: sdk.MustNewDecFromStr("1.0")},
	)
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(1, 0, 0)

	testCases := []struct {
		name          string
		msg           *types.MsgCreateFixedAmountPlan
		expectedError error
	}{
		{
			name: "happy case",
			msg: types.NewMsgCreateFixedAmountPlan(
				farmingPoolAddr, stakingCoinWeights, startTime, endTime,
				DefaultEpochDays, sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1))},
			),
			expectedError: nil,
		},
		{
			name: "invalid farming pool address",
			msg: types.NewMsgCreateFixedAmountPlan(
				sdk.AccAddress{}, stakingCoinWeights, startTime, endTime,
				DefaultEpochDays, sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1))},
			),
			expectedError: sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid farming pool address %q: %v", "empty address string is not allowed", nil),
		},
		// {
		// 	"epoch amount must not be empty",
		// 	types.NewMsgCreateFixedAmountPlan(
		// 		farmingPoolAddr, stakingCoinWeights, startTime, endTime,
		// 		DefaultEpochDays, sdk.Coins{},
		// 	),
		// },
		// {
		// 	"empty staking coin weights ",
		// 	types.NewMsgCreateFixedAmountPlan(
		// 		sdk.AccAddress{}, stakingCoinWeights, startTime, endTime,
		// 		DefaultEpochDays, sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1))},
		// 	),
		// },
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCreateFixedAmountPlan{}, tc.msg)
		require.Equal(t, types.TypeMsgCreateFixedAmountPlan, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()

		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedError == nil {
				require.Nil(t, err)
				signers := tc.msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, tc.msg.GetPlanCreator(), signers[0])
			} else {
				fmt.Println("tc.expectedError.Error(): ", tc.expectedError.Error())
				require.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}

	// for _, tc := range testCases {
	// 	require.IsType(t, &types.MsgCreateFixedAmountPlan{}, tc.msg)
	// 	require.Equal(t, types.TypeMsgCreateFixedAmountPlan, tc.msg.Type())
	// 	require.Equal(t, types.RouterKey, tc.msg.Route())
	// 	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

	// 	err := tc.msg.ValidateBasic()
	// 	if tc.expectedError != nil {
	// 		require.Nil(t, err)
	// 		signers := tc.msg.GetSigners()
	// 		require.Len(t, signers, 1)
	// 		require.Equal(t, tc.msg.GetPlanCreator(), signers[0])

	// 	} else {
	// 		// require.ErrorIs(t require.TestingT, err error, target error, msgAndArgs ...interface{})
	// 		// require.ErrorAs(t, err, tc.expectedErr)
	// 		fmt.Println("tc.expectedError.Error(): ", tc.expectedError.Error())
	// 		require.EqualError(t, err, tc.expectedError.Error())
	// 	}
	// }
}
