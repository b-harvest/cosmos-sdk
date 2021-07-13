<!-- order: 2 -->

 # State

The tax module keeps track of the staking and rewards states.

## Tax Interface

The tax interface exposes methods to read and write standard tax tax information. Note that all of these methods operate on a tax struct confirming to the interface and in order to write the tax to the store, the tax keeper will need to be used. 

```go
// TaxI is an interface used to store tax records within state.
type TaxI interface {
    proto.Message
    
    GetId() uint64
    SetId(uint64) error
    
    GetType() int32
    SetType(int32) error

    GetTaxPoolAddress() sdk.AccAddress
    SetTaxPoolAddress(sdk.AccAddress) error

    GetRewardPoolAddress() sdk.AccAddress
    SetRewardPoolAddress(sdk.AccAddress) error

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

    String() string
}
```

## Base Tax

A base tax is the simplest and most common tax type, which just stores all requisite fields directly in a struct.

```go
// BaseTax defines a base tax type. It contains all the necessary fields
// for basic tax tax functionality. Any custom tax tax type should extend this
// type for additional functionality (e.g. fixed amount tax, ratio tax).
type BaseTax struct {
    Id                       uint64       // index of the tax
    Type                     TaxType     // type of the tax; public or private
    TaxPoolAddress       string       // bech32-encoded tax pool address
    RewardPoolAddress        string       // bech32-encoded reward pool address
    TerminationAddress       string       // bech32-encoded termination address
    StakingReserveAddress    string       // bech32-encoded staking reserve address
    StakingCoinWeights       sdk.DecCoins // coin weights for the tax
    StartTime                time.Time    // start time of the tax
    EndTime                  time.Time    // end time of the tax
    EpochDays                uint32       // distributing epoch measuring in days
}
```

```go
// FixedAmountTax defines a fixed amount tax that fixed amount of coins are distributed for every epoch day.
type FixedAmountTax struct {
    *BaseTax

    EpochAmount      sdk.Coins // distributing amount for each epoch
}
```

```go
// RatioTax defines a ratio tax that ratio of total coins in tax pool address is distributed for every epoch day.
type RatioTax struct {
    *BaseTax

    EpochRatio            sdk.Dec // distributing amount by ratio
}
```
## Tax Types

```go
// TaxType enumerates the valid types of a tax.
type TaxType int32

const (
    // PLAN_TYPE_UNSPECIFIED defines the default tax type.
    TaxTypeNil TaxType = 0
    // PLAN_TYPE_PUBLIC defines the public tax type.
    TaxTypePublic TaxType = 1
    // PLAN_TYPE_PRIVATE defines the private tax type.
    TaxTypePrivate TaxType = 2
)
```

The parameters of the Tax state are:

- ModuleName, RouterKey, StoreKey, QuerierRoute: `tax`
- Tax: `0x11 | Id -> ProtocolBuffer(Tax)`
- TaxByFarmerAddrIndex: `0x12 | FarmerAddrLen (1 byte) | FarmerAddr -> Id`
    - iterable for several `TaxId` results by indexed `FarmerAddr`
- LastEpochTime: `0x13 | Id -> time.Time`
- GlobalTaxTaxIdKey: `[]byte("globalTaxTaxId") -> ProtocolBuffer(uint64)`
    - store latest tax id
- ModuleName, RouterKey, StoreKey, QuerierRoute: `tax`


## Staking

```go
// Staking defines a farmer's staking information.
type Staking struct {
    TaxId                   uint64
    Farmer                   string
    StakedCoins              sdk.Coins
    QueuedCoins              sdk.Coins
}
```

The parameters of the Staking state are:

- Staking: `0x21 | TaxId | FarmerAddrLen (1 byte) | FarmerAddr -> ProtocolBuffer(Staking)`

## Reward

```go
// Reward defines a record of tax rewards.
type Reward struct {
    TaxId                   uint64
    Farmer                   string
    RewardCoins              sdk.Coins
}
```

The parameters of the Reward state are:

- Reward: `0x31 | TaxId | FarmerAddrLen (1 byte) | FarmerAddr -> ProtocolBuffer(Reward)`

## Examples 

An example of `FixedAmountTax`

```json
{
  "base_tax": {
    "id": 0,
    "type": 0,
    "taxPoolAddress": "cosmos1...",
    "rewardPoolAddress": "cosmos1...",
    "stakingReserveAddress": "cosmos1...",
    "stakingCoinWeights": [
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
    "terminationAddress": "cosmos1..."
  },
  "epochAmount": {
    "denom": "uatom",
    "amount": "10000000"
  }
}
```

An example of `RatioTax`

```json
{
  "base_tax": {
    "id": 0,
    "type": 0,
    "taxPoolAddress": "cosmos1...",
    "rewardPoolAddress": "cosmos1...",
    "stakingReserveAddress": "cosmos1...",
    "stakingCoinWeights": [
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
    "terminationAddress": "cosmos1..."
  },
  "epochRatio": "0.01"
}
```




