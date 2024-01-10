package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	epochingtypes "cosmossdk.io/x/epoching/types"

	"cosmossdk.io/x/bls/types"
)

type msgServer struct {
	k Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper}
}

var _ types.MsgServer = msgServer{}

// AddBlsSig adds BLS sig messages and changes a raw checkpoint status to SEALED if sufficient voting power is accumulated
func (m msgServer) AddBlsSig(goCtx context.Context, msg *types.MsgAddBlsSig) (*types.MsgAddBlsSigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(fmt.Sprintf("received BLS sig for epoch %d from %s", msg.BlsSig.EpochNum, msg.GetSigners()))

	err := m.k.addBlsSig(ctx, msg.BlsSig)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddBlsSigResponse{}, nil
}

// TODO: add become bls validator, pop msg

// WrappedCreateValidator registers validator's BLS public key
// and forwards corresponding MsgCreateValidator message to
// the epoching module
func (m msgServer) WrappedCreateValidator(goCtx context.Context, msg *types.MsgWrappedCreateValidator) (*types.MsgWrappedCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// stateless checks on the inside `MsgCreateValidator` msg
	//if err := m.k.epochingKeeper.CheckMsgCreateValidator(ctx, msg.MsgCreateValidator); err != nil {
	//	return nil, err
	//}

	if msg.VerifyPoP() != true {
		return nil, fmt.Errorf("the proof-of-possession is not valid")
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.MsgCreateValidator.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// store BLS public key
	err = m.k.CreateRegistration(ctx, *msg.Key.Pubkey, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.MsgWrappedCreateValidatorResponse{}, err
}

// WrappedCreateValidatorOrigin registers validator's BLS public key
// and forwards corresponding MsgCreateValidator message to
// the epoching module
func (m msgServer) WrappedCreateValidatorOrigin(goCtx context.Context, msg *types.MsgWrappedCreateValidator) (*types.MsgWrappedCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// stateless checks on the inside `MsgCreateValidator` msg
	if err := m.k.epochingKeeper.CheckMsgCreateValidator(ctx, msg.MsgCreateValidator); err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.MsgCreateValidator.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// store BLS public key
	err = m.k.CreateRegistration(ctx, *msg.Key.Pubkey, valAddr)
	if err != nil {
		return nil, err
	}

	// enqueue the msg into the epoching module
	queueMsg := epochingtypes.QueuedMessage{
		Msg: &epochingtypes.QueuedMessage_MsgCreateValidator{MsgCreateValidator: msg.MsgCreateValidator},
	}

	m.k.epochingKeeper.EnqueueMsg(ctx, queueMsg)

	return &types.MsgWrappedCreateValidatorResponse{}, err
}
