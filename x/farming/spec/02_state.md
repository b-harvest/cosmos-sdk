<!-- order: 2 -->

 # State

The farming module `x/farming` keeps...

## Plan

Plan stores information about the farming plan.

Plan type has the following structure.

```go
type PlanI interface {
    proto.Message
    
    GetId() uint64
    SetId(uint64) error
    
    GetType() int32
    SetType(int32) error

    GetFarmingPoolAddress() sdk.AccAddress
    SetFarmingPoolAddress(sdk.AccAddress) error

    GetDistributionPoolAddress() sdk.AccAddress
    SetDistributionPoolAddress(sdk.AccAddress) error

    GetTerminationAddress() sdk.AccAddress
    SetTerminationAddress(sdk.AccAddress) error
    
    GetStakingReserveAddress() sdk.AccAddress
    SetStakingReserveAddress(sdk.AccAddress) error
        
    GetStakingCoinsWeight() sdk.DecCoins
    SetStakingCoinsWeight(sdk.DecCoins) error
    
    GetStartTime() time.Time
    SetStartTime(time.Time) error

    GetEndTime() time.Time
    SetEndTime(time.Time) error

    GetEpochDays() uint32
    SetEpochDays(uint32) error

    GetDistributionMethod() string
    GetDistributionThisEpoch() sdk.Coins
    IsTerminated() bool

    String() string
}
```

```go
type BasePlan struct {
    Id                       uint64          // index of this farming plan
    Type                     int32           // type of this farming plan, private/public
    FarmingPoolAddress       string          // 
    DistributionPoolAddress  string          // 
    TerminationAddress       string
    StakingReserveAddress    string     
    StakingCoinsWeight       []CoinWeight
    StartTime                time.Time
    EndTime                  time.Time
    EpochDays                uint32
}
```

```go
type FixedAmountPlan struct {
    *BasePlan

    EpochAmount      sdk.Coins
}

```

```go
type RatioPlan struct {
    *BasePlan

    EpochRatio            sdk.Dec
}

```

```go
type CoinWeight sdk.DecCoin
type PlanType int32

const (
    TypePublic  PlanType = 0
    TypePrivate PlanType = 1
)
```

The parameters of the Plan state are:

- Plan: `0x11 | Id -> ProtocolBuffer(Plan)`
- PlanByFarmerAddrIndex: `0x12 | FarmerAddrLen (1 byte) | FarmerAddr -> Id`
    - iterable for several `PlanId` results by indexed `FarmerAddr`
- LastEpochTime: `0x13 | Id -> time.Time`
- GlobalFarmingPlanIdKey: `[]byte("globalFarmingPlanId") -> LatestPlanId`
- ModuleName, RouterKey, StoreKey, QuerierRoute: `farming`

- example of `FixedAmountPlan`

    ```json
    {
      "base_plan": {
        "id": 0,
        "type": 0,
        "farmingPoolAddress": "xxx",
        "distributionPoolAddress": "yyy",
        "stakingReserveAddress": "kkk",
        "stakingCoinsWithWeight": [
          {
            "denom": "xxx",
            "amount": "0.200000000000000000"
          },
          {
            "denom": "yyy",
            "amount": "0.300000000000000000"
          },
          {
            "denom": "zzz",
            "amount": "0.500000000000000000"
          }
        ],
        "startTime": "2021-10-01T00:00:00Z",
        "endTime": "2022-04-01T00:00:00Z",
        "epochDays": 1,
        "terminationAddress": "zzz"
      },
      "epochAmount": {
        "denom": "uatom",
        "amount": "10000000"
      }
    }
    ```

- example of `RatioPlan`

    ```json
    {
      "base_plan": {
        "id": 0,
        "type": 0,
        "farmingPoolAddress": "xxx",
        "distributionPoolAddress": "yyy",
        "stakingReserveAddress": "kkk",
        "stakingCoinsWithWeight": [
          {
            "denom": "xxx",
            "amount": "0.200000000000000000"
          },
          {
            "denom": "yyy",
            "amount": "0.300000000000000000"
          },
          {
            "denom": "zzz",
            "amount": "0.500000000000000000"
          }
        ],
        "startTime": "2021-10-01T00:00:00Z",
        "endTime": "2022-04-01T00:00:00Z",
        "epochDays": 1,
        "terminationAddress": "zzz"
      },
      "epochRatio": "0.01"
    }
    ```

```go
type Staking struct {
    PlanId                   uint64
    Farmer                   string
    StakedCoins              sdk.Coins
    QueuedCoins              sdk.Coins
}
```

The parameters of the Staking state are:

- Staking: `0x21 | PlanId | FarmerAddrLen (1 byte) | FarmerAddr -> ProtocolBuffer(Staking)`

```go
type Reward struct {
    PlanId                   uint64
    Farmer                   string
    RewardCoins              sdk.Coins
}
```

The parameters of the Reward state are:

- Reward: `0x31 | PlanId | FarmerAddrLen (1 byte) | FarmerAddr -> ProtocolBuffer(Reward)`



