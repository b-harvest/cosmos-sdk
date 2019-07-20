package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//codec "github.com/cosmos/cosmos-sdk/x/fund/internal/types/codec"
)

type AssetRatePair struct {
	Asset sdk.Coin
	Rate  sdk.Dec
}

//type AssetRatePairs []AssetRatePair
type AssetRatePairs sdk.Coins

//	micro like uatom, need to * 1,000,000

type Fund struct {
	Name           string         `json:"name" yaml:"name"` // TODO: need to checking length limit
	ManagerAccount sdk.AccAddress `json:"manager_account" yaml:"manager_account"`
	FundID         uint64         `json:"fund_id" yaml:"fund_id"`
	Supply         sdk.Int        `json:"supply" yaml:"supply"`
	AssetRatePairs AssetRatePairs
	//AssetRates sdk.Coins
}

// input

// TODO: MustMarshalFund, MustUnMarshalFund

//func NewFund()

// return the fund
func MustMarshalFund(cdc *codec.Codec, fund Fund) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(fund)
}

// return the fund
func MustUnmarshalFund(cdc *codec.Codec, value []byte) Fund {
	fund, err := UnmarshalFund(cdc, value)
	if err != nil {
		panic(err)
	}
	return fund
}

// return the fund
func UnmarshalFund(cdc *codec.Codec, value []byte) (fund Fund, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &fund)
	return fund, err
}
