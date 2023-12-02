package types_test

import (
	"testing"

	"github.com/babylonchain/babylon/privval"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"
)

func TestProofOfPossession_IsValid(t *testing.T) {
	valPrivKey := ed25519.GenPrivKey()
	blsPrivKey := bls12381.GenPrivKey()
	pop, err := privval.BuildPoP(valPrivKey, blsPrivKey)
	require.NoError(t, err)
	valpk, err := codec.FromTmPubKeyInterface(valPrivKey.PubKey())
	require.NoError(t, err)
	require.True(t, pop.IsValid(blsPrivKey.PubKey(), valpk))
}
