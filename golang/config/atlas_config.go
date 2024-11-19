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

type ContractsJSON struct {
	Atlas             string `json:"atlas"`
	AtlasVerification string `json:"atlasVerification"`
	Sorter            string `json:"sorter"`
	Simulator         string `json:"simulator"`
	Multicall3        string `json:"multicall3"`
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
	AtlasV100    = "1.0.0"
	AtlasV101    = "1.0.1"
	AtlasV110    = "1.1.0"
	AtlasVLatest = AtlasV110

	chainConfig map[string]map[string]ChainConfig
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

	var chainConfigJSON map[string]map[string]ChainConfigJSON
	err = json.Unmarshal(byteValue, &chainConfigJSON)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling config: %v", err))
	}

	chainConfig = make(map[string]map[string]ChainConfig)
	for chainID, versionConfig := range chainConfigJSON {
		chainConfig[chainID] = make(map[string]ChainConfig)

		for version, config := range versionConfig {
			chainConfig[chainID][version] = ChainConfig{
				Contracts: Contracts{
					Atlas:             common.HexToAddress(config.Contracts.Atlas),
					AtlasVerification: common.HexToAddress(config.Contracts.AtlasVerification),
					Sorter:            common.HexToAddress(config.Contracts.Sorter),
					Simulator:         common.HexToAddress(config.Contracts.Simulator),
					Multicall3:        common.HexToAddress(config.Contracts.Multicall3),
				},
				EIP712Domain: config.EIP712Domain,
			}
		}
	}
}

// GetChainConfig returns the chain configuration for the given chainId
func GetChainConfig(chainId uint64, version *string) (*ChainConfig, error) {
	if version == nil {
		version = &AtlasVLatest
	}

	configMutex.RLock()
	defer configMutex.RUnlock()

	versionConfig, ok := chainConfig[strconv.FormatUint(chainId, 10)]
	if !ok {
		return nil, fmt.Errorf("chain configuration not found for chainId: %d", chainId)
	}

	config, ok := versionConfig[*version]
	if !ok {
		return nil, fmt.Errorf("chain configuration not found for chainId: %d, version: %s", chainId, *version)
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
func GetAllChainConfigs() map[string]map[string]ChainConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Create a deep copy to avoid potential race conditions
	configCopy := make(map[string]map[string]ChainConfig, len(chainConfig))
	for k, v := range chainConfig {
		configCopy[k] = make(map[string]ChainConfig, len(v))
		for k2, v2 := range v {
			configCopy[k][k2] = v2
		}
	}

	return configCopy
}

// MergeChainConfigs merges the provided chain configs with the original
func MergeChainConfigs(providedConfigs map[string]map[string]ChainConfig) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	for chainId, versionConfig := range providedConfigs {
		for version, providedConfig := range versionConfig {
			if _, ok := chainConfig[chainId][version]; ok {
				// Existing chain/config
				if isFullChainConfig(providedConfig) {
					// Prioritize new complete config
					chainConfig[chainId][version] = providedConfig
				} else {
					// Merge partial config
					chainConfig[chainId][version] = mergeConfigs(chainConfig[chainId][version], providedConfig)
				}
			} else {
				// New chain: ensure full config is provided
				if !isFullChainConfig(providedConfig) {
					return fmt.Errorf("full chain configuration must be provided for new chainId: %s, version: %s", chainId, version)
				}

				if _, ok := chainConfig[chainId]; !ok {
					chainConfig[chainId] = make(map[string]ChainConfig)
				}

				chainConfig[chainId][version] = providedConfig
			}
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

func GetAtlasAddress(chainId uint64, version *string) (common.Address, error) {
	if version == nil {
		version = &AtlasVLatest
	}

	config, err := GetChainConfig(chainId, version)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Atlas, nil
}

func GetAtlasVerificationAddress(chainId uint64, version *string) (common.Address, error) {
	if version == nil {
		version = &AtlasVLatest
	}

	config, err := GetChainConfig(chainId, version)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.AtlasVerification, nil
}

func GetSorterAddress(chainId uint64, version *string) (common.Address, error) {
	if version == nil {
		version = &AtlasVLatest
	}

	config, err := GetChainConfig(chainId, version)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Sorter, nil
}

func GetSimulatorAddress(chainId uint64, version *string) (common.Address, error) {
	if version == nil {
		version = &AtlasVLatest
	}

	config, err := GetChainConfig(chainId, version)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Simulator, nil
}

func GetMulticall3Address(chainId uint64, version *string) (common.Address, error) {
	if version == nil {
		version = &AtlasVLatest
	}

	config, err := GetChainConfig(chainId, version)
	if err != nil {
		return common.Address{}, err
	}
	return config.Contracts.Multicall3, nil
}

func GetEIP712Domain(chainId uint64, version *string) (*apitypes.TypedDataDomain, error) {
	if version == nil {
		version = &AtlasVLatest
	}

	config, err := GetChainConfig(chainId, version)
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
