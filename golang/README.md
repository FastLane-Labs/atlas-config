# @fastlane-labs/atlas-config (Go Version)

A configuration package for Atlas Protocol, providing essential chain configuration data for EVM networks. This package is designed to be used with the Atlas protocol and contains all relevant smart contract addresses for various supported chains.

## Installation

To use this configuration package in your Go project, clone or download the source code into your desired location.

## Usage

Import and use the configuration data in your Go project:

```go
import (
	"fmt"
	"github.com/fastlane-labs/atlas-config/config"
)
```

### Access the entire chain configuration

```go
chainConfig := config.GetAllChainConfigs()
fmt.Println(chainConfig)
```

### Get configuration for a specific chain (e.g., Polygon mainnet)

```go
polygonConfig, err := config.GetChainConfig(137)
if err != nil {
	fmt.Println("Error:", err)
} else {
	fmt.Println(polygonConfig)
}
```

### Attempting to get configuration for an unknown chain will return an error

```go
unknownConfig, err := config.GetChainConfig(999999)
if err != nil {
	fmt.Println("Error:", err) // "chain configuration not found for chainId: 999999"
} else {
	fmt.Println(unknownConfig)
}
```

### Get all chain configurations

```go
allConfigs := config.GetAllChainConfigs()
fmt.Println(allConfigs)
```

### Get supported chain IDs

```go
supportedChainIds := config.GetSupportedChainIds()
fmt.Println(supportedChainIds)
```

### Merge provided chain configurations with existing ones

```go
additionalConfig := map[string]config.ChainConfig{
	"137": {
		Contracts: config.Contracts{
			Atlas: config.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
	},
	"999999": {
		Contracts: config.Contracts{
			Atlas:             config.HexToAddress("0x0987654321098765432109876543210987654321"),
			AtlasVerification: config.HexToAddress("0x0987654321098765432109876543210987654321"),
			Sorter:            config.HexToAddress("0x0987654321098765432109876543210987654321"),
			Simulator:         config.HexToAddress("0x0987654321098765432109876543210987654321"),
			Multicall3:        config.HexToAddress("0x0987654321098765432109876543210987654321"),
		},
		EIP712Domain: config.EIP712Domain{
			Name:              "New Chain",
			Version:           "1",
			ChainId:           999999,
			VerifyingContract: "0x1111111111111111111111111111111111111111",
		},
	},
}

err := config.MergeChainConfigs(additionalConfig)
if err != nil {
	fmt.Println("Error merging configs:", err)
} else {
	fmt.Println("Configs merged successfully")
}
```

## Configuration Structure

The `ChainConfig` struct describes the structure of the configuration for each supported chain:

```go
type ChainConfig struct {
	Contracts    Contracts
	EIP712Domain EIP712Domain
}

type Contracts struct {
	Atlas             common.Address
	AtlasVerification common.Address
	Sorter            common.Address
	Simulator         common.Address
	Multicall3        common.Address
}

type EIP712Domain struct {
	Name              string
	Version           string
	ChainId           uint64
	VerifyingContract string
}
```

## Supported Chains

This package includes configurations for various Ethereum networks. Use the `GetChainConfig` function with the appropriate chain ID to retrieve the configuration for a specific network.

### Mainnets
- Polygon (Chain ID: 137)
- Base (Chain ID: 8453)
- Arbitrum (Chain ID: 42161)

### Testnets
- Ethereum Sepolia (Chain ID: 11155111)
- Polygon Amoy (Chain ID: 80002)

Each chain configuration includes contract addresses and EIP-712 domain information specific to that network. Use the appropriate chain ID when calling `GetChainConfig()` to retrieve the configuration for your desired network.

## Integration with Atlas Protocol

This configuration package is designed to work seamlessly with the Atlas protocol. It provides all the necessary smart contract addresses and network-specific information required for interacting with the Atlas protocol.

## Contributing

If you'd like to contribute to this project, please submit a pull request or open an issue on our GitHub repository.

## License

This project is licensed under the [MIT License](LICENSE).

## Support

For questions, issues, or feature requests, please open an issue on our GitHub repository or contact our support team at support@fastlanelabs.com.

## Disclaimer

This package is part of the Atlas protocol ecosystem. Make sure to use it in conjunction with other Atlas-related packages and follow best practices for blockchain development and security.