package privval

import (
	"errors"

	tmcrypto "github.com/cometbft/cometbft/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"
	blstypes "github.com/cosmos/cosmos-sdk/x/bls/types"
)

type ValidatorKeys struct {
	ValPubkey tmcrypto.PubKey
	BlsPubkey bls12381.PublicKey
	PoP       *blstypes.ProofOfPossession

	valPrivkey tmcrypto.PrivKey
	blsPrivkey bls12381.PrivateKey
}

func NewValidatorKeys(valPrivkey tmcrypto.PrivKey, blsPrivKey bls12381.PrivateKey) (*ValidatorKeys, error) {
	pop, err := BuildPoP(valPrivkey, blsPrivKey)
	if err != nil {
		return nil, err
	}
	return &ValidatorKeys{
		ValPubkey:  valPrivkey.PubKey(),
		BlsPubkey:  blsPrivKey.PubKey(),
		valPrivkey: valPrivkey,
		blsPrivkey: blsPrivKey,
		PoP:        pop,
	}, nil
}

// BuildPoP builds a proof-of-possession by PoP=sign(key = BLS_sk, data = sign(key = Ed25519_sk, data = BLS_pk))
// where valPrivKey is Ed25519_sk and blsPrivkey is BLS_sk
func BuildPoP(valPrivKey tmcrypto.PrivKey, blsPrivkey bls12381.PrivateKey) (*blstypes.ProofOfPossession, error) {
	if valPrivKey == nil {
		return nil, errors.New("validator private key is empty")
	}
	if blsPrivkey == nil {
		return nil, errors.New("BLS private key is empty")
	}
	data, err := valPrivKey.Sign(blsPrivkey.PubKey().Bytes())
	if err != nil {
		return nil, err
	}
	pop := bls12381.Sign(blsPrivkey, data)
	return &blstypes.ProofOfPossession{
		Ed25519Sig: data,
		BlsSig:     &pop,
	}, nil
}
