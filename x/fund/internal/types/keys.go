package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// the one key to use for the keeper store
var FundKeyPrefix = []byte{0x00}
var FundByNameKey = []byte{0x01}
var FundsByManagerAddrKey = []byte{0x02}

// nolint
const (
	// module name
	ModuleName = "fund"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
)

// TODO: below functions
// gets the key for delegator bond with validator
// VALUE: staking/Delegation
func GetFundKey(fundID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, fundID)
	return append(FundKeyPrefix, bz...)
}

//
//// SetProposal set a proposal to store
//func (keeper Keeper) SetProposal(ctx sdk.Context, proposal Proposal) {
//	store := ctx.KVStore(keeper.storeKey)
//	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal)
//	store.Set(ProposalKey(proposal.ProposalID), bz)
//}
//// ProposalKey gets a specific proposal from the store
//func ProposalKey(proposalID uint64) []byte {
//	bz := make([]byte, 8)
//	binary.LittleEndian.PutUint64(bz, proposalID)
//	return append(ProposalsKeyPrefix, bz...)
//}

func GetFundKeyByName(fundName string) []byte {
	return append(FundByNameKey, []byte(fundName)...)
}

// gets the prefix for a delegator for all validators
func GetFundsKey(mgrAddr sdk.AccAddress) []byte {
	return append(FundsByManagerAddrKey, mgrAddr.Bytes()...)
}
