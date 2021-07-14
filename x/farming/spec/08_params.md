<!-- order: 8 -->

# Parameters

The farming module contains the following parameters:

| Key                        | Type      | Example                                  |
| -------------------------- | --------- | ---------------------------------------- |
| PrivatePlanCreationFee     | sdk.Coins | [{"denom":"stake","amount":"100000000"}] |
| StakingPositionCreationFee | sdk.Coins | [{"denom":"stake","amount":"100000"}]    |
| EpochDays                  | uint64    | 1                                        |

## PrivatePlanCreationFee

Fee paid for to create a Private type Farming plan. This fee prevents spamming and is collected in in the community pool of the distribution module.

## StakingPositionCreationFee

When a farmer creates new `StakingPosition`, the farmer needs to pay `StakingPositionCreationFee` to prevent spam on the `StakingPosition` struct.

## EpochDays

The universal epoch length in number of days. Every process for staking and reward distribution is executed with this `EpochDays` frequency.
