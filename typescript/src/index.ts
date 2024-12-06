/// <reference types="node" />

import type { ChainConfig, VersionConfig, PartialChainConfig } from './types';

let config: { [chainId: string]: ChainConfig } = {};

try {
  // Try different possible paths for the config file
  const paths = [
    '../configs/chain-configs-multi-version.json',
    './chain-configs-multi-version.json',
    'chain-configs-multi-version.json'
  ];
  
  for (const path of paths) {
    try {
      config = require(path);
      break;
    } catch (e) {
      continue;
    }
  }
  
  if (Object.keys(config).length === 0) {
    throw new Error('Could not load config file');
  }
} catch (e) {
  console.error('Error loading chain-configs-multi-version.json:', e);
}

export const chainConfig: { [chainId: string]: ChainConfig } = config;

export function getChainConfig(chainId: number, version?: string): VersionConfig {
  const chainConf = chainConfig[chainId.toString()];
  if (!chainConf) {
    throw new Error(`Chain configuration not found for chainId: ${chainId}`);
  }

  // If version is not specified, return the latest version
  if (!version) {
    const versions = Object.keys(chainConf).sort((a, b) => {
      return parseFloat(b) - parseFloat(a);
    });
    version = versions[0];
  }

  const versionConfig = chainConf[version];
  if (!versionConfig) {
    throw new Error(`Version ${version} not found for chainId: ${chainId}`);
  }

  return versionConfig;
}

export function getSupportedChainIds(): number[] {
  return Object.keys(chainConfig).map(Number);
}

export function getVersionsForChain(chainId: number): string[] {
  const chainConf = chainConfig[chainId.toString()];
  if (!chainConf) {
    return [];
  }
  return Object.keys(chainConf);
}

export function getAllChainConfigs(): { chainId: number; version: string; config: VersionConfig }[] {
  const configs: { chainId: number; version: string; config: VersionConfig }[] = [];
  
  for (const [chainId, chainConf] of Object.entries(chainConfig)) {
    for (const [version, versionConfig] of Object.entries(chainConf)) {
      configs.push({
        chainId: Number(chainId),
        version,
        config: versionConfig
      });
    }
  }
  
  return configs;
}

function isFullVersionConfig(config: any): config is VersionConfig {
  return (
    typeof config.contracts === 'object' &&
    typeof config.contracts.atlas === 'string' &&
    typeof config.contracts.atlasVerification === 'string' &&
    typeof config.contracts.sorter === 'string' &&
    typeof config.contracts.simulator === 'string' &&
    typeof config.contracts.multicall3 === 'string' &&
    typeof config.eip712Domain === 'object' &&
    typeof config.eip712Domain.name === 'string' &&
    typeof config.eip712Domain.version === 'string' &&
    typeof config.eip712Domain.chainId === 'number' &&
    typeof config.eip712Domain.verifyingContract === 'string'
  );
}

export function mergeChainConfigs(providedConfigs: { [chainId: string]: PartialChainConfig }): { [chainId: string]: ChainConfig } {
  const mergedConfig = { ...chainConfig };
  
  for (const [chainId, providedVersions] of Object.entries(providedConfigs)) {
    if (!mergedConfig[chainId]) {
      mergedConfig[chainId] = {};
    }
    
    for (const [version, providedVersionConfig] of Object.entries(providedVersions)) {
      if (mergedConfig[chainId][version]) {
        // Existing version
        if (isFullVersionConfig(providedVersionConfig)) {
          // Prioritize new complete config
          mergedConfig[chainId][version] = providedVersionConfig;
        } else {
          // Merge partial config
          mergedConfig[chainId][version] = {
            ...mergedConfig[chainId][version],
            contracts: {
              ...mergedConfig[chainId][version].contracts,
              ...providedVersionConfig.contracts,
            },
            eip712Domain: {
              ...mergedConfig[chainId][version].eip712Domain,
              ...providedVersionConfig.eip712Domain,
            },
          };
        }
      } else {
        // New version: ensure full config is provided
        if (isFullVersionConfig(providedVersionConfig)) {
          mergedConfig[chainId][version] = providedVersionConfig;
        } else {
          throw new Error(`Full version configuration must be provided for new version ${version} on chainId: ${chainId}`);
        }
      }
    }
  }

  return mergedConfig;
}
