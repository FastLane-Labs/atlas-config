/// <reference types="node" />

import type { ChainConfig, PartialChainConfig } from "./types";
import { AtlasVLatest } from "./types";

let config: {
  [chainId: string]: { [version: string]: ChainConfig };
};

try {
  config = require("chain-config.json") as {
    [chainId: string]: { [version: string]: ChainConfig };
  };
} catch (e) {
  try {
    config = require("./chain-config.json") as {
      [chainId: string]: { [version: string]: ChainConfig };
    };
  } catch (e) {
    console.error("Error loading ./chain-config.json:", e);
    config = {};
  }
}

export const chainConfig: {
  [chainId: string]: { [version: string]: ChainConfig };
} = config;

export function getChainConfig(chainId: number, version?: string): ChainConfig {
  if (!version) {
    version = AtlasVLatest;
  }

  const config = chainConfig[chainId.toString()][version];
  if (!config) {
    throw new Error(`Chain configuration not found for chainId: ${chainId}`);
  }
  return config;
}

// New function to return all supported chainIds
export function getSupportedChainIds(): number[] {
  return Object.keys(chainConfig).map(Number);
}

// New function to return all chain configs as a list
export function getAllChainConfigs(): ChainConfig[] {
  return Object.values(chainConfig).reduce<ChainConfig[]>(
    (acc, versionConfig) => {
      return [...acc, ...Object.values(versionConfig)];
    },
    []
  );
}

// Updated function to merge provided chain configs with the original
export function mergeChainConfigs(providedConfigs: {
  [chainId: string]: { [version: string]: PartialChainConfig | ChainConfig };
}): { [chainId: string]: { [version: string]: ChainConfig } } {
  const mergedConfig = { ...chainConfig };

  for (const [chainId, versionConfig] of Object.entries(providedConfigs)) {
    for (const [version, providedConfig] of Object.entries(versionConfig)) {
      if (mergedConfig[chainId][version]) {
        // Existing chain
        if (isFullChainConfig(providedConfig)) {
          // Prioritize new complete config
          mergedConfig[chainId][version] = providedConfig as ChainConfig;
        } else {
          // Merge partial config
          mergedConfig[chainId][version] = {
            ...mergedConfig[chainId][version],
            ...providedConfig,
            contracts: {
              ...mergedConfig[chainId][version].contracts,
              ...(providedConfig.contracts as {
                [key: string]: { address: string };
              }),
            },
            eip712Domain: {
              ...mergedConfig[chainId][version].eip712Domain,
              ...providedConfig.eip712Domain,
            },
          };
        }
      } else {
        // New chain: ensure full config is provided
        if (isFullChainConfig(providedConfig)) {
          mergedConfig[chainId][version] = providedConfig as ChainConfig;
        } else {
          throw new Error(
            `Full chain configuration must be provided for new chainId: ${chainId}`
          );
        }
      }
    }
  }

  return mergedConfig;
}

// Helper function to check if a provided config is a full ChainConfig
function isFullChainConfig(
  config: PartialChainConfig | ChainConfig
): config is ChainConfig {
  return (
    typeof config.contracts === "object" &&
    typeof config.contracts.atlas === "string" &&
    typeof config.contracts.atlasVerification === "string" &&
    typeof config.contracts.sorter === "string" &&
    typeof config.contracts.simulator === "string" &&
    typeof config.contracts.multicall3 === "string" &&
    typeof config.eip712Domain === "object" &&
    typeof config.eip712Domain.name === "string" &&
    typeof config.eip712Domain.version === "string" &&
    typeof config.eip712Domain.chainId === "number" &&
    typeof config.eip712Domain.verifyingContract === "string"
  );
}
