package datagen

import (
	ed255192 "github.com/cometbft/cometbft/crypto/ed25519"

	"cosmossdk.io/x/bls/types"

	"cosmossdk.io/privval"

	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GenerateGenesisKey() *types.GenesisKey {
	accPrivKey := secp256k1.GenPrivKey()
	tmValPrivKey := ed255192.GenPrivKey()
	blsPrivKey := bls12381.GenPrivKey()
	tmValPubKey := tmValPrivKey.PubKey()
	valPubKey, err := codec.FromTmPubKeyInterface(tmValPubKey)
	if err != nil {
		panic(err)
	}

	blsPubKey := blsPrivKey.PubKey()
	address := sdk.ValAddress(accPrivKey.PubKey().Address())
	pop, err := privval.BuildPoP(tmValPrivKey, blsPrivKey)
	if err != nil {
		panic(err)
	}

	gk, err := types.NewGenesisKey(address, &blsPubKey, pop, valPubKey)
	if err != nil {
		panic(err)
	}

	return gk
}
