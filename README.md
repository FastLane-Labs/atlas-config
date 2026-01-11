# Chain Configurations

A cross-language configuration package that provides a single source of truth for chain configurations, specifically for Atlas smart contracts. This package is designed to be used across multiple programming languages, including TypeScript, Go, and Rust, and is intended to be consumed by integrators or used with the Atlas SDK, which is available in both TypeScript and Go versions.

## Installation

### TypeScript (Node.js)

Install the package via npm:

```bash
npm install @fastlane-labs/atlas-config
```

### Go

Get the package using `go get`:

```bash
go get github.com/fastLane-labs/atlas-config/golang/config
```

## Usage

### TypeScript

Import the `getChainConfig` function from the package:

```typescript
import { getChainConfig } from '@fastlane-labs/atlas-config';

const config = getChainConfig(137); // Polygon Mainnet
console.log(config.name); // Outputs: "Polygon Mainnet"
```
For more detailed documentation on using the TypeScript package, please refer to the [TypeScript README](./typescript/README.md).

### Go

Import the package and use the `GetChainConfig` function:

```go
import (
    "fmt"
    "github.com/fastlane-labs/atlas-config/golang/config"
)

func main() {
    config, err := chainconfig.GetChainConfig(137)
    if err != nil {
        // Handle error
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(config.Name) // Outputs: "Polygon Mainnet"
}
```

For more detailed documentation on using the Go package, please refer to the [Go README](./golang/README.md).

## Configuration Structure

The configuration data is stored in a JSON file and includes details like contract addresses and EIP-712 domain configurations. This data is parsed and made available through language-specific interfaces.

## Supported Chains

### Mainnets
- **Polygon Mainnet** (Chain ID: `137`)
- **Arbitrum Mainnet** (Chain ID: `42161`)
- **Base Mainnet** (Chain ID: `8453`)

### Testnets
- **Base Sepolia** (Chain ID: `84532`)
- **Polygon Amoy** (Chain ID: `80002`)
- **Ethereum Sepolia** (Chain ID: `11155111`)

## Development
To contribute or make changes to the configurations:

1. Update the configuration file located at `configs/chain-config.json`.
2. Build or compile the packages in their respective language directories.
3. Publish the updated packages to the appropriate registries.
   - **TypeScript**: Run `npm publish` in the `typescript` directory.
   - **Go**: Tag the new version and push to GitHub.
   - **Rust**: Publish to [crates.io](https://crates.io/) or update the Git repository.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

This structure allows you to maintain a single source of truth for your chain configurations while providing easy-to-use packages for multiple programming languages. Remember to update the version numbers and publish the packages whenever you make changes to the configuration.

