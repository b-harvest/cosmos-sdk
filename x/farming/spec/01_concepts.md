<!-- order: 1 -->

 # Concepts
## Farming Module

The `farming` module provides farming rewards to farmers. 

## Plans

In the `farming` module, there are two types of farming plans in the `farming` module.

### Public Farming Plan

A public farming plan can only be created through governance proposal. 

### Private Farming Plan

 A private farming plan can be created with any account. The plan creator's account is used as distributing account `FarmingPoolAddress` that will be distributed to farmers automatically. There is a fee `PlanCreationFee` paid upon plan creation to prevent from spamming attack. 

## Distribution Methods

In the `farming` module, there are two types of distribution methods. 
### Fixed Amount Plan

A `FixedAmountPlan` distributes fixed amount of coins to farmers for every epoch day. 
If the plan creators's account `FarmingPoolAddress` is depleted, then there is no more coins to distribute unless it is filled up with more coins.

### Ratio Plan

A `RatioPlan` distributes to farmers by ratio for every epoch day. If the plan creators's account `FarmingPoolAddress` is depleted, then there is no more coins to distribute unless it is filled up with more coins.

