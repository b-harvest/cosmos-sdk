package testepoching

import (
	"cosmossdk.io/math"
	tmcrypto "github.com/cometbft/cometbft/crypto"
	tmtypes "github.com/cometbft/cometbft/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/x/epoching/testepoching/datagen"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetTmConsPubKey gets the validator's public key as a tmcrypto.PubKey.
func GetTmConsPubKey(v stakingtypes.Validator) (tmcrypto.PubKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}

	return cryptocodec.ToTmPubKeyInterface(pk)
}

// ToTmValidator casts an SDK validator to a tendermint type Validator.
func ToTmValidator(v stakingtypes.Validator, r math.Int) (*tmtypes.Validator, error) {
	tmPk, err := GetTmConsPubKey(v)
	if err != nil {
		return nil, err
	}

	return tmtypes.NewValidator(tmPk, v.ConsensusPower(r)), nil
}

//// ToTmValidators casts all validators to the corresponding tendermint type.
//func ToTmValidators(v stakingtypes.Validators, r math.Int) ([]*tmtypes.Validator, error) {
//	validators := make([]*tmtypes.Validator, len(v))
//	var err error
//	for i, val := range v {
//		validators[i], err = ToTmValidator(val, r)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return validators, nil
//}

// GenTmValidatorSet generates a set with `numVals` Tendermint validators
func GenTmValidatorSet(numVals int) (*tmtypes.ValidatorSet, error) {
	vals := []*tmtypes.Validator{}
	for i := 0; i < numVals; i++ {
		privVal := datagen.NewPV()
		pubKey, err := privVal.GetPubKey()
		if err != nil {
			return nil, err
		}
		val := tmtypes.NewValidator(pubKey, 1)
		vals = append(vals, val)
	}
	return tmtypes.NewValidatorSet(vals), nil
}
