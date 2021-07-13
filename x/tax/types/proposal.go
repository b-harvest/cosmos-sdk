package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"

	proto "github.com/gogo/protobuf/proto"
)

const (
	ProposalTypePublicTax string = "PublicTax"
)

// Implements Proposal Interface
var _ gov.Content = &SetTaxesProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypePublicTax)
	gov.RegisterProposalTypeCodec(&SetTaxesProposal{}, "cosmos-sdk/SetTaxesProposal")
}

func NewSetTaxesProposal(title, description string, taxes []Tax) (gov.Content, error) {
	taxesAny, err := PackTaxes(taxes)
	if err != nil {
		panic(err)
	}

	return &SetTaxesProposal{
		Title:       title,
		Description: description,
		Taxes:       taxesAny,
	}, nil
}

func (p *SetTaxesProposal) GetTitle() string { return p.Title }

func (p *SetTaxesProposal) GetDescription() string { return p.Description }

func (p *SetTaxesProposal) ProposalRoute() string { return RouterKey }

func (p *SetTaxesProposal) ProposalType() string { return ProposalTypePublicTax }

func (p *SetTaxesProposal) ValidateBasic() error {
	for _, tax := range p.Taxes {
		_, ok := tax.GetCachedValue().(TaxI)
		if !ok {
			return fmt.Errorf("expected taxI")
		}
		// TODO: TaxI needs ValidateBasic()?
		// if err := p.ValidateBasic(); err != nil {
		// 	return err
		// }
	}
	return gov.ValidateAbstract(p)
}

func (p SetTaxesProposal) String() string {
	return fmt.Sprintf(`Create FixedAmountTax Proposal:
  Title:       %s
  Description: %s
  Taxes: 	   %s
`, p.Title, p.Description, p.Taxes)
}

// PackTaxes converts TaxIs to Any slice.
func PackTaxes(taxes []TaxI) ([]*types.Any, error) {
	taxesAny := make([]*types.Any, len(taxes))
	for i, tax := range taxes {
		msg, ok := tax.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("cannot proto marshal %T", tax)
		}
		any, err := types.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
		taxesAny[i] = any
	}

	return taxesAny, nil
}

// UnpackTaxes converts Any slice to TaxIs.
func UnpackTaxes(taxesAny []*types.Any) ([]TaxI, error) {
	taxes := make([]TaxI, len(taxesAny))
	for i, any := range taxesAny {
		p, ok := any.GetCachedValue().(TaxI)
		if !ok {
			return nil, fmt.Errorf("expected taxI")
		}
		taxes[i] = p
	}

	return taxes, nil
}

// UnpackTax converts Any slice to TaxI.
func UnpackTax(any *types.Any) (TaxI, error) {
	p, ok := any.GetCachedValue().(TaxI)
	if !ok {
		return nil, fmt.Errorf("expected taxI")
	}

	return p, nil
}
