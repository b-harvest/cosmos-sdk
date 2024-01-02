package keeper_test

//import (
//	"math/rand"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/require"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/x/epoching/testepoching"
//	"github.com/cosmos/cosmos-sdk/x/epoching/types"
//	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
//)
//
//// TODO (fuzz tests): replace the following tests with fuzz ones
//func TestMsgWrappedDelegate(t *testing.T) {
//	r := rand.New(rand.NewSource(time.Now().Unix()))
//	helper := testepoching.NewHelper(t)
//	msgSrvr := helper.MsgSrvr
//	// enter 1st epoch, in which BBN starts handling validator-related msgs
//	ctx := helper.GenAndApplyEmptyBlock(r)
//	wctx := sdk.WrapSDKContext(ctx)
//
//	testCases := []struct {
//		name      string
//		req       *stakingtypes.MsgDelegate
//		expectErr bool
//	}{
//		{
//			"empty wrapped msg",
//			&stakingtypes.MsgDelegate{},
//			true,
//		},
//	}
//	for _, tc := range testCases {
//		wrappedMsg := types.NewMsgWrappedDelegate(tc.req)
//		_, err := msgSrvr.WrappedDelegate(wctx, wrappedMsg)
//		if tc.expectErr {
//			require.Error(t, err)
//		} else {
//			require.NoError(t, err)
//		}
//	}
//}
//
//func TestMsgWrappedUndelegate(t *testing.T) {
//	r := rand.New(rand.NewSource(time.Now().Unix()))
//	helper := testepoching.NewHelper(t)
//	msgSrvr := helper.MsgSrvr
//	// enter 1st epoch, in which BBN starts handling validator-related msgs
//	ctx := helper.GenAndApplyEmptyBlock(r)
//	wctx := sdk.WrapSDKContext(ctx)
//
//	testCases := []struct {
//		name      string
//		req       *stakingtypes.MsgUndelegate
//		expectErr bool
//	}{
//		{
//			"empty wrapped msg",
//			&stakingtypes.MsgUndelegate{},
//			true,
//		},
//	}
//	for _, tc := range testCases {
//		wrappedMsg := types.NewMsgWrappedUndelegate(tc.req)
//		_, err := msgSrvr.WrappedUndelegate(wctx, wrappedMsg)
//		if tc.expectErr {
//			require.Error(t, err)
//		} else {
//			require.NoError(t, err)
//		}
//	}
//}
//
//func TestMsgWrappedBeginRedelegate(t *testing.T) {
//	r := rand.New(rand.NewSource(time.Now().Unix()))
//	helper := testepoching.NewHelper(t)
//	msgSrvr := helper.MsgSrvr
//	// enter 1st epoch, in which BBN starts handling validator-related msgs
//	ctx := helper.GenAndApplyEmptyBlock(r)
//	wctx := sdk.WrapSDKContext(ctx)
//
//	testCases := []struct {
//		name      string
//		req       *stakingtypes.MsgBeginRedelegate
//		expectErr bool
//	}{
//		{
//			"empty wrapped msg",
//			&stakingtypes.MsgBeginRedelegate{},
//			true,
//		},
//	}
//	for _, tc := range testCases {
//		wrappedMsg := types.NewMsgWrappedBeginRedelegate(tc.req)
//
//		_, err := msgSrvr.WrappedBeginRedelegate(wctx, wrappedMsg)
//		if tc.expectErr {
//			require.Error(t, err)
//		} else {
//			require.NoError(t, err)
//		}
//	}
//}
