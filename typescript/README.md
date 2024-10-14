# @fastlane-labs/atlas-config

A configuration package for Atlas Protocol, providing essential chain configuration data for EVM networks. This package is designed to be used in conjunction with the [atlas-sdk](https://www.npmjs.com/package/@fastlane-labs/atlas-sdk) and contains all relevant smart contract addresses for the Atlas protocol.

## Installation

Install the package using npm:

```bash
npm install @fastlane-labs/atlas-config
```

```bash
yarn add @fastlane-labs/atlas-config
```

## Usage

Import and use the configuration data in your TypeScript or JavaScript project

```typescript
import { chainConfig, getChainConfig, getAllChainConfigs, getSupportedChainIds, mergeChainConfigs } from '@fastlane-labs/atlas-config';
```

### Access the entire chain configuration

```typescript
console.log(chainConfig);
```

### Get configuration for a specific chain (e.g., Polygon mainnet)

```typescript
const polygonConfig = getChainConfig(137);
console.log(polygonConfig);
```

### Attempting to get configuration for an unknown chain will throw an error

```typescript
try {
  const unknownConfig = getChainConfig(999999);
} catch (error) {
  console.error(error); // "Chain configuration not found for chainId: 999999"
}
```

### Get all chain configurations

```typescript
const allConfigs = getAllChainConfigs();
console.log(allConfigs);
```

### Get supported chain IDs

```typescript
const supportedChainIds = getSupportedChainIds();
console.log(supportedChainIds);
```

### Merge provided chain configurations with existing ones

```typescript
const additionalConfig = {
  137: {
    contracts: {
      atlas: { address: '0x1234567890123456789012345678901234567890' }
    }
  },
  999999: {
    contracts: {
      atlas: { address: '0x0987654321098765432109876543210987654321' },
      atlasVerification: { address: '0x0987654321098765432109876543210987654321' },
      sorter: { address: '0x0987654321098765432109876543210987654321' },
      simulator: { address: '0x0987654321098765432109876543210987654321' },
      multicall3: { address: '0x0987654321098765432109876543210987654321' }
    },
    eip712Domain: {
      name: 'New Chain',
      version: '1',
      chainId: 999999,
      verifyingContract: '0x1111111111111111111111111111111111111111'
    }
  }
};

const mergedConfigs = mergeChainConfigs(additionalConfig);
console.log(mergedConfigs);
```

````

## Configuration Structure

The `ChainConfig` interface describes the structure of the configuration for each supported chain:

```typescript
interface ChainConfig {
  contracts: {
    atlas: object;
    atlasVerification: object;
    sorter: object;
    simulator: object;
    multicall3: object;
  };
  eip712Domain: {
    name: string;
    version: string;
    chainId: number;
    verifyingContract: string;
  };
}
````

## Supported Chains

This package includes configurations for various Ethereum networks. Use the `getChainConfig` function with the appropriate chain ID to retrieve the configuration for a specific network.

### Mainnets

- Polygon (Chain ID: 137)
- Base (Chain ID: 8453)
- Arbitrum (Chain ID: 42161)

### Testnets

- Ethereum Sepolia (Chain ID: 11155111)
- Polygon Amoy (Chain ID: 80002)

Each chain configuration includes contract addresses and EIP-712 domain information specific to that network. Use the appropriate chain ID when calling `getChainConfig()` to retrieve the configuration for your desired network.

## Integration with atlas-sdk

This configuration package is designed to work seamlessly with the [Atlas Typescript SDK](https://www.npmjs.com/package/@fastlane-labs/atlas-sdk). It provides all the necessary smart contract addresses and network-specific information required for interacting with the Atlas protocol.

## Contributing

If you'd like to contribute to this project, please submit a pull request or open an issue on our GitHub repository.

## License

This project is licensed under the [MIT License](LICENSE).

## Support

For questions, issues, or feature requests, please open an issue on our GitHub repository or contact our support team at [support@fastlanelabs.com](mailto\:support@fastlanelabs.com).

## Disclaimer

This package is part of the Atlas protocol ecosystem. Make sure to use it in conjunction with other Atlas-related packages and follow best practices for blockchain development and security.

