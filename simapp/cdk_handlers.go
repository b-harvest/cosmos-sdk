package simapp

import (
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

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
