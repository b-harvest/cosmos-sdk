package auth

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// StoreKey is string representation of the store key for auth
	StoreKey = "acc"

	// FeeStoreKey is a string representation of the store key for fees
	FeeStoreKey = "fee"

	// QuerierRoute is the querier route for acc
	QuerierRoute = StoreKey
)

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	globalAccountNumberKey = []byte("globalAccountNumber")

	globalSubKeyNumberKey = []byte("globalSubKeyNumber")

	RouterKeyList = []string{"bank", "crisis", "distr", "gov", "slashing", "staking", "auth"}
)

// AccountKeeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical Account constructor.
	proto func() Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace params.Subspace
}

// NewAccountKeeper returns a new sdk.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore params.Subspace, proto func() Account,
) AccountKeeper {

	return AccountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// NewAccountWithAddress implements sdk.AccountKeeper.
func (ak AccountKeeper) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) Account {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	err = acc.SetAccountNumber(ak.GetNextAccountNumber(ctx))
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return acc
}

// NewAccount creates a new account
func (ak AccountKeeper) NewAccount(ctx sdk.Context, acc Account) Account {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

// AddressStoreKey turn an address to key used to get it from the account store
func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}

// GetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) Account {
	store := ctx.KVStore(ak.key)
	bz := store.Get(AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx sdk.Context) []Account {
	accounts := []Account{}
	appendAccount := func(acc Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	ak.IterateAccounts(ctx, appendAccount)
	return accounts
}

// SetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) SetAccount(ctx sdk.Context, acc Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	bz, err := ak.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	store.Set(AddressStoreKey(addr), bz)
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	store.Delete(AddressStoreKey(addr))
}

// IterateAccounts implements sdk.AccountKeeper.
func (ak AccountKeeper) IterateAccounts(ctx sdk.Context, process func(Account) (stop bool)) {
	store := ctx.KVStore(ak.key)
	iter := sdk.KVStorePrefixIterator(store, AddressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := ak.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// GetPubKey Returns the PubKey of the account at address
func (ak AccountKeeper) GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (crypto.PubKey, sdk.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (ak AccountKeeper) GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, sdk.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetSequence(), nil
}

func (ak AccountKeeper) setSequence(ctx sdk.Context, addr sdk.AccAddress, newSequence uint64) sdk.Error {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	if err := acc.SetSequence(newSequence); err != nil {
		panic(err)
	}

	ak.SetAccount(ctx, acc)
	return nil
}

// GetNextAccountNumber Returns and increments the global account number counter
func (ak AccountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(ak.key)
	bz := store.Get(globalAccountNumberKey)
	if bz == nil {
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(globalAccountNumberKey, bz)

	return accNumber
}

// GetNextSubKeyNumber Returns and increments the global subKey number counter
func (ak AccountKeeper) GetNextSubKeyNumber(ctx sdk.Context) uint64 {
	var subKeyNumber uint64
	store := ctx.KVStore(ak.key)
	bz := store.Get(globalSubKeyNumberKey)
	if bz == nil {
		subKeyNumber = 1  // 0 is reserved for NoSubKey
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &subKeyNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(subKeyNumber + 1)
	store.Set(globalSubKeyNumberKey, bz)

	return subKeyNumber
}

func (ak AccountKeeper) NewSubKey(ctx sdk.Context, pubKey crypto.PubKey, permissionedRoutes []string, dailySpendableAllowance sdk.Coins) SubKey {
	subKey := SubKey{
		PubKey:                  pubKey,
		SubKeyNumber:            ak.GetNextSubKeyNumber(ctx),
		PermissionedRoutes:      permissionedRoutes,
		DailySpendableAllowance: dailySpendableAllowance,
	}
	return subKey
}

// CreateSubKeyAccount
func (ak AccountKeeper) MigrateSubKeyAccount(ctx sdk.Context, acc Account) Account {
	vacc, ok := acc.(VestingAccount)
	if ok {
		panic(sdk.ErrUnknownAddress(fmt.Sprintf("Vesting Account %s can't be SubKey Account", vacc.GetAddress())))
	}

	params := ak.GetParams(ctx)

	dailySpendableAllowance := acc.GetCoins()
	for i := range dailySpendableAllowance {
		dailySpendableAllowance[i].Amount = sdk.NewInt(NoLimitAllowance) // // -1, nolimit dailySpendableAllowance
	}

	masterSubKey := ak.NewSubKey(ctx, acc.GetPubKey(), RouterKeyList, dailySpendableAllowance)

	// Deduct MinDepositAmount
	deposit := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(params.MinDepositAmount)))}
	err := acc.SetCoins(acc.GetCoins().Sub(deposit))
	if err != nil {
		panic(err)
	}

	newSubKeyAcc := &SubKeyAccount{
		BaseAccount:     &BaseAccount{Address: acc.GetAddress()},
		Deposit:         deposit,
	}
	newSubKeyAcc.SubKeys = append(newSubKeyAcc.SubKeys, masterSubKey)
	ak.SetAccount(ctx, acc)

	return acc
}

// CreateSubKeyAccountWithoutPrivateKey
func (ak AccountKeeper) NewSubKeyAccount(ctx sdk.Context, acc Account, pubKey crypto.PubKey) Account {
	newSubKeyAccAddress := sdk.AccAddress(crypto.AddressHash([]byte(fmt.Sprintf("%s %d", acc.GetAddress(), acc.GetSequence()))))
	newBaseAcc := NewBaseAccountWithAddress(newSubKeyAccAddress)
	err := newBaseAcc.SetAccountNumber(ak.GetNextAccountNumber(ctx))
	if err != nil {
		panic(err)
	}
	params := ak.GetParams(ctx)

	dailySpendableAllowance := acc.GetCoins()
	for i := range dailySpendableAllowance {
		dailySpendableAllowance[i].Amount = sdk.NewInt(NoLimitAllowance) // // -1, nolimit dailySpendableAllowance
	}

	masterSubKey := ak.NewSubKey(ctx, pubKey, RouterKeyList, dailySpendableAllowance)

	// Deduct MinDepositAmount
	deposit := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(params.MinDepositAmount)))}
	err = acc.SetCoins(acc.GetCoins().Sub(deposit))
	if err != nil {
		panic(err)
	}

	newSubKeyAcc := &SubKeyAccount{
		BaseAccount:     &newBaseAcc,
		Deposit:         deposit,
	}
	newSubKeyAcc.SubKeys = append(newSubKeyAcc.SubKeys, masterSubKey)
	ak.SetAccount(ctx, acc)

	return newSubKeyAcc
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (ak AccountKeeper) SetParams(ctx sdk.Context, params Params) {
	ak.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (ak AccountKeeper) GetParams(ctx sdk.Context) (params Params) {
	ak.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Misc.

func (ak AccountKeeper) decodeAccount(bz []byte) (acc Account) {
	err := ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
