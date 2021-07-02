<!-- order: 4 -->

 # Messages

Messages (Msg) are objects that trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps farming module messages from transactions.

## **MsgCreatePlan**

```go
type MsgCreateFixedAmountPlan struct {
    FarmingPoolAddress  string
    StakingCoinsWeight  sdk.DecCoins
    StartTime           time.Time
    EndTime             time.Time
    EpochDays           uint32
		EpochAmount         sdk.Coins
}
```

```go
type MsgCreateRatioPlan struct {
    FarmingPoolAddress  string
    StakingCoinsWeight  sdk.DecCoins
    StartTime           time.Time
    EndTime             time.Time
    EpochDays           uint32
		EpochRatio          sdk.Dec
}
```

## **MsgStake**

- To become eligible farmer, he/she should stake staking coins into the farming plan

```go
type MsgStake struct {
    PlanId       uint64
    Farmer       string
    StakingCoins sdk.Coins
}
```

## **MsgUnstake**

- If farmer wants to terminate staking, he/she can unstake the coins
- All accumulated farming rewards is automatically withdrawn to farmer when unstake

```go
type MsgUnstake struct {
    PlanId         uint64
    Farmer         string
    UnstakingCoins sdk.Coins
}

```