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
	config, err := GetChainConfig(11155111, &AtlasV100)
	if assert.NoError(t, err) {
		assert.NotNil(t, config)
		assert.Equal(t, "0x9EE12d2fed4B43F4Be37F69930CcaD9B65133482", config.Contracts.Atlas.Hex())
	}

	// Test for a non-existing chain ID
	_, err = GetChainConfig(999999, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chain configuration not found")
}

func TestGetAtlasAddress(t *testing.T) {
	address, err := GetAtlasAddress(11155111, &AtlasV100)
	assert.NoError(t, err)
	assert.Equal(t, "0x9EE12d2fed4B43F4Be37F69930CcaD9B65133482", address.Hex())

	_, err = GetAtlasAddress(999999, nil)
	assert.Error(t, err)
}

func TestGetAtlasVerificationAddress(t *testing.T) {
	address, err := GetAtlasVerificationAddress(11155111, &AtlasV100)
	assert.NoError(t, err)
	assert.Equal(t, "0xB6F66a1b7cec02324D83c8DEA192818cA23A08B3", address.Hex())

	_, err = GetAtlasVerificationAddress(999999, nil)
	assert.Error(t, err)
}

func TestGetSorterAddress(t *testing.T) {
	address, err := GetSorterAddress(11155111, &AtlasV100)
	assert.NoError(t, err)
	assert.Equal(t, "0xFE3c655d4D305Ac7f1c2F6306C79397560Afea0C", address.Hex())

	_, err = GetSorterAddress(999999, nil)
	assert.Error(t, err)
}

func TestGetSimulatorAddress(t *testing.T) {
	address, err := GetSimulatorAddress(11155111, &AtlasV100)
	assert.NoError(t, err)
	assert.Equal(t, "0xc3ab39ebd49D80bc36208545021224BAF6d2Bdb0", address.Hex())

	_, err = GetSimulatorAddress(999999, nil)
	assert.Error(t, err)
}

func TestGetMulticall3Address(t *testing.T) {
	address, err := GetMulticall3Address(11155111, &AtlasV100)
	assert.NoError(t, err)
	assert.Equal(t, "0xcA11bde05977b3631167028862bE2a173976CA11", address.Hex())

	_, err = GetMulticall3Address(999999, nil)
	assert.Error(t, err)
}

func TestGetEIP712Domain(t *testing.T) {
	domain, err := GetEIP712Domain(11155111, &AtlasV100)
	assert.NoError(t, err)
	assert.Equal(t, "AtlasVerification", domain.Name)
	assert.Equal(t, "1.0", domain.Version)
	assert.Equal(t, big.NewInt(11155111), (*big.Int)(domain.ChainId))
	assert.Equal(t, "0xB6F66a1b7cec02324D83c8DEA192818cA23A08B3", domain.VerifyingContract)

	_, err = GetEIP712Domain(999999, nil)
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
	assert.NotEmpty(t, sepoliaConfig[AtlasV100].Contracts.Atlas)
	assert.NotEmpty(t, sepoliaConfig[AtlasV100].EIP712Domain.Name)
}

