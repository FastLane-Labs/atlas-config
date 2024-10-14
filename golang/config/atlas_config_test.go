package config

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetSupportedChainIds(t *testing.T) {
	chainIds := GetSupportedChainIds()
	assert.NotEmpty(t, chainIds)
	assert.Contains(t, chainIds, uint64(11155111)) // Assuming Sepolia testnet is always supported
}

func TestGetAllChainConfigs(t *testing.T) {
	configs := GetAllChainConfigs()
	assert.NotEmpty(t, configs)
	assert.Contains(t, configs, "11155111") // Assuming Sepolia testnet is always supported

	// Check structure of a config
	sepoliaConfig, ok := configs["11155111"]
	require.True(t, ok)
	assert.NotEmpty(t, sepoliaConfig.Contracts.Atlas)
	assert.NotEmpty(t, sepoliaConfig.EIP712Domain.Name)
}

func TestMergeChainConfigs(t *testing.T) {
	// Test merging with an existing chain
	existingChainId := "11155111"
	newAtlasAddress := "0x1234567890123456789012345678901234567890"
	providedConfigs := map[string]ChainConfig{
		existingChainId: {
			Contracts: Contracts{
				Atlas: common.HexToAddress(newAtlasAddress),
			},
		},
	}

	err := MergeChainConfigs(providedConfigs)
	assert.NoError(t, err)

	updatedConfig, err := GetChainConfig(11155111)
	assert.NoError(t, err)
	assert.Equal(t, newAtlasAddress, updatedConfig.Contracts.Atlas.Hex())

	// Test adding a new chain
	newChainId := "999999"
	newChainConfig := ChainConfig{
		Contracts: Contracts{
			Atlas:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
			AtlasVerification: common.HexToAddress("0x2222222222222222222222222222222222222222"),
			Sorter:            common.HexToAddress("0x3333333333333333333333333333333333333333"),
			Simulator:         common.HexToAddress("0x4444444444444444444444444444444444444444"),
			Multicall3:        common.HexToAddress("0x5555555555555555555555555555555555555555"),
		},
		EIP712Domain: EIP712Domain{
			Name:              "NewChain",
			Version:           "1.0",
			ChainId:           999999,
			VerifyingContract: "0x6666666666666666666666666666666666666666",
		},
	}
	providedConfigs = map[string]ChainConfig{
		newChainId: newChainConfig,
	}

	err = MergeChainConfigs(providedConfigs)
	assert.NoError(t, err)

	addedConfig, err := GetChainConfig(999999)
	assert.NoError(t, err)
	assert.Equal(t, newChainConfig.Contracts.Atlas, addedConfig.Contracts.Atlas)
	assert.Equal(t, newChainConfig.EIP712Domain.Name, addedConfig.EIP712Domain.Name)

	// Test adding an incomplete new chain config (should fail)
	incompleteChainId := "888888"
	incompleteConfig := ChainConfig{
		Contracts: Contracts{
			Atlas: common.HexToAddress("0x7777777777777777777777777777777777777777"),
		},
	}
	providedConfigs = map[string]ChainConfig{
		incompleteChainId: incompleteConfig,
	}

	err = MergeChainConfigs(providedConfigs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "full chain configuration must be provided for new chainId: 888888")
}

func TestThreadSafety(t *testing.T) {
	// This test is not comprehensive, but it can help catch obvious race conditions
	done := make(chan bool)
	go func() {
		_, _ = GetChainConfig(11155111)
		done <- true
	}()
	go func() {
		_ = GetSupportedChainIds()
		done <- true
	}()
	go func() {
		_ = GetAllChainConfigs()
		done <- true
	}()
	go func() {
		_ = MergeChainConfigs(map[string]ChainConfig{})
		done <- true
	}()

	for i := 0; i < 4; i++ {
		<-done
	}
}
