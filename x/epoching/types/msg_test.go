package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/x/epoching/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Most of the code below is adapted from https://github.com/cosmos/cosmos-sdk/blob/v0.45.5/x/staking/types/msg_test.go

var (
	DefaultBondDenom = "stake"

	pk1      = ed25519.GenPrivKey().PubKey()
	pk2      = ed25519.GenPrivKey().PubKey()
	pk3      = ed25519.GenPrivKey().PubKey()
	valAddr1 = sdk.ValAddress(pk1.Address())
	valAddr2 = sdk.ValAddress(pk2.Address())
	valAddr3 = sdk.ValAddress(pk3.Address())

	emptyAddr sdk.ValAddress

	coinPos  = sdk.NewInt64Coin(DefaultBondDenom, 1000)
	coinZero = sdk.NewInt64Coin(DefaultBondDenom, 0)
)

func TestMsgDecode(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)
	stakingtypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// pubkey serialisation/deserialisation
	pk1bz, err := cdc.MarshalInterface(pk1)
	require.NoError(t, err)
	var pkUnmarshaled cryptotypes.PubKey
	err = cdc.UnmarshalInterface(pk1bz, &pkUnmarshaled)
	require.NoError(t, err)
	require.True(t, pk1.Equals(pkUnmarshaled.(*ed25519.PubKey)))

	// create unwrapped msg
	msgUnwrapped := stakingtypes.NewMsgDelegate(sdk.AccAddress(valAddr1).String(), valAddr2.String(), coinPos)

	// wrap and marshal msg
	msg := types.NewMsgWrappedDelegate(msgUnwrapped)
	msgSerialized, err := cdc.MarshalInterface(msg)
	require.NoError(t, err)

	var msgUnmarshaled sdk.Msg
	err = cdc.UnmarshalInterface(msgSerialized, &msgUnmarshaled)
	require.NoError(t, err)
	msg2, ok := msgUnmarshaled.(*types.MsgWrappedDelegate)
	require.True(t, ok)
	require.Equal(t, msg.Msg.Amount, msg2.Msg.Amount)
	require.Equal(t, msg.Msg.DelegatorAddress, msg2.Msg.DelegatorAddress)
	require.Equal(t, msg.Msg.ValidatorAddress, msg2.Msg.ValidatorAddress)

	var qmsgUnmarshaled sdk.Msg
	var msgCreateValUnmarshaled sdk.Msg

	commission1 := stakingtypes.NewCommissionRates(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec())
	msgcreateval1, err := stakingtypes.NewMsgCreateValidator(valAddr1.String(), pk1, coinPos, stakingtypes.Description{}, commission1, math.OneInt())
	require.NoError(t, err)
	qmsg, err := types.NewQueuedMessage(1, time.Now(), []byte("tx id 1"), msgcreateval1)
	require.NoError(t, err)
	msgCreateval1Ser, err := cdc.MarshalInterface(msgcreateval1)
	require.NoError(t, err)
	err = cdc.UnmarshalInterface(msgCreateval1Ser, &msgCreateValUnmarshaled)
	require.NoError(t, err)
	msgcreateval3 := msgCreateValUnmarshaled.(*stakingtypes.MsgCreateValidator)
	require.NotNil(t, msgcreateval3.Pubkey.GetCachedValue())

	qmsgSer, err := cdc.MarshalInterface(&qmsg)
	require.NoError(t, err)
	err = cdc.UnmarshalInterface(qmsgSer, &qmsgUnmarshaled)
	qmsg2, ok := qmsgUnmarshaled.(*types.QueuedMessage)
	msgcreateval2 := qmsg2.UnwrapToSdkMsg().(*stakingtypes.MsgCreateValidator)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, qmsg.MsgId, qmsg2.MsgId)
	require.True(t, msgcreateval1.Pubkey.Equal(msgcreateval2.Pubkey))
}

