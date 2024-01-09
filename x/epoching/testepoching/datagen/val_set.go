package datagen

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	blstypes "cosmossdk.io/x/bls/types"

	epochingtypes "cosmossdk.io/x/epoching/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"
)

func GenRandomValSet(n int) epochingtypes.ValidatorSet {
	power := int64(10)
	var valSet []epochingtypes.Validator
	for i := 0; i < n; i++ {
		address := GenRandomValidatorAddress()
		val := epochingtypes.Validator{
			Addr:  address,
			Power: power,
		}
		valSet = append(valSet, val)
	}

	return epochingtypes.NewSortedValidatorSet(valSet)
}

func GenRandomValidatorAddress() sdk.ValAddress {
	return sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address())
}

func GenRandomPubkeysAndSigs(n int, msg []byte) ([]bls12381.PublicKey, []bls12381.Signature) {
	var blsPubkeys []bls12381.PublicKey
	var blsSigs []bls12381.Signature
	for i := 0; i < n; i++ {
		privKey := bls12381.GenPrivKey()
		pubkey := bls12381.GenPrivKey().PubKey()
		sig := bls12381.Sign(privKey, msg)
		blsPubkeys = append(blsPubkeys, pubkey)
		blsSigs = append(blsSigs, sig)
	}

	return blsPubkeys, blsSigs
}

func GenerateValidatorSetWithBLSPrivKeys(n int) (*blstypes.ValidatorWithBlsKeySet, []bls12381.PrivateKey) {
	valSet := &blstypes.ValidatorWithBlsKeySet{
		ValSet: make([]*blstypes.ValidatorWithBlsKey, n),
	}
	blsPrivKeys := make([]bls12381.PrivateKey, n)

	for i := 0; i < n; i++ {
		addr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
		blsPrivkey := bls12381.GenPrivKey()
		val := &blstypes.ValidatorWithBlsKey{
			ValidatorAddress: addr.String(),
			BlsPubKey:        blsPrivkey.PubKey(),
			VotingPower:      1000,
		}
		valSet.ValSet[i] = val
		blsPrivKeys[i] = blsPrivkey
	}

	return valSet, blsPrivKeys
}
