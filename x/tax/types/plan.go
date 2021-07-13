package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


// NewTax creates a new Tax object
//nolint:interfacer
func NewTax(name, poolAddr, taxSourceAddr string, rate sdk.Dec, startTime, endTime time.Time) *Tax {
	baseTax := &Tax{
		Name:    name,
		PoolAddress:    poolAddr,
		TaxSourceAddress:    taxSourceAddr,
		Rate:    rate,
		StartTime:             startTime,
		EndTime:               endTime,
	}
	return baseTax
}

// Validate checks for errors on the Tax fields
func (tax Tax) Validate() error {
	if _, err := sdk.AccAddressFromBech32(tax.PoolAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid tax pool address %q: %v", tax.PoolAddress, err)
	}
	if _, err := sdk.AccAddressFromBech32(tax.TaxSourceAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid tax source address %q: %v", tax.TaxSourceAddress, err)
	}
	// TODO: rate, name max length
	if !tax.EndTime.After(tax.StartTime) {
		return sdkerrors.Wrapf(ErrInvalidTaxEndTime, "end time %s must be greater than start time %s", tax.EndTime, tax.StartTime)
	}
	return nil
}

func (tax Tax) String() string {
	out, _ := tax.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an Tax.
func (tax Tax) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &tax)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}