//// test ValidateBasic for MsgWrappedDelegate
//func TestMsgWrappedDelegate(t *testing.T) {
//	tests := []struct {
//		name          string
//		delegatorAddr sdk.AccAddress
//		validatorAddr sdk.ValAddress
//		bond          sdk.Coin
//		expectPass    bool
//	}{
//		{"basic good", sdk.AccAddress(valAddr1), valAddr2, coinPos, true},
//		{"no wrapped msg", nil, nil, coinPos, false},
//		{"self bond", sdk.AccAddress(valAddr1), valAddr1, coinPos, true},
//		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, coinPos, false},
//		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, coinPos, false},
//		{"empty bond", sdk.AccAddress(valAddr1), valAddr2, coinZero, false},
//		{"nil bold", sdk.AccAddress(valAddr1), valAddr2, sdk.Coin{}, false},
//	}
//
//	for _, tc := range tests {
//		var msg *types.MsgWrappedDelegate
//		if tc.delegatorAddr == nil {
//			msg = types.NewMsgWrappedDelegate(nil)
//		} else {
//			msgUnwrapped := stakingtypes.NewMsgDelegate(tc.delegatorAddr.String(), tc.validatorAddr.String(), tc.bond)
//			msg = types.NewMsgWrappedDelegate(msgUnwrapped)
//		}
//		if tc.expectPass {
//			require.NoError(t, msg.ValidateBasic(), "test: %v", tc.name)
//		} else {
//			require.Error(t, msg.ValidateBasic(), "test: %v", tc.name)
//		}
//	}
//}
//
//// test ValidateBasic for MsgWrappedBeginRedelegate
//func TestMsgWrappedBeginRedelegate(t *testing.T) {
//	tests := []struct {
//		name             string
//		delegatorAddr    sdk.AccAddress
//		validatorSrcAddr sdk.ValAddress
//		validatorDstAddr sdk.ValAddress
//		amount           sdk.Coin
//		expectPass       bool
//	}{
//		{"regular", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(DefaultBondDenom, 1), true},
//		{"no wrapped msg", nil, nil, nil, coinPos, false},
//		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(DefaultBondDenom, 0), false},
//		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.Coin{}, false},
//		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, valAddr3, sdk.NewInt64Coin(DefaultBondDenom, 1), false},
//		{"empty source validator", sdk.AccAddress(valAddr1), emptyAddr, valAddr3, sdk.NewInt64Coin(DefaultBondDenom, 1), false},
//		{"empty destination validator", sdk.AccAddress(valAddr1), valAddr2, emptyAddr, sdk.NewInt64Coin(DefaultBondDenom, 1), false},
//	}
//
//	for _, tc := range tests {
//		var msg *types.MsgWrappedBeginRedelegate
//		if tc.delegatorAddr == nil {
//			msg = types.NewMsgWrappedBeginRedelegate(nil)
//		} else {
//			msgUnwrapped := stakingtypes.NewMsgBeginRedelegate(tc.delegatorAddr.String(), tc.validatorSrcAddr.String(), tc.validatorDstAddr.String(), tc.amount)
//			msg = types.NewMsgWrappedBeginRedelegate(msgUnwrapped)
//		}
//		if tc.expectPass {
//			require.NoError(t, msg.ValidateBasic(), "test: %v", tc.name)
//		} else {
//			require.Error(t, msg.ValidateBasic(), "test: %v", tc.name)
//		}
//	}
//}
//
//// test ValidateBasic for MsgWrappedUndelegate
//func TestMsgWrappedUndelegate(t *testing.T) {
//	tests := []struct {
//		name          string
//		delegatorAddr sdk.AccAddress
//		validatorAddr sdk.ValAddress
//		amount        sdk.Coin
//		expectPass    bool
//	}{
//		{"regular", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(DefaultBondDenom, 1), true},
//		{"no wrapped msg", nil, nil, coinPos, false},
//		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(DefaultBondDenom, 0), false},
//		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, sdk.Coin{}, false},
//		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, sdk.NewInt64Coin(DefaultBondDenom, 1), false},
//		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, sdk.NewInt64Coin(DefaultBondDenom, 1), false},
//	}
//
//	for _, tc := range tests {
//		var msg *types.MsgWrappedUndelegate
//		if tc.delegatorAddr == nil {
//			msg = types.NewMsgWrappedUndelegate(nil)
//		} else {
//			msgUnwrapped := stakingtypes.NewMsgUndelegate(tc.delegatorAddr.String(), tc.validatorAddr.String(), tc.amount)
//			msg = types.NewMsgWrappedUndelegate(msgUnwrapped)
//		}
//		if tc.expectPass {
//			require.NoError(t, msg.ValidateBasic(), "test: %v", tc.name)
//		} else {
//			require.Error(t, msg.ValidateBasic(), "test: %v", tc.name)
//		}
//	}
//}
