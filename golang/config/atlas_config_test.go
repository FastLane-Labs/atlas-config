package config

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChainConfig(t *testing.T) {
	// Test for an existing chain ID
	config, err := GetChainConfig(11155111)
	if assert.NoError(t, err) {
		assert.NotNil(t, config)
		assert.Equal(t, "0x9EE12d2fed4B43F4Be37F69930CcaD9B65133482", config.Contracts.Atlas.Hex())
	}

	// Test for a non-existing chain ID
	_, err = GetChainConfig(999999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chain configuration not found")
}

func TestGetAtlasAddress(t *testing.T) {
	address, err := GetAtlasAddress(11155111)
	assert.NoError(t, err)
	assert.Equal(t, "0x9EE12d2fed4B43F4Be37F69930CcaD9B65133482", address.Hex())

	_, err = GetAtlasAddress(999999)
	assert.Error(t, err)
}

func TestGetAtlasVerificationAddress(t *testing.T) {
	address, err := GetAtlasVerificationAddress(11155111)
	assert.NoError(t, err)
	assert.Equal(t, "0xB6F66a1b7cec02324D83c8DEA192818cA23A08B3", address.Hex())

	_, err = GetAtlasVerificationAddress(999999)
	assert.Error(t, err)
}

func TestGetSorterAddress(t *testing.T) {
	address, err := GetSorterAddress(11155111)
	assert.NoError(t, err)
	assert.Equal(t, "0xFE3c655d4D305Ac7f1c2F6306C79397560Afea0C", address.Hex())

	_, err = GetSorterAddress(999999)
	assert.Error(t, err)
}

func TestGetSimulatorAddress(t *testing.T) {
	address, err := GetSimulatorAddress(11155111)
	assert.NoError(t, err)
	assert.Equal(t, "0xc3ab39ebd49D80bc36208545021224BAF6d2Bdb0", address.Hex())

	_, err = GetSimulatorAddress(999999)
	assert.Error(t, err)
}

func TestGetMulticall3Address(t *testing.T) {
	address, err := GetMulticall3Address(11155111)
	assert.NoError(t, err)
	assert.Equal(t, "0xcA11bde05977b3631167028862bE2a173976CA11", address.Hex())

	_, err = GetMulticall3Address(999999)
	assert.Error(t, err)
}

func TestGetEIP712Domain(t *testing.T) {
	domain, err := GetEIP712Domain(11155111)
	assert.NoError(t, err)
	assert.Equal(t, "AtlasVerification", domain.Name)
	assert.Equal(t, "1.0", domain.Version)
	assert.Equal(t, big.NewInt(11155111), (*big.Int)(domain.ChainId))
	assert.Equal(t, "0xB6F66a1b7cec02324D83c8DEA192818cA23A08B3", domain.VerifyingContract)

	_, err = GetEIP712Domain(999999)
	assert.Error(t, err)
}
