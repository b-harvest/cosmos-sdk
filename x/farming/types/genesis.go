package types

// NewGenesisState returns new GenesisState.
func NewGenesisState(params Params, planRecords []PlanRecord, stakings []Staking, rewards []Reward) *GenesisState {
	return &GenesisState{
		Params:      params,
		PlanRecords: planRecords,
		Stakings:    stakings,
		Rewards:     rewards,
	}
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []PlanRecord{}, []Staking{}, []Reward{})
}

// ValidateGenesis validates GenesisState.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	for _, record := range data.PlanRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}
	// TODO: unimplemented
	//for _, staking := range data.Stakings {
	//	if err := staking.Validate(); err != nil {
	//		return err
	//	}
	//}
	//for _, reward := range data.Rewards {
	//	if err := reward.Validate(); err != nil {
	//		return err
	//	}
	//}
	return nil
}

// Validate validates PlanRecord.
func (record PlanRecord) Validate() error {
	// TODO: unimplemented
	return nil
}
