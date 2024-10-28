import { fail } from 'assert';
import { chainConfig, getChainConfig, getSupportedChainIds, getAllChainConfigs, mergeChainConfigs } from '../src/index';
import type { ChainConfig, PartialChainConfig } from '../src/types';
import { describe, it, expect, beforeEach } from '@jest/globals';

describe('Chain Config', () => {
  beforeEach(() => {
    // Reset chainConfig before each test with a complete configuration
    (global as any).chainConfig = {
      '1': {
        contracts: {
          atlas: { address: '0x1000000000000000000000000000000000000000' },
          atlasVerification: { address: '0x2000000000000000000000000000000000000000' },
          sorter: { address: '0x3000000000000000000000000000000000000000' },
          simulator: { address: '0x4000000000000000000000000000000000000000' },
          multicall3: { address: '0x5000000000000000000000000000000000000000' },
        },
        eip712Domain: {
          name: 'Initial',
          version: '1',
          chainId: 1,
          verifyingContract: '0x6000000000000000000000000000000000000000',
        },
      },
    };
  });

  it('should export chainConfig object', () => {
    expect(chainConfig).toBeInstanceOf(Object);
    expect(Object.keys(chainConfig).length).toBeGreaterThan(0);
  });

  describe('getChainConfig', () => {
    it('should return the correct config for a valid chainId', () => {
      // Assuming chainId 11155111 (Sepolia testnet) exists in the config
      const config = getChainConfig(11155111);
      
      // Check if the chainId exists in chainConfig
      expect(chainConfig).toHaveProperty('11155111');

      // Only run these tests if config is not undefined
      if (config) {
        expect(config).toBeInstanceOf(Object);
        expect(config).toHaveProperty('contracts');
        expect(config).toHaveProperty('eip712Domain');
      } else {
        fail('Config should not be undefined for chainId 11155111');
      }
    });

    it('should throw an error for an unknown chainId', () => {
      expect(() => getChainConfig(999999)).toThrow('Chain configuration not found for chainId: 999999');
    });

    it('should throw an error for chainId 0', () => {
      expect(() => getChainConfig(0)).toThrow('Chain configuration not found for chainId: 0');
    });

    it('should have the correct structure for a chain config', () => {
      const chainId = 11155111; // Use a known valid chainId
      const config = getChainConfig(chainId);
      
      expect(config).toBeInstanceOf(Object);
      expect(config?.contracts).toEqual(
        expect.objectContaining({
          atlas: expect.any(String),
          atlasVerification: expect.any(String),
          sorter: expect.any(String),
          simulator: expect.any(String),
          multicall3: expect.any(String)
        })
      );
      expect(config?.eip712Domain).toEqual(
        expect.objectContaining({
          name: expect.any(String),
          version: expect.any(String),
          chainId: expect.any(Number),
          verifyingContract: expect.any(String)
        })
      );
    });
  });

  describe('getSupportedChainIds', () => {
    it('should return an array of numbers', () => {
      const chainIds = getSupportedChainIds();
      expect(Array.isArray(chainIds)).toBe(true);
      expect(chainIds.length).toBeGreaterThan(0);
      chainIds.forEach(id => expect(typeof id).toBe('number'));
    });

    it('should return all chain IDs from the chainConfig', () => {
      const chainIds = getSupportedChainIds();
      const configChainIds = Object.keys(chainConfig).map(Number);
      expect(chainIds).toEqual(expect.arrayContaining(configChainIds));
      expect(chainIds.length).toBe(configChainIds.length);
    });
  });

  describe('getAllChainConfigs', () => {
    it('should return an array of ChainConfig objects', () => {
      const configs = getAllChainConfigs();
      expect(Array.isArray(configs)).toBe(true);
      expect(configs.length).toBeGreaterThan(0);
      configs.forEach(config => {
        expect(config).toHaveProperty('contracts');
        expect(config).toHaveProperty('eip712Domain');
      });
    });

    it('should return all configs from the chainConfig', () => {
      const configs = getAllChainConfigs();
      expect(configs.length).toBe(Object.keys(chainConfig).length);
    });
  });

  describe('mergeChainConfigs', () => {
    it('should merge partial config for existing chain', () => {
      const existingChainId = Object.keys(chainConfig)[0];  // Get the first chain ID from the existing config
      const providedConfigs: { [chainId: string]: PartialChainConfig } = {
        [existingChainId]: {
          contracts: {
            atlas: '0x7000000000000000000000000000000000000000'
          }
        }
      };

      const mergedConfigs = mergeChainConfigs(providedConfigs);

      expect(mergedConfigs[existingChainId].contracts.atlas).toBe('0x7000000000000000000000000000000000000000');
      expect(mergedConfigs[existingChainId].contracts.atlasVerification).toBe(chainConfig[existingChainId].contracts.atlasVerification);
    });

    it('should add new chain config', () => {
      const newChainId = '999999';
      const newConfig: ChainConfig = {
        contracts: {
          atlas: '0x9000000000000000000000000000000000000000',
          atlasVerification: '0xa000000000000000000000000000000000000000',
          sorter: '0xb000000000000000000000000000000000000000',
          simulator: '0xc000000000000000000000000000000000000000',
          multicall3: '0xd000000000000000000000000000000000000000'
        },
        eip712Domain: {
          name: 'New Chain',
          version: '1',
          chainId: 999999,
          verifyingContract: '0xe000000000000000000000000000000000000000'
        }
      };

      const providedConfigs = { [newChainId]: newConfig };
      const mergedConfigs = mergeChainConfigs(providedConfigs);

      expect(mergedConfigs[newChainId]).toEqual(newConfig);
    });

    it('should throw error when adding incomplete new chain config', () => {
      const newChainId = '888888';
      const incompleteConfig: PartialChainConfig = {
        contracts: {
          atlas: '0xf000000000000000000000000000000000000000'
        }
      };

      expect(() => mergeChainConfigs({ [newChainId]: incompleteConfig })).toThrow(
        `Full chain configuration must be provided for new chainId: ${newChainId}`
      );
    });
  });
});
