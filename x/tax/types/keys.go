package types

const (
	// ModuleName is the name of the tax module
	ModuleName = "tax"

	// RouterKey is the message router key for the tax module
	RouterKey = ModuleName

	// StoreKey is the default store key for the tax module
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the tax module
	QuerierRoute = ModuleName
)

//var (
//	// param key for global tax tax IDs
//	GlobalTaxTaxIDKey = []byte("globalTaxTaxId")
//
//	TaxKeyPrefix                  = []byte{0x11}
//	TaxByFarmerAddrIndexKeyPrefix = []byte{0x12}
//	LastEpochTimeKeyPrefix         = []byte{0x13}
//
//	StakingKeyPrefix = []byte{0x21}
//
//	RewardKeyPrefix = []byte{0x31}
//)
//
//// GetTaxKey returns kv indexing key of the tax
//func GetTaxKey(taxID uint64) []byte {
//	key := make([]byte, 9)
//	key[0] = TaxKeyPrefix[0]
//	copy(key[1:], sdk.Uint64ToBigEndian(taxID))
//	return key
//}
//
//// GetTaxesByFarmerAddrIndexKey returns kv indexing key of the tax indexed by reserve account
//func GetTaxesByFarmerAddrIndexKey(farmerAcc sdk.AccAddress) []byte {
//	return append(TaxByFarmerAddrIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...)
//}
//
//// GetTaxByFarmerAddrIndexKey returns kv indexing key of the tax indexed by reserve account
//func GetTaxByFarmerAddrIndexKey(farmerAcc sdk.AccAddress, taxID uint64) []byte {
//	return append(append(TaxByFarmerAddrIndexKeyPrefix, address.MustLengthPrefix(farmerAcc.Bytes())...), sdk.Uint64ToBigEndian(taxID)...)
//}
//
//// GetStakingPrefix returns prefix of staking records in the tax
//func GetStakingPrefix(taxID uint64) []byte {
//	key := make([]byte, 9)
//	key[0] = StakingKeyPrefix[0]
//	copy(key[1:9], sdk.Uint64ToBigEndian(taxID))
//	return key
//}
//
//// GetStakingIndexKey returns key for farmer's staking of corresponding the tax id
//func GetStakingIndexKey(taxID uint64, farmerAcc sdk.AccAddress) []byte {
//	// TODO: review for addrLen,  <addrLen (1 Byte)><addrBytes>
//	return append(append(StakingKeyPrefix, sdk.Uint64ToBigEndian(taxID)...), address.MustLengthPrefix(farmerAcc.Bytes())...)
//}
//
//// GetRewardPrefix returns prefix of reward records in the tax
//func GetRewardPrefix(taxID uint64) []byte {
//	key := make([]byte, 9)
//	key[0] = RewardKeyPrefix[0]
//	copy(key[1:9], sdk.Uint64ToBigEndian(taxID))
//	return key
//}
//
//// GetRewardIndexKey returns key for farmer's reward of corresponding the tax id
//func GetRewardIndexKey(taxID uint64, farmerAcc sdk.AccAddress) []byte {
//	// TODO: review for addrLen,  <addrLen (1 Byte)><addrBytes>
//	return append(append(RewardKeyPrefix, sdk.Uint64ToBigEndian(taxID)...), address.MustLengthPrefix(farmerAcc.Bytes())...)
//}
