package keeper

import (
	"fmt"
	"time"

	epochingtypes "github.com/cosmos/cosmos-sdk/x/epoching/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"
	"github.com/cosmos/cosmos-sdk/x/bls/types"
)

type BlsSigner interface {
	GetAddress() sdk.ValAddress
	SignMsgWithBls(msg []byte) (bls12381.Signature, error)
	GetBlsPubkey() (bls12381.PublicKey, error)
}

// TODO: CDK fix to vote extension

// SendBlsSig prepares a BLS signature message and sends it to Tendermint
func (k Keeper) SendBlsSig(ctx sdk.Context, epochNum uint64, lch types.LastCommitHash, valSet epochingtypes.ValidatorSet) error {
	// get self address
	addr := k.blsSigner.GetAddress()

	// check if itself is the validator
	_, _, err := valSet.FindValidatorWithIndex(addr)
	if err != nil {
		// only send the BLS sig when the node itself is a validator, not being a validator is not an error
		return nil
	}

	// get BLS signature by signing
	signBytes := types.GetSignBytes(epochNum, lch)
	blsSig, err := k.blsSigner.SignMsgWithBls(signBytes)
	if err != nil {
		return err
	}

	// create MsgAddBlsSig message
	msg := types.NewMsgAddBlsSig(k.clientCtx.GetFromAddress(), epochNum, lch, blsSig, addr)

	// keep sending the message to Tendermint until success or timeout
	// TODO should read the parameters from config file
	var res *sdk.TxResponse
	err = retry.Do(1*time.Second, 1*time.Minute, func() error {
		res, err = tx.SendMsgToTendermint(k.clientCtx, msg)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Failed to send the BLS sig tx for epoch %v: %v", epochNum, err))
		return err
	}

	ctx.Logger().Info(fmt.Sprintf("Successfully sent BLS-sig tx for epoch %d, tx hash: %s, gas used: %d, gas wanted: %d", epochNum, res.TxHash, res.GasUsed, res.GasWanted))

	return nil
}
