package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding tax type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.TaxKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.TaxByFarmerAddrIndexKeyPrefix):
			var pA, pB types.BaseTax
			cdc.MustUnmarshal(kvA.Value, &pA)
			cdc.MustUnmarshal(kvA.Value, &pB)
			return fmt.Sprintf("%v\n%v", pA, pB)

		case bytes.Equal(kvA.Key[:1], types.StakingKeyPrefix):
			var sA, sB types.Staking
			cdc.MustUnmarshal(kvA.Value, &sA)
			cdc.MustUnmarshal(kvA.Value, &sB)
			return fmt.Sprintf("%v\n%v", sA, sB)

		case bytes.Equal(kvA.Key[:1], types.RewardKeyPrefix):
			var rA, rB types.Reward
			cdc.MustUnmarshal(kvA.Value, &rA)
			return fmt.Sprintf("%v\n%v", rA, rB)

		default:
			panic(fmt.Sprintf("invalid tax key prefix %X", kvA.Key[:1]))
		}
	}
}
