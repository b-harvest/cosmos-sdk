package types

import (
	"fmt"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeFixedAmountPlan string = "FixedAmountPlan"
	ProposalTypeRatioPlan       string = "RatioPlan"
)

// Implements Proposal Interface
var _ gov.Content = &FixedAmountPlanProposal{}
var _ gov.Content = &RatioPlanProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeFixedAmountPlan)
	gov.RegisterProposalTypeCodec(&FixedAmountPlanProposal{}, "cosmos-sdk/FixedAmountPlanProposal")
	gov.RegisterProposalType(ProposalTypeRatioPlan)
	gov.RegisterProposalTypeCodec(&RatioPlanProposal{}, "cosmos-sdk/RatioPlanProposal")
}

func NewFixedAmountPlanProposal(title, description string, plan FixedAmountPlan) gov.Content {
	return &FixedAmountPlanProposal{title, description, plan}
}

func (p *FixedAmountPlanProposal) GetTitle() string { return p.Title }

func (p *FixedAmountPlanProposal) GetDescription() string { return p.Description }

func (p *FixedAmountPlanProposal) ProposalRoute() string { return RouterKey }

func (p *FixedAmountPlanProposal) ProposalType() string { return ProposalTypeFixedAmountPlan }

func (p *FixedAmountPlanProposal) ValidateBasic() error {
	// TODO: more checks? (ex. length check?)
	if err := p.Plan.Validate(); err != nil {
		return err
	}
	return gov.ValidateAbstract(p)
}

func (p FixedAmountPlanProposal) String() string {
	return fmt.Sprintf(`Create FixedAmountPlan Proposal:
  Title:       %s
  Description: %s
  Plan: 	   %s
`, p.Title, p.Description, p.Plan)
}

func NewRatioPlanProposal(title, description string, plan RatioPlan) gov.Content {
	return &RatioPlanProposal{title, description, plan}
}

func (p *RatioPlanProposal) GetTitle() string { return p.Title }

func (p *RatioPlanProposal) GetDescription() string { return p.Description }

func (p *RatioPlanProposal) ProposalRoute() string { return RouterKey }

func (p *RatioPlanProposal) ProposalType() string { return ProposalTypeRatioPlan }

func (p *RatioPlanProposal) ValidateBasic() error {
	// TODO: more checks? (ex. length check?)
	if err := p.Plan.Validate(); err != nil {
		return err
	}
	return gov.ValidateAbstract(p)
}

func (p RatioPlanProposal) String() string {
	return fmt.Sprintf(`Create RatioPlan Proposal:
  Title:       %s
  Description: %s
  Plan:        %s
`, p.Title, p.Description, p.Plan)
}
