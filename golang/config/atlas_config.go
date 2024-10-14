package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type ContractJSON struct {
	Address string `json:"address"`
}

type ContractsJSON struct {
	Atlas             ContractJSON `json:"atlas"`
	AtlasVerification ContractJSON `json:"atlasVerification"`
	Sorter            ContractJSON `json:"sorter"`
	Simulator         ContractJSON `json:"simulator"`
	Multicall3        ContractJSON `json:"multicall3"`
}

type Contracts struct {
	Atlas             common.Address
	AtlasVerification common.Address
	Sorter            common.Address
	Simulator         common.Address
	Multicall3        common.Address
}

type EIP712Domain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainId           uint64 `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

type ChainConfigJSON struct {
	Contracts    ContractsJSON `json:"contracts"`
	EIP712Domain EIP712Domain  `json:"eip712Domain"`
}

type ChainConfig struct {
	Contracts    Contracts
	EIP712Domain EIP712Domain
}

var (
	chainConfig map[string]ChainConfig
	configMutex sync.RWMutex
)

func init() {
	// Load the config file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current file path")
	}
	packageDir := filepath.Dir(currentFile)
	configPath := filepath.Join(packageDir, "chain-config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Config file not found: %s", configPath))
	}

	byteValue, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	var chainConfigJSON map[string]ChainConfigJSON
	err = json.Unmarshal(byteValue, &chainConfigJSON)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling config: %v", err))
	}

	chainConfig = make(map[string]ChainConfig)
	for chainID, config := range chainConfigJSON {
		chainConfig[chainID] = ChainConfig{
			Contracts: Contracts{
				Atlas:             common.HexToAddress(config.Contracts.Atlas.Address),
				AtlasVerification: common.HexToAddress(config.Contracts.AtlasVerification.Address),
				Sorter:            common.HexToAddress(config.Contracts.Sorter.Address),
				Simulator:         common.HexToAddress(config.Contracts.Simulator.Address),
				Multicall3:        common.HexToAddress(config.Contracts.Multicall3.Address),
			},
			EIP712Domain: config.EIP712Domain,
		}
	}
}

// GetChainConfig returns the chain configuration for the given chainId
func GetChainConfig(chainId uint64) (*ChainConfig, error) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	config, ok := chainConfig[strconv.FormatUint(chainId, 10)]
	if !ok {
		return nil, fmt.Errorf("chain configuration not found for chainId: %d", chainId)
	}
	return &config, nil
}

// GetSupportedChainIds returns all supported chain IDs
func GetSupportedChainIds() []uint64 {
	configMutex.RLock()
	defer configMutex.RUnlock()

	chainIds := make([]uint64, 0, len(chainConfig))
	for chainIdStr := range chainConfig {
		chainId, _ := strconv.ParseUint(chainIdStr, 10, 64)
		chainIds = append(chainIds, chainId)
	}
	return chainIds
}

// GetAllChainConfigs returns all chain configurations
func GetAllChainConfigs() map[string]ChainConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Create a deep copy to avoid potential race conditions
	configCopy := make(map[string]ChainConfig, len(chainConfig))
	for k, v := range chainConfig {
		configCopy[k] = v
	}
	return configCopy
}

// MergeChainConfigs merges the provided chain configs with the original
func MergeChainConfigs(providedConfigs map[string]ChainConfig) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	for chainId, providedConfig := range providedConfigs {
		if existingConfig, ok := chainConfig[chainId]; ok {
			// Existing chain: merge config
			mergedConfig := mergeConfigs(existingConfig, providedConfig)
			chainConfig[chainId] = mergedConfig
		} else {
			// New chain: ensure full config is provided
			if !isFullChainConfig(providedConfig) {
				return fmt.Errorf("full chain configuration must be provided for new chainId: %s", chainId)
			}
			chainConfig[chainId] = providedConfig
		}
	}
	return nil
}

// mergeConfigs merges two ChainConfig structs
func mergeConfigs(existing, provided ChainConfig) ChainConfig {
	merged := existing

	// Merge Contracts
	if provided.Contracts.Atlas != (common.Address{}) {
		merged.Contracts.Atlas = provided.Contracts.Atlas
	}
	if provided.Contracts.AtlasVerification != (common.Address{}) {
		merged.Contracts.AtlasVerification = provided.Contracts.AtlasVerification
	}
	if provided.Contracts.Sorter != (common.Address{}) {
		merged.Contracts.Sorter = provided.Contracts.Sorter
	}
	if provided.Contracts.Simulator != (common.Address{}) {
		merged.Contracts.Simulator = provided.Contracts.Simulator
	}
	if provided.Contracts.Multicall3 != (common.Address{}) {
		merged.Contracts.Multicall3 = provided.Contracts.Multicall3
	}

	// Merge EIP712Domain
	if provided.EIP712Domain.Name != "" {
		merged.EIP712Domain.Name = provided.EIP712Domain.Name
	}
	if provided.EIP712Domain.Version != "" {
		merged.EIP712Domain.Version = provided.EIP712Domain.Version
	}
	if provided.EIP712Domain.ChainId != 0 {
		merged.EIP712Domain.ChainId = provided.EIP712Domain.ChainId
	}
	if provided.EIP712Domain.VerifyingContract != "" {
		merged.EIP712Domain.VerifyingContract = provided.EIP712Domain.VerifyingContract
	}

	return merged
}

// isFullChainConfig checks if the provided config is a full ChainConfig
func isFullChainConfig(config ChainConfig) bool {
	return config.Contracts.Atlas != (common.Address{}) &&
		config.Contracts.AtlasVerification != (common.Address{}) &&
		config.Contracts.Sorter != (common.Address{}) &&
		config.Contracts.Simulator != (common.Address{}) &&
		config.Contracts.Multicall3 != (common.Address{}) &&
		config.EIP712Domain.Name != "" &&
		config.EIP712Domain.Version != "" &&
		config.EIP712Domain.ChainId != 0 &&
		config.EIP712Domain.VerifyingContract != ""
}

func GetAtlasAddress(chainId uint64) (common.Address, error) {
	config, err := GetChainConfig(chainId)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Atlas, nil
}

func GetAtlasVerificationAddress(chainId uint64) (common.Address, error) {
	config, err := GetChainConfig(chainId)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.AtlasVerification, nil
}

func GetSorterAddress(chainId uint64) (common.Address, error) {
	config, err := GetChainConfig(chainId)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Sorter, nil
}

func GetSimulatorAddress(chainId uint64) (common.Address, error) {
	config, err := GetChainConfig(chainId)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Simulator, nil
}

func GetMulticall3Address(chainId uint64) (common.Address, error) {
	config, err := GetChainConfig(chainId)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Multicall3, nil
}

func GetEIP712Domain(chainId uint64) (*apitypes.TypedDataDomain, error) {
	config, err := GetChainConfig(chainId)
	if err != nil {
		return nil, err
	}

	return &apitypes.TypedDataDomain{
		Name:              config.EIP712Domain.Name,
		Version:           config.EIP712Domain.Version,
		ChainId:           (*math.HexOrDecimal256)(new(big.Int).SetUint64(config.EIP712Domain.ChainId)),
		VerifyingContract: config.EIP712Domain.VerifyingContract,
	}, nil
}
