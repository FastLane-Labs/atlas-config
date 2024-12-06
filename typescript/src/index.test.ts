import { describe, expect, it } from '@jest/globals';
import {
  getChainConfig,
  getSupportedChainIds,
  getVersionsForChain,
  getAllChainConfigs,
  mergeChainConfigs,
  chainConfig
} from './index';
import type { PartialChainConfig, PartialVersionConfig } from './types';

describe('Chain Config Tests', () => {
  describe('getChainConfig', () => {
    it('should return the latest version config when no version specified', () => {
      const config = getChainConfig(84532);
      expect(config.contracts.atlas).toBe('0xda92A5d54149c30a1cA54155a9fF55F32c2E0570'); // v1.2
    });

    it('should return specific version config when version is specified', () => {
      const config = getChainConfig(84532, '1.1');
      expect(config.contracts.atlas).toBe('0xa55051bd82eFeA1dD487875C84fE9c016859659B'); // v1.1
    });

    it('should throw error for invalid chain ID', () => {
      expect(() => getChainConfig(999999)).toThrow('Chain configuration not found');
    });

    it('should throw error for invalid version', () => {
      expect(() => getChainConfig(84532, '9.9')).toThrow('Version 9.9 not found');
    });
  });

  describe('getSupportedChainIds', () => {
    it('should return all supported chain IDs', () => {
      const chainIds = getSupportedChainIds();
      expect(chainIds).toContain(137);  // Polygon
      expect(chainIds).toContain(84532); // Base Sepolia
      expect(chainIds.length).toBeGreaterThan(0);
    });
  });

  describe('getVersionsForChain', () => {
    it('should return all versions for a chain', () => {
      const versions = getVersionsForChain(84532);
      expect(versions).toContain('1.1');
      expect(versions).toContain('1.2');
      expect(versions.length).toBe(2);
    });

    it('should return empty array for invalid chain', () => {
      const versions = getVersionsForChain(999999);
      expect(versions).toEqual([]);
    });
  });

  describe('getAllChainConfigs', () => {
    it('should return all chain configs with versions', () => {
      const configs = getAllChainConfigs();
      
      // Find Base Sepolia v1.2 config
      const baseSepoliaV12 = configs.find(c => 
        c.chainId === 84532 && c.version === '1.2'
      );
      
      expect(baseSepoliaV12).toBeDefined();
      expect(baseSepoliaV12?.config.contracts.atlas)
        .toBe('0xda92A5d54149c30a1cA54155a9fF55F32c2E0570');
    });
  });

  describe('mergeChainConfigs', () => {
    it('should merge partial configs correctly', () => {
      const partialConfig: { [chainId: string]: PartialChainConfig } = {
        '84532': {
          '1.2': {
            contracts: {
              atlas: '0xNewAddress'
            }
          } as PartialVersionConfig
        }
      };

      const merged = mergeChainConfigs(partialConfig);
      expect(merged['84532']['1.2'].contracts.atlas).toBe('0xNewAddress');
      // Other contracts should remain unchanged
      expect(merged['84532']['1.2'].contracts.simulator)
        .toBe('0xabdBe9C2b1Dd55e5A661410E2f81a6B8C1f12D0C');
    });

    it('should throw error when adding new chain without full config', () => {
      const invalidConfig: { [chainId: string]: PartialChainConfig } = {
        '999999': {
          '1.0': {
            contracts: {
              atlas: '0xNewAddress'
            }
          } as PartialVersionConfig
        }
      };

      expect(() => mergeChainConfigs(invalidConfig))
        .toThrow('Full version configuration must be provided');
    });

    it('should add new version when full config provided', () => {
      const newVersionConfig: { [chainId: string]: PartialChainConfig } = {
        '84532': {
          '2.0': {
            contracts: {
              atlas: '0xNewAtlas',
              atlasVerification: '0xNewVerification',
              sorter: '0xNewSorter',
              simulator: '0xNewSimulator',
              multicall3: '0xNewMulticall'
            },
            eip712Domain: {
              name: 'AtlasVerification',
              version: '1.0',
              chainId: 84532,
              verifyingContract: '0xNewVerification'
            }
          }
        }
      };

      const merged = mergeChainConfigs(newVersionConfig);
      expect(merged['84532']['2.0']).toBeDefined();
      expect(merged['84532']['2.0'].contracts.atlas).toBe('0xNewAtlas');
    });

    it('should handle multiple versions for the same chain', () => {
      const multiVersionConfig: { [chainId: string]: PartialChainConfig } = {
        '84532': {
          '1.1': {
            contracts: {
              atlas: '0xVersion1',
              atlasVerification: '0xVerification1',
              sorter: '0xSorter1',
              simulator: '0xSimulator1',
              multicall3: '0xMulticall3'
            },
            eip712Domain: {
              name: 'AtlasVerification',
              version: '1.0',
              chainId: 84532,
              verifyingContract: '0xVerification1'
            }
          },
          '1.2': {
            contracts: {
              atlas: '0xVersion2',
              atlasVerification: '0xVerification2',
              sorter: '0xSorter2',
              simulator: '0xSimulator2',
              multicall3: '0xMulticall3'
            },
            eip712Domain: {
              name: 'AtlasVerification',
              version: '1.0',
              chainId: 84532,
              verifyingContract: '0xVerification2'
            }
          }
        }
      };

      const merged = mergeChainConfigs(multiVersionConfig);
      expect(merged['84532']['1.1'].contracts.atlas).toBe('0xVersion1');
      expect(merged['84532']['1.2'].contracts.atlas).toBe('0xVersion2');
    });
  });
});
