package bls

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/core/appmodule"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/x/bls/types"

	"cosmossdk.io/x/bls/keeper"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block.
// Upon each BeginBlock, if reaching the second block after the epoch begins, then
// - extract the LastCommitHash from the block
// - create a raw checkpoint with the status of ACCUMULATING
// - start a BLS signer which creates a BLS sig transaction and distributes it to the network
func BeginBlocker(cctx context.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	ctx := sdk.UnwrapSDKContext(cctx)

	// if this block is the second block of an epoch
	epoch := k.GetEpoch(ctx)
	if epoch.IsFirstBlock(ctx) {
		err := k.InitValidatorBLSSet(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to store validator BLS set: %w", err))
		}
	}
	if epoch.IsSecondBlock(ctx) {
		// note that this epochNum is obtained after the BeginBlocker of the epoching module is executed
		// meaning that the epochNum has been incremented upon a new epoch
		lch := ctx.BlockHeader().LastCommitHash
		ckpt, err := k.BuildRawCheckpoint(ctx, epoch.EpochNumber-1, lch)
		if err != nil {
			panic("failed to generate a raw checkpoint")
		}

		// emit BeginEpoch event
		err = ctx.EventManager().EmitTypedEvent(
			&types.EventCheckpointAccumulating{
				Checkpoint: ckpt,
			},
		)
		if err != nil {
			panic(err)
		}
		curValSet := k.GetValidatorSet(ctx, epoch.EpochNumber-1)

		go func() {
			err := k.SendBlsSig(ctx, epoch.EpochNumber-1, lch, curValSet)
			if err != nil {
				// failing to send a BLS-sig causes a panicking situation
				panic(err)
			}
		}()
	}
}

func PreBlocker(ctx context.Context, k keeper.Keeper) (appmodule.ResponsePreBlock, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := k.Logger(ctx)
	cp := sdkCtx.ConsensusParams()
	currentHeight := sdkCtx.HeaderInfo().Height
	targetEnableHeight := ToSetVoteExtensionEnableHeight(cp, currentHeight)
	if cp.Abci.VoteExtensionsEnableHeight != targetEnableHeight {
		cp.Abci.VoteExtensionsEnableHeight = targetEnableHeight
		if err := k.ParamManager.StoreConsensusParams(sdkCtx, cp); err != nil {
			panic(err)
			// TODO: panic handling
			return &sdk.ResponsePreBlock{
				ConsensusParamsChanged: false,
			}, err
		}
		logger.Info("[VE] Changed for %d, on %d, from %d", targetEnableHeight, currentHeight, cp.Abci.VoteExtensionsEnableHeight)
		return &sdk.ResponsePreBlock{
			ConsensusParamsChanged: true,
		}, nil
	}
	return &sdk.ResponsePreBlock{
		ConsensusParamsChanged: false,
	}, nil
}

func ToSetVoteExtensionEnableHeight(cp cmtproto.ConsensusParams, currentHeight int64) int64 {
	nextHeight := currentHeight + 1
	IsCheckpointNextHeight := IsCheckpointHeight(nextHeight)

	IsExtsEnabled := func(height int64) bool {
		return cp.Abci != nil && height >= cp.Abci.VoteExtensionsEnableHeight && cp.Abci.VoteExtensionsEnableHeight != 0
	}

	extsEnabled := IsExtsEnabled(nextHeight)
	if IsCheckpointNextHeight {
		if !extsEnabled {
			// enable for next height
			return nextHeight
		} else {
			// TODO: need to consider
			// stay current param
			return cp.Abci.VoteExtensionsEnableHeight
		}
	} else {
		if !extsEnabled {
			// stay current param
			return cp.Abci.VoteExtensionsEnableHeight
		} else {
			// disable
			return 0
		}
	}
}

// TODO: temporary mocking funcion, TBD
func IsCheckpointHeight(height int64) bool {
	if height%10 == 0 {
		return true
	} else {
		return false
	}
}
