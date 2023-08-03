package simulation

// DONTCOVER

import (
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

const (
	keyCommunityTax        = "communitytax"
	keyBaseProposerReward  = "baseproposerreward"
	keyBonusProposerReward = "bonusproposerreward"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	// don't allow param changes
	return nil
}
