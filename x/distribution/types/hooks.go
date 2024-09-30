package types

import (
	"context"
)

var _ DistributionHooks = &MultiDistributionHooks{}

type MultiDistributionHooks []DistributionHooks

func NewMultiDistributionHooks(hooks ...DistributionHooks) MultiDistributionHooks {
	return hooks
}

func (h MultiDistributionHooks) ValidatorVotingPowersPerAsset(ctx context.Context) (map[string][2]int64, error) {
	votingPowers := make(map[string][2]int64)
	err := error(nil)
	for i := range h {
		if votingPowers, err = h[i].ValidatorVotingPowersPerAsset(ctx); err != nil {
			return nil, err
		}
	}
	return votingPowers, err
}
