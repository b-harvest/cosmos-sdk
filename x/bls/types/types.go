package types

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	epochingtypes "github.com/cosmos/cosmos-sdk/x/epoching/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// HashSize is the size in bytes of a hash
	HashSize   = sha256.Size
	BitmapBits = 104 // 104 bits for 104 validators at top
)

type LastCommitHash []byte

type BlsSigHash []byte

type RawCkptHash []byte

func NewCheckpoint(epochNum uint64, lch LastCommitHash) *RawCheckpoint {
	return &RawCheckpoint{
		EpochNum:       epochNum,
		LastCommitHash: &lch,
		Bitmap:         bitmap.New(BitmapBits), // 13 bytes, holding 100 validators
		BlsMultiSig:    nil,
	}
}

func NewCheckpointWithMeta(ckpt *RawCheckpoint, status CheckpointStatus) *RawCheckpointWithMeta {
	return &RawCheckpointWithMeta{
		Ckpt:      ckpt,
		Status:    status,
		Lifecycle: []*CheckpointStateUpdate{},
	}
}

// Accumulate does the following things
// 1. aggregates the BLS signature
// 2. aggregates the BLS public key
// 3. updates Bitmap
// 4. accumulates voting power
// it returns nil if the checkpoint is updated, otherwise it returns an error
func (cm *RawCheckpointWithMeta) Accumulate(
	vals epochingtypes.ValidatorSet,
	signerAddr sdk.ValAddress,
	signerBlsKey bls12381.PublicKey,
	sig bls12381.Signature,
	totalPower int64) error {

	// the checkpoint should be accumulating
	if cm.Status != Accumulating {
		return ErrCkptNotAccumulating
	}

	// get validator and its index
	val, index, err := vals.FindValidatorWithIndex(signerAddr)
	if err != nil {
		return err
	}

	// return an error if the validator has already voted
	if bitmap.Get(cm.Ckpt.Bitmap, index) {
		return ErrCkptAlreadyVoted
	}

	// aggregate BLS sig
	if cm.Ckpt.BlsMultiSig != nil {
		aggSig, err := bls12381.AggrSig(*cm.Ckpt.BlsMultiSig, sig)
		if err != nil {
			return err
		}
		cm.Ckpt.BlsMultiSig = &aggSig
	} else {
		cm.Ckpt.BlsMultiSig = &sig
	}

	// aggregate BLS public key
	if cm.BlsAggrPk != nil {
		aggPK, err := bls12381.AggrPK(*cm.BlsAggrPk, signerBlsKey)
		if err != nil {
			return err
		}
		cm.BlsAggrPk = &aggPK
	} else {
		cm.BlsAggrPk = &signerBlsKey
	}

	// update bitmap
	bitmap.Set(cm.Ckpt.Bitmap, index, true)

	// accumulate voting power and update status when the threshold is reached
	cm.PowerSum += uint64(val.Power)
	if int64(cm.PowerSum) > totalPower/3 {
		cm.Status = Sealed
	}

	return nil
}

func (cm *RawCheckpointWithMeta) IsMoreMatureThanStatus(status CheckpointStatus) bool {
	return cm.Status > status
}

// RecordStateUpdate appends a new state update to the raw ckpt with meta
// where the time/height are captured by the current ctx
func (cm *RawCheckpointWithMeta) RecordStateUpdate(ctx sdk.Context, status CheckpointStatus) {
	height, time := ctx.BlockHeight(), ctx.BlockTime()
	stateUpdate := &CheckpointStateUpdate{
		State:       status,
		BlockHeight: uint64(height),
		BlockTime:   &time,
	}
	cm.Lifecycle = append(cm.Lifecycle, stateUpdate)
}

func NewLastCommitHashFromHex(s string) (LastCommitHash, error) {
	bz, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	var lch LastCommitHash

	err = lch.Unmarshal(bz)
	if err != nil {
		return nil, err
	}

	return lch, nil
}

func (lch *LastCommitHash) Unmarshal(bz []byte) error {
	if len(bz) != HashSize {
		return errors.New("invalid lastCommitHash length")
	}
	*lch = bz
	return nil
}

func (lch *LastCommitHash) Size() (n int) {
	if lch == nil {
		return 0
	}
	return len(*lch)
}

func (lch *LastCommitHash) Equal(l LastCommitHash) bool {
	return lch.String() == l.String()
}

func (lch *LastCommitHash) String() string {
	return hex.EncodeToString(*lch)
}

func (lch *LastCommitHash) MustMarshal() []byte {
	bz, err := lch.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func (lch *LastCommitHash) Marshal() ([]byte, error) {
	return *lch, nil
}

func (lch LastCommitHash) MarshalTo(data []byte) (int, error) {
	copy(data, lch)
	return len(data), nil
}

func (lch *LastCommitHash) ValidateBasic() error {
	if lch == nil {
		return errors.New("invalid lastCommitHash")
	}
	if len(*lch) != HashSize {
		return errors.New("invalid lastCommitHash")
	}
	return nil
}

func RawCkptToBytes(cdc codec.BinaryCodec, ckpt *RawCheckpoint) []byte {
	return cdc.MustMarshal(ckpt)
}

// ValidateBasic does sanity checks on a raw checkpoint
func (ckpt RawCheckpoint) ValidateBasic() error {
	if ckpt.Bitmap == nil {
		return ErrInvalidRawCheckpoint.Wrapf("bitmap cannot be empty")
	}
	err := ckpt.LastCommitHash.ValidateBasic()
	if err != nil {
		return ErrInvalidRawCheckpoint.Wrapf(err.Error())
	}
	err = ckpt.BlsMultiSig.ValidateBasic()
	if err != nil {
		return ErrInvalidRawCheckpoint.Wrapf(err.Error())
	}

	return nil
}

func CkptWithMetaToBytes(cdc codec.BinaryCodec, ckptWithMeta *RawCheckpointWithMeta) []byte {
	return cdc.MustMarshal(ckptWithMeta)
}

func BytesToCkptWithMeta(cdc codec.BinaryCodec, bz []byte) (*RawCheckpointWithMeta, error) {
	ckptWithMeta := new(RawCheckpointWithMeta)
	err := cdc.Unmarshal(bz, ckptWithMeta)
	return ckptWithMeta, err
}