func TestMergeChainConfigs(t *testing.T) {
	// Setup initial config
	chainConfig = map[string]map[string]ChainConfig{
		"1": {
			"0.0.1": {
				Contracts: Contracts{
					Atlas:             common.HexToAddress("0x1000000000000000000000000000000000000000"),
					AtlasVerification: common.HexToAddress("0x2000000000000000000000000000000000000000"),
					Sorter:            common.HexToAddress("0x3000000000000000000000000000000000000000"),
					Simulator:         common.HexToAddress("0x4000000000000000000000000000000000000000"),
					Multicall3:        common.HexToAddress("0x5000000000000000000000000000000000000000"),
				},
				EIP712Domain: EIP712Domain{
					Name:              "Initial",
					Version:           "1",
					ChainId:           1,
					VerifyingContract: "0x6000000000000000000000000000000000000000",
				},
			},
		},
	}

	t.Run("Merge partial config for existing chain", func(t *testing.T) {
		err := MergeChainConfigs(map[string]map[string]ChainConfig{
			"1": {
				"0.0.1": {
					Contracts: Contracts{
						AtlasVerification: common.HexToAddress("0x7000000000000000000000000000000000000000"),
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, common.HexToAddress("0x1000000000000000000000000000000000000000"), chainConfig["1"]["0.0.1"].Contracts.Atlas)
		assert.Equal(t, common.HexToAddress("0x7000000000000000000000000000000000000000"), chainConfig["1"]["0.0.1"].Contracts.AtlasVerification)
	})

	t.Run("Prioritize new complete config for existing chain", func(t *testing.T) {
		newCompleteConfig := ChainConfig{
			Contracts: Contracts{
				Atlas:             common.HexToAddress("0x8000000000000000000000000000000000000000"),
				AtlasVerification: common.HexToAddress("0x9000000000000000000000000000000000000000"),
				Sorter:            common.HexToAddress("0xa000000000000000000000000000000000000000"),
				Simulator:         common.HexToAddress("0xb000000000000000000000000000000000000000"),
				Multicall3:        common.HexToAddress("0xc000000000000000000000000000000000000000"),
			},
			EIP712Domain: EIP712Domain{
				Name:              "New",
				Version:           "2",
				ChainId:           1,
				VerifyingContract: "0xd000000000000000000000000000000000000000",
			},
		}
		err := MergeChainConfigs(map[string]map[string]ChainConfig{"1": {"0.0.1": newCompleteConfig}})
		assert.NoError(t, err)
		assert.Equal(t, newCompleteConfig, chainConfig["1"]["0.0.1"])
	})

	t.Run("Add new chain config", func(t *testing.T) {
		newChainConfig := ChainConfig{
			Contracts: Contracts{
				Atlas:             common.HexToAddress("0xe000000000000000000000000000000000000000"),
				AtlasVerification: common.HexToAddress("0xf000000000000000000000000000000000000000"),
				Sorter:            common.HexToAddress("0x1000000000000000000000000000000000000001"),
				Simulator:         common.HexToAddress("0x2000000000000000000000000000000000000001"),
				Multicall3:        common.HexToAddress("0x3000000000000000000000000000000000000001"),
			},
			EIP712Domain: EIP712Domain{
				Name:              "New Chain",
				Version:           "1",
				ChainId:           2,
				VerifyingContract: "0x4000000000000000000000000000000000000001",
			},
		}
		err := MergeChainConfigs(map[string]map[string]ChainConfig{"2": {"0.0.1": newChainConfig}})
		assert.NoError(t, err)
		assert.Equal(t, newChainConfig, chainConfig["2"]["0.0.1"])
	})

	t.Run("Attempt to add incomplete new chain config", func(t *testing.T) {
		err := MergeChainConfigs(map[string]map[string]ChainConfig{
			"3": {
				"0.0.1": {
					Contracts: Contracts{
						Atlas: common.HexToAddress("0x5000000000000000000000000000000000000001"),
					},
				},
			},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "full chain configuration must be provided for new chainId: 3")
	})

	t.Run("Merge multiple chains at once", func(t *testing.T) {
		existingChainId := "1"
		newChainId := "4"
		providedConfigs := map[string]map[string]ChainConfig{
			existingChainId: {
				"0.0.1": {
					Contracts: Contracts{
						Atlas: common.HexToAddress("0x6000000000000000000000000000000000000001"),
					},
				},
			},
			newChainId: {
				"0.0.1": {
					Contracts: Contracts{
						Atlas:             common.HexToAddress("0x7000000000000000000000000000000000000001"),
						AtlasVerification: common.HexToAddress("0x8000000000000000000000000000000000000001"),
						Sorter:            common.HexToAddress("0x9000000000000000000000000000000000000001"),
						Simulator:         common.HexToAddress("0xa000000000000000000000000000000000000001"),
						Multicall3:        common.HexToAddress("0xb000000000000000000000000000000000000001"),
					},
					EIP712Domain: EIP712Domain{
						Name:              "Another New Chain",
						Version:           "1",
						ChainId:           4,
						VerifyingContract: "0xc000000000000000000000000000000000000001",
					},
				},
			},
		}

		err := MergeChainConfigs(providedConfigs)
		assert.NoError(t, err)
		assert.Equal(t, common.HexToAddress("0x6000000000000000000000000000000000000001"), chainConfig[existingChainId]["0.0.1"].Contracts.Atlas)
		assert.Equal(t, common.HexToAddress("0x7000000000000000000000000000000000000001"), chainConfig[newChainId]["0.0.1"].Contracts.Atlas)
	})

	t.Run("Merge empty config", func(t *testing.T) {
		originalConfig := make(map[string]map[string]ChainConfig)
		for k, v := range chainConfig {
			originalConfig[k] = make(map[string]ChainConfig)
			for k2, v2 := range v {
				originalConfig[k][k2] = v2
			}
		}
		err := MergeChainConfigs(map[string]map[string]ChainConfig{})
		assert.NoError(t, err)
		assert.Equal(t, originalConfig, chainConfig)
	})

	t.Run("Merge config with no changes", func(t *testing.T) {
		originalConfig := make(map[string]map[string]ChainConfig)
		for k, v := range chainConfig {
			originalConfig[k] = make(map[string]ChainConfig)
			for k2, v2 := range v {
				originalConfig[k][k2] = v2
			}
		}
		err := MergeChainConfigs(originalConfig)
		assert.NoError(t, err)
		assert.Equal(t, originalConfig, chainConfig)
	})

	t.Run("Merge partial config that updates multiple fields", func(t *testing.T) {
		existingChainId := "1"
		providedConfig := map[string]map[string]ChainConfig{
			existingChainId: {
				"0.0.1": {
					Contracts: Contracts{
						Atlas:  common.HexToAddress("0xd000000000000000000000000000000000000001"),
						Sorter: common.HexToAddress("0xe000000000000000000000000000000000000001"),
					},
					EIP712Domain: EIP712Domain{
						Name:    "Updated Chain",
						Version: "2.0",
					},
				},
			},
		}

		err := MergeChainConfigs(providedConfig)
		assert.NoError(t, err)
		assert.Equal(t, common.HexToAddress("0xd000000000000000000000000000000000000001"), chainConfig[existingChainId]["0.0.1"].Contracts.Atlas)
		assert.Equal(t, common.HexToAddress("0xe000000000000000000000000000000000000001"), chainConfig[existingChainId]["0.0.1"].Contracts.Sorter)
		assert.Equal(t, "Updated Chain", chainConfig[existingChainId]["0.0.1"].EIP712Domain.Name)
		assert.Equal(t, "2.0", chainConfig[existingChainId]["0.0.1"].EIP712Domain.Version)
	})
}

func TestThreadSafety(t *testing.T) {
	// This test is not comprehensive, but it can help catch obvious race conditions
	done := make(chan bool)
	go func() {
		_, _ = GetChainConfig(11155111, &AtlasV100)
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
		_ = MergeChainConfigs(map[string]map[string]ChainConfig{})
		done <- true
	}()

	for i := 0; i < 4; i++ {
		<-done
	}
}
