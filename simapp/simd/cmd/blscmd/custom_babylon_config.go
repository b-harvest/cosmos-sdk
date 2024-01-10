package cmd

import (
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

const (
	defaultKeyName       = ""
	defaultGasPrice      = "0.01ubbn"
	defaultGasAdjustment = 1.5
)

func defaultSignerConfig() SignerConfig {
	return SignerConfig{
		KeyName:       defaultKeyName,
		GasPrice:      defaultGasPrice,
		GasAdjustment: defaultGasAdjustment,
	}
}

type SignerConfig struct {
	KeyName       string  `mapstructure:"key-name"`
	GasPrice      string  `mapstructure:"gas-price"`
	GasAdjustment float64 `mapstructure:"gas-adjustment"`
}

type BabylonAppConfig struct {
	serverconfig.Config `mapstructure:",squash"`

	SignerConfig SignerConfig `mapstructure:"signer-config"`
}

func DefaultBabylonConfig() *BabylonAppConfig {
	return &BabylonAppConfig{
		Config:       *serverconfig.DefaultConfig(),
		SignerConfig: defaultSignerConfig(),
	}
}

func DefaultBabylonTemplate() string {
	return serverconfig.DefaultConfigTemplate + `
###############################################################################
###                      Babylon BLS configuration                      ###
###############################################################################

[signer-config]

# Configures which key that the BLS signer uses to sign BLS-sig transactions
key-name = "{{ .SignerConfig.KeyName }}"
# Configures the gas-price that the signer would like to pay
gas-price = "{{ .SignerConfig.GasPrice }}"
# Configures the adjustment of the gas cost of estimation
gas-adjustment = "{{ .SignerConfig.GasAdjustment }}"
`
}
