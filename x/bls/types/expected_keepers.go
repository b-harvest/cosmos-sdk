package types

import (
	"context"

	"cosmossdk.io/core/address"

	sdk "github.com/cosmos/cosmos-sdk/types"
	epochingtypes "github.com/cosmos/cosmos-sdk/x/epoching/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	AddressCodec() address.Codec
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// EpochingKeeper defines the expected interface needed to retrieve epoch info
type EpochingKeeper interface {
	GetEpoch(ctx sdk.Context) *epochingtypes.Epoch
	EnqueueMsg(ctx sdk.Context, msg epochingtypes.QueuedMessage)
	GetValidatorSet(ctx sdk.Context, epochNumer uint64) epochingtypes.ValidatorSet
	GetTotalVotingPower(ctx sdk.Context, epochNumber uint64) int64
	CheckMsgCreateValidator(ctx sdk.Context, msg *stakingtypes.MsgCreateValidator) error
}

// Event Hooks
// These can be utilized to communicate between a checkpointing keeper and another
// keeper which must take particular actions when raw checkpoints change
// state. The second keeper must implement this interface, which then the
// checkpointing keeper can call.

// CheckpointingHooks event hooks for raw checkpoint object (noalias)
type CheckpointingHooks interface {
	// TODO: bump 50 , sdk -> context
	AfterBlsKeyRegistered(ctx context.Context, valAddr sdk.ValAddress) error         // Must be called when a BLS key is registered
	AfterRawCheckpointConfirmed(ctx context.Context, epoch uint64) error             // Must be called when a raw checkpoint is CONFIRMED
	AfterRawCheckpointForgotten(ctx context.Context, ckpt *RawCheckpoint) error      // Must be called when a raw checkpoint is FORGOTTEN
	AfterRawCheckpointFinalized(ctx context.Context, epoch uint64) error             // Must be called when a raw checkpoint is FINALIZED
	AfterRawCheckpointBlsSigVerified(ctx context.Context, ckpt *RawCheckpoint) error // Must be called when a raw checkpoint's multi-sig is verified
}
