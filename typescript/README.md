# Atlas Config - TypeScript

TypeScript package that provides chain configurations for Atlas Protocol smart contracts. This package is part of the Atlas Protocol suite and is designed to work seamlessly with the Atlas SDK.

## Installation

```bash
npm install @fastlane-labs/atlas-config
# or
pnpm add @fastlane-labs/atlas-config
# or
yarn add @fastlane-labs/atlas-config
```

## Usage

### Basic Usage

```typescript
import { getChainConfig } from '@fastlane-labs/atlas-config';

// Get config for a specific chain
const sepoliaConfig = getChainConfig(11155111); // Sepolia testnet
console.log(sepoliaConfig.contracts.atlas); // Atlas contract address
console.log(sepoliaConfig.eip712Domain); // EIP-712 domain config
```

### Getting Supported Chain IDs

```typescript
import { getSupportedChainIds } from '@fastlane-labs/atlas-config';

const chainIds = getSupportedChainIds();
console.log("Supported chain IDs:", chainIds); // [137, 11155111, ...]
```

### Getting All Chain Configs

```typescript
import { getAllChainConfigs } from '@fastlane-labs/atlas-config';

const allConfigs = getAllChainConfigs();
console.log("All chain configs:", allConfigs.map(config => ({
  chainId: config.chainId,
  name: config.config.eip712Domain.name
})));
```

### Merging Custom Configurations

You can merge your own configurations with the default ones. This is useful for testing or using custom contract deployments.

```typescript
import { mergeChainConfigs } from '@fastlane-labs/atlas-config';

// Example: Updating a single contract address
const partialUpdate = {
  '11155111': { // Sepolia
    '1.0': {
      contracts: {
        atlas: { address: '0x1234567890123456789012345678901234567890' }
      }
    }
  }
};

// Example: Adding a new chain with complete configuration
const newChainConfig = {
  '999999': {
    '1.0': {
      contracts: {
        atlas: { address: '0x0987654321098765432109876543210987654321' },
        atlasVerification: { address: '0x0987654321098765432109876543210987654321' },
        sorter: { address: '0x0987654321098765432109876543210987654321' },
        simulator: { address: '0x0987654321098765432109876543210987654321' },
        multicall3: { address: '0x0987654321098765432109876543210987654321' }
      },
      eip712Domain: {
        name: 'New Test Chain',
        version: '1.0',
        chainId: 999999,
        verifyingContract: '0x1111111111111111111111111111111111111111'
      }
    }
  }
};

// Merge configurations
try {
  const mergedConfigs = mergeChainConfigs({
    ...partialUpdate,
    ...newChainConfig
  });
  console.log("Updated Sepolia config:", mergedConfigs['11155111']);
  console.log("New chain config:", mergedConfigs['999999']);
} catch (error) {
  console.error("Error merging configs:", error);
}
```

### Integration with Atlas SDK

The package is designed to work seamlessly with the Atlas SDK:

```typescript
import { getChainConfig } from '@fastlane-labs/atlas-config';
import { OperationBuilder } from '@fastlane-labs/atlas-sdk';

// Example: Creating a solver operation
OperationBuilder.newSolverOperation({
  from: "0x...",
  to: "0x...",
  value: BigInt(0),
  gas: BigInt(0),
  maxFeePerGas: BigInt(0),
  deadline: BigInt(0),
  solver: "0x...",
  control: "0x...",
  userOpHash: "0x...",
  bidToken: "0x...",
  bidAmount: BigInt(10000000000000000), // 0.01 ETH
  data: "0x...",
  signature: "0x..."
});
```

## Configuration Types

### Contract Configuration
```typescript
type ContractConfig = {
  atlas: string;
  atlasVerification: string;
  sorter: string;
  simulator: string;
  multicall3: string;
};
```

### EIP-712 Domain Configuration
```typescript
type EIP712Domain = {
  name: string;
  version: string;
  chainId: number;
  verifyingContract: string;
};
```

### Version Configuration
```typescript
type VersionConfig = {
  contracts: ContractConfig;
  eip712Domain: EIP712Domain;
};
```

## Error Handling

The package includes proper error handling for common scenarios:

```typescript
// Invalid chain ID
try {
  const config = getChainConfig(999999);
} catch (error) {
  console.error("Chain not supported:", error);
}

// Invalid version
try {
  const config = getChainConfig(11155111, '9.9');
} catch (error) {
  console.error("Version not found:", error);
}

// Incomplete configuration when merging
try {
  const incompleteConfig = {
    '888888': {
      '1.0': {
        contracts: {
          atlas: { address: '0x1234567890123456789012345678901234567890' }
        }
      }
    }
  };
  mergeChainConfigs(incompleteConfig);
} catch (error) {
  console.error("Expected error:", error);
}
```

## Development

To contribute to this package:

1. Clone the repository
2. Install dependencies: `pnpm install`
3. Run tests: `pnpm test`
4. Build: `pnpm build`

## License

MIT License - see the [LICENSE](../LICENSE) file for details.

