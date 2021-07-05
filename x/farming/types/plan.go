package types

import (
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ PlanI = (*FixedAmountPlan)(nil)
	_ PlanI = (*RatioPlan)(nil)
)

// NewBasePlan creates a new BasePlan object
//nolint:interfacer
func NewBasePlan(id uint64, typ PlanType, farmingPoolAddr, distPoolAddr, terminationAddr, reserveAddr string, coinWeights sdk.DecCoins, startTime, endTime time.Time, epochDays uint32) *BasePlan {
	basePlan := &BasePlan{
		Id:                      id,
		Type:                    typ,
		FarmingPoolAddress:      farmingPoolAddr,
		DistributionPoolAddress: distPoolAddr,
		TerminationAddress:      terminationAddr,
		StakingReserveAddress:   reserveAddr,
		StakingCoinWeights:      coinWeights,
		StartTime:               startTime,
		EndTime:                 endTime,
		EpochDays:               epochDays,
	}
	return basePlan
}

func (plan BasePlan) GetId() uint64 {
	return plan.Id
}

func (plan *BasePlan) SetId(id uint64) error {
	plan.Id = id
	return nil
}

func (plan BasePlan) GetType() PlanType {
	return plan.Type
}

func (plan *BasePlan) SetType(typ PlanType) error {
	plan.Type = typ
	return nil
}

func (plan BasePlan) GetFarmingPoolAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(plan.FarmingPoolAddress)
	return addr
}

func (plan *BasePlan) SetFarmingPoolAddress(addr sdk.AccAddress) error {
	plan.FarmingPoolAddress = addr.String()
	return nil
}

func (plan BasePlan) GetDistributionPoolAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(plan.DistributionPoolAddress)
	return addr
}

func (plan *BasePlan) SetDistributionPoolAddress(addr sdk.AccAddress) error {
	plan.DistributionPoolAddress = addr.String()
	return nil
}

func (plan BasePlan) GetTerminationAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(plan.TerminationAddress)
	return addr
}

func (plan *BasePlan) SetTerminationAddress(addr sdk.AccAddress) error {
	plan.TerminationAddress = addr.String()
	return nil
}

func (plan BasePlan) GetStakingReserveAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(plan.StakingReserveAddress)
	return addr
}

func (plan *BasePlan) SetStakingReserveAddress(addr sdk.AccAddress) error {
	plan.StakingReserveAddress = addr.String()
	return nil
}

func (plan BasePlan) GetStakingCoinWeights() sdk.DecCoins {
	return plan.StakingCoinWeights
}

func (plan *BasePlan) SetStakingCoinWeights(coinWeights sdk.DecCoins) error {
	plan.StakingCoinWeights = coinWeights
	return nil
}

func (plan BasePlan) GetStartTime() time.Time {
	return plan.StartTime
}

func (plan *BasePlan) SetStartTime(t time.Time) error {
	plan.StartTime = t
	return nil
}

func (plan BasePlan) GetEndTime() time.Time {
	return plan.EndTime
}

func (plan *BasePlan) SetEndTime(t time.Time) error {
	plan.EndTime = t
	return nil
}

func (plan BasePlan) GetEpochDays() uint32 {
	return plan.EpochDays
}

func (plan *BasePlan) SetEpochDays(days uint32) error {
	plan.EpochDays = days
	return nil
}

// Validate checks for errors on the Plan fields
func (plan BasePlan) Validate() error {
	// TODO: Unimplemented
	return nil
}

func (plan BasePlan) String() string {
	out, _ := plan.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an Plan.
func (plan BasePlan) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &plan)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

func NewFixedAmountPlan(basePlan *BasePlan, epochAmount sdk.Coins) *FixedAmountPlan {
	return &FixedAmountPlan{
		BasePlan:    basePlan,
		EpochAmount: epochAmount,
	}
}

func NewRatioPlan(basePlan *BasePlan, epochRatio sdk.Dec) *RatioPlan {
	return &RatioPlan{
		BasePlan:   basePlan,
		EpochRatio: epochRatio,
	}
}

type PlanI interface {
	proto.Message

	GetId() uint64
	SetId(uint64) error

	GetType() PlanType
	SetType(PlanType) error

	GetFarmingPoolAddress() sdk.AccAddress
	SetFarmingPoolAddress(sdk.AccAddress) error

	GetDistributionPoolAddress() sdk.AccAddress
	SetDistributionPoolAddress(sdk.AccAddress) error

	GetTerminationAddress() sdk.AccAddress
	SetTerminationAddress(sdk.AccAddress) error

	GetStakingReserveAddress() sdk.AccAddress
	SetStakingReserveAddress(sdk.AccAddress) error

	GetStakingCoinWeights() sdk.DecCoins
	SetStakingCoinWeights(sdk.DecCoins) error

	GetStartTime() time.Time
	SetStartTime(time.Time) error

	GetEndTime() time.Time
	SetEndTime(time.Time) error

	GetEpochDays() uint32
	SetEpochDays(uint32) error

	String() string
}

func (s Staking) String() string {
	// TODO: not implemented
	return ""
}

func (r Reward) String() string {
	// TODO: not implemented
	return ""
}
