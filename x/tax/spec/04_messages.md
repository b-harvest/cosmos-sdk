<!-- order: 4 -->

 # Messages

Messages (Msg) are objects that trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps tax module messages from transactions.

## MsgCreateFixedAmountTax

```go
type MsgCreateFixedAmountTax struct {
    TaxPoolAddress  string
    StakingCoinWeights  sdk.DecCoins
    StartTime           time.Time
    EndTime             time.Time
    EpochDays           uint32
    EpochAmount         sdk.Coins
}
```

## MsgCreateRatioTax

```go
type MsgCreateRatioTax struct {
    TaxPoolAddress  string
    StakingCoinWeights  sdk.DecCoins
    StartTime           time.Time
    EndTime             time.Time
    EpochDays           uint32
    EpochRatio          sdk.Dec
}
```
## MsgStake

A farmer must have sufficient coins to stake into a tax tax. The farmer becomes eligible to receive rewards once the farmer stakes some coins.

```go
type MsgStake struct {
    TaxId       uint64
    Farmer       string
    StakingCoins sdk.Coins
}
```
## MsgUnstake

A farmer must have some staking coins in the tax to trigger this message. Unlike `x/staking` module, there is no unbonding period of time required to unstake coins from the tax. All accumulated tax rewards are automatically withdrawn to the farmer once unstaking event is triggered.

```go
type MsgUnstake struct {
    TaxId         uint64
    Farmer         string
    UnstakingCoins sdk.Coins
}

```
## MsgClaim

A farmer should claim their tax rewards. The rewards are not automatically distributed. This is similar mechanism with `x/distribution` module.

```go
type MsgClaim struct {
	TaxId uint64
	Farmer string
}
```