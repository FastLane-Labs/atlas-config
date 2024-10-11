package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

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

var chainConfig map[string]ChainConfig

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

func GetChainConfig(chainId uint64) (*ChainConfig, error) {
	config, ok := chainConfig[strconv.FormatUint(chainId, 10)]
	if !ok {
		return nil, fmt.Errorf("chain configuration not found for chainId: %d", chainId)
	}
	return &config, nil
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
