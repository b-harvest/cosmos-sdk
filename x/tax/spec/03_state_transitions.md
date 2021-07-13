<!-- order: 3 -->

 # State Transitions

This document describes the state transaction operations pertaining to the tax module. 

As stated in [01_concepts.md](01_concepts.md), there are public and private tax taxes available in the `tax` module. Public tax can be created by any account whereas private tax can only be created through governance proposal.

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

- Staking Coins for Tax
    - Each `taxTax` predefines list of `stakingCoinWeights` using `sdk.DecCoins`
    - `weight` mean that each group of stakers with each coin `denom` will receive each predefined `weight` of the total rewards
- Multiple Tax Coins within a `taxPoolAddress`
    - If `taxPoolAddress` has multiple kinds of coins, then all coins are identically distributed following the given `taxTax`
- Time Parameters
    - Each `taxTax` has its own `startTime` and `endTime`
    - Each `taxTax` has its own `epochDays` : tax rewards distribution frequency
- Distribution Method
    - `FixedAmountTax`
        - fixed amount of coins are distributed for each `epochDays`
        - amount in `sdk.Coins`
    - `RatioTax`
        - `epochRatio` of total assets in `taxPoolAddress` is distributed for each `epochDays`
        - `epochRatio` in percentage
- Termination Address
    - When the tax ends after the `endTime`, transfer the balance of `taxPoolAddress` to  `terminationAddress`.