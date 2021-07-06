package types

// Event types for the farming module.
const (
	EventTypeCreateFixedAmountPlan = "create_fixed_amount_plan"
	EventTypeCreateRatioPlan       = "create_ratio_plan"
	EventTypeStake                 = "stake"
	EventTypeUnstake               = "unstake"
	EventTypeClaim                 = "claim"

	AttributeKeyFarmingPoolAddress      = "farming_pool_address"
	AttributeKeyStakingReserveAddress   = "staking_reserve_address"
	AttributeKeyDistributionPoolAddress = "distribution_pool_address"
	AttributeKeyTerminationAddress      = "termination_address"
	AttributeKeyStakingCoinWeights      = "staking_coin_weights"
	AttributeKeyStartTime               = "start_time"
	AttributeKeyEndTime                 = "end_time"
	AttributeKeyEpochDays               = "epoch_days"
	AttributeKeyEpochAmount             = "epoch_amount"
	AttributeKeyEpochRatio              = "epoch_ratio"
)
