<!-- order: 4 -->

 # Messages

Messages (Msg) are objects that trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps farming module messages from transactions.

## **MsgCreateFixedAmountPlan**

```go
type MsgCreateFixedAmountPlan struct {
    FarmingPoolAddress  string
    StakingCoinWeights  sdk.DecCoins
    StartTime           time.Time
    EndTime             time.Time
    EpochDays           uint32
    EpochAmount         sdk.Coins
}
```

## **MsgCreateRatioPlan**

```go
type MsgCreateRatioPlan struct {
    FarmingPoolAddress  string
    StakingCoinWeights  sdk.DecCoins
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

- If farmer wants to terminate staking, the farmer can unstake the coins at any time
- There is no unbonding period of time required to unstake coins from the plan
- All accumulated farming rewards are automatically withdrawn to farmer when unstake event occurs

```go
type MsgUnstake struct {
    PlanId         uint64
    Farmer         string
    UnstakingCoins sdk.Coins
}

```

## **MsgClaim**

- A farmer should claim their farming rewards. This is similar mechanism with claiming rewards from distribution module.

```go
type MsgClaim struct {
	PlanId uint64
	Farmer string
}
```