<!-- order: 3 -->

 # State Transitions

These messages (Msg) in the farming module trigger state transitions.

- Private/Public type Farming plan
    - Private type Farming plan
        - Anyone can create a private `farmingPlan` with his/her own account as `farmingPoolAddress`
        - Plan creator should pay `planCreationFee`
        - Plan creator should provide a signature with his/her own account which will be used for `farmingPoolAddress`
        - Funds inside the `farmingPoolAddress` will be distributed to farmers automatically
    - Public type Farming plan
        - Only governance can create a public `farmingPlan` with a module account as `farmingPoolAddress`
        - Funds inside the `farmingPoolAddress` will be distributed to farmers automatically

    ```go
    type PlanType int32

    const (
        TypePublic  PlanType = 0
        TypePrivate PlanType = 1
    )
    ```

- Staking Coins for Farming
    - Each `farmingPlan` predefines list of `stakingCoinsWithWeight` using `sdk.DecCoins`
    - `weight` mean that each group of stakers with each coin `denom` will receive each predefined `weight` of the total rewards
- Multiple Farming Coins within a `farmingPoolAddress`
    - If `farmingPoolAddress` has multiple kinds of coins, then all coins are identically distributed following the given `farmingPlan`
- Time Parameters
    - Each `farmingPlan` has its own `startTime` and `endTime`
    - Each `farmingPlan` has its own `epochDays` : farming rewards distribution frequency
- Distribution Method
    - `FixedAmountPlan`
        - fixed amount of coins are distributed for each `epochDays`
        - amount in `sdk.Coins`
    - `RatioPlan`
        - `epochRatio` of total assets in `farmingPoolAddress` is distributed for each `epochDays`
        - `epochRatio` in percentage
- Termination Address
    - When the plan ends after the `endTime`, transfer the balance of `farmingPoolAddress` to  `terminationAddress`.