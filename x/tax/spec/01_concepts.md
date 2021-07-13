<!-- order: 1 -->

 # Concepts
## Tax Module

`x/tax` is a Cosmos SDK module that implements tax functionality that collecting taxes from `tax_source_address` and distribute to each `pool_address` at a rate

### 1. Atom Inflation Distribution

- Current : distribution module reward flow

    1) Gas fees collected in ante handler, sent to `feeCollectorName` module account

    2) Atom inflation minted in mint module, sent to `feeCollectorName` module account

    3) In distribution module

    a) Send all rewards in `feeCollectorName` to distribution module account

    b) From `distributionModuleAccount`, substitute `communityTax`

    c) Rest are distributed to proposer and validator reward pools

    d) Substituted amount for community tax is saved in kv store

- **Implementation with Tax Module**
    - Implementation Independency
        - Tax Module is **100% independent** from other existing modules
            - We don’t need to change **any module** at all to adopt Tax Module!
            - Tax Module even works **without** Distribution Module or Mint Module!
        - Begin Block Processing Order : Mint → **(Tax)** → Distribution
    - Functionalities
        - Distribute Atom inflation and gas fees to different tax purposes
            - Atom inflation and gas fees are accumulated in `feeCollectorName` module account
            - Distribute tax amounts from `feeCollectorName` module account to **each tax pool** module account
            - Rest amounts **stay** in `feeCollectorName` so that “Distribution Module” can use it for community fund and staking rewards distribution
        - Create, modify or remove tax taxes via governance
            - Tax taxes can be created, modified or removed by **usual parameter governance**
    - Coin Flow with Tax module
        - In Mint Module
            - Atom inflation to `feeCollectorName` module account
            - [https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/mint/abci.go#L27-L40](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/mint/abci.go#L27-L40)
            - [https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/mint/keeper/keeper.go#L108-L110](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/mint/keeper/keeper.go#L108-L110)
        - In Ante Handler
            - Gas fees to `feeCollectorName` module account
            - [https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/auth/ante/fee.go#L112-L135](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/auth/ante/fee.go#L112-L135)
        - In Tax Module
            - Part of Atom inflation and gas fees go to different tax pools
            - Rest stays within `feeCollectorName` module account
        - In Distribution Module
            - All amounts in `feeCollectorName` module account go to community fund and staking rewards
            - [https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/distribution/keeper/allocation.go#L82-L101](https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-rc0/x/distribution/keeper/allocation.go#L82-L101)

    ![https://s3-us-west-2.amazonaws.com/secure.notion-static.com/702b4723-b59a-4db6-bd41-a3d9344fa394/Untitled.png](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/702b4723-b59a-4db6-bd41-a3d9344fa394/Untitled.png)

    - tax module parameter example

        ```json
        {
          "mint": {
            "params": {
              "taxes": [
                {
                  "name": "liquidity_farming_20213Q-20221Q",
                  "rate": "0.300000000000000000",
                  "pool_address": "xxx",
                  "tax_source_address": "cosmos17xpfvakm2amg962yls6f84z3kell8c5lserqta",
                  "start_time": "2021-10-01T00:00:00Z",
                  "end_time": "2022-04-01T00:00:00Z"
                }
              ]
            }
          }
        }
        ```