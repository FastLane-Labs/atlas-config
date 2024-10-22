import { chainConfig, mergeChainConfigs } from './index';
import { ChainConfig } from './types';
import { beforeEach, describe, expect, it } from '@jest/globals';

describe('mergeChainConfigs', () => {
  beforeEach(() => {
    // Reset chainConfig before each test
    (global as any).chainConfig = {
      '137': {
        contracts: {
          atlas: '0x892F8f6779ca6927c1A6Cc74319e03d2abEf18D5',
          atlasVerification: '0xc05DDBe9745ce9DB45C32F5e4C1DA7a3c4FDa220',
          sorter: '0x81f1E70A11A9E10Fa314cC093D149E5ec56EE97f',
          simulator: '0xfBc81A39459E0D82EC31B4e585f7A318AFAdB49B',
          multicall3: '0xcA11bde05977b3631167028862bE2a173976CA11',
        },
        eip712Domain: {
          name: 'AtlasVerification',
          version: '1.0',
          chainId: 137,
          verifyingContract: '0xc05DDBe9745ce9DB45C32F5e4C1DA7a3c4FDa220',
        },
      },
    };
  });

  it('should merge partial config for existing chain', () => {
    const result = mergeChainConfigs({
      '137': {
        contracts: {
          atlasVerification: '0x7000000000000000000000000000000000000000',
        },
      },
    });

    expect(result['137'].contracts.atlas).toBe('0x892F8f6779ca6927c1A6Cc74319e03d2abEf18D5');
    expect(result['137'].contracts.atlasVerification).toBe('0x7000000000000000000000000000000000000000');
  });

  it('should prioritize new complete config for existing chain', () => {
    const newConfig: ChainConfig = {
      contracts: {
        atlas: '0x3000000000000000000000000000000000000000',
        atlasVerification: '0x4000000000000000000000000000000000000000',
        sorter: '0x5000000000000000000000000000000000000000',
        simulator: '0x6000000000000000000000000000000000000000',
        multicall3: '0x7000000000000000000000000000000000000000',
      },
      eip712Domain: {
        name: 'New',
        version: '2',
        chainId: 1,
        verifyingContract: '0x8000000000000000000000000000000000000000',
      },
    };

    const result = mergeChainConfigs({ '1': newConfig });

    expect(result['1']).toEqual(newConfig);
  });

  it('should add new chain config', () => {
    const newConfig: ChainConfig = {
      contracts: {
        atlas:  '0x9000000000000000000000000000000000000000',
        atlasVerification: '0xa000000000000000000000000000000000000000',
        sorter: '0xb000000000000000000000000000000000000000',
        simulator: '0xc000000000000000000000000000000000000000',
        multicall3: '0xd000000000000000000000000000000000000000',
      },
      eip712Domain: {
        name: 'New Chain',
        version: '1',
        chainId: 2,
        verifyingContract: '0xe000000000000000000000000000000000000000',
      },
    };

    const result = mergeChainConfigs({ '2': newConfig });

    expect(result['2']).toEqual(newConfig);
  });

  it('should throw error when adding incomplete new chain config', () => {
    expect(() => {
      mergeChainConfigs({
        '3': {
          contracts: {
            atlas: '0xf000000000000000000000000000000000000000',
          },
        } as ChainConfig,
      });
    }).toThrow('Full chain configuration must be provided for new chainId: 3');
  });

  it('should merge multiple chains at once', () => {
    const existingChainId = Object.keys(chainConfig)[0];
    const newChainId = '999999';
    const providedConfigs = {
      [existingChainId]: {
        contracts: {
          atlas: '0x7000000000000000000000000000000000000000'
        }
      },
      [newChainId]: {
        contracts: {
          atlas: '0x8000000000000000000000000000000000000000',
          atlasVerification: '0x9000000000000000000000000000000000000000',
          sorter: '0xa000000000000000000000000000000000000000',
          simulator: '0xb000000000000000000000000000000000000000',
          multicall3: '0xc000000000000000000000000000000000000000'
        },
        eip712Domain: {
          name: 'New Chain',
          version: '1',
          chainId: 999999,
          verifyingContract: '0xd000000000000000000000000000000000000000'
        }
      }
    };

    const mergedConfigs = mergeChainConfigs(providedConfigs);

    expect(mergedConfigs[existingChainId].contracts.atlas).toBe('0x7000000000000000000000000000000000000000');
    expect(mergedConfigs[newChainId].contracts.atlas).toBe('0x8000000000000000000000000000000000000000');
  });

  it('should not change anything when merging an empty config', () => {
    const originalConfig = { ...chainConfig };
    const mergedConfigs = mergeChainConfigs({});

    expect(mergedConfigs).toEqual(originalConfig);
  });

  it('should return the same config when merging a config with no changes', () => {
    const originalConfig = { ...chainConfig };
    const mergedConfigs = mergeChainConfigs(originalConfig);

    expect(mergedConfigs).toEqual(originalConfig);
  });

  it('should merge a partial config that updates multiple fields', () => {
    const existingChainId = Object.keys(chainConfig)[0];
    const providedConfig = {
      [existingChainId]: {
        contracts: {
          atlas: '0x7000000000000000000000000000000000000000',
          sorter: '0x8000000000000000000000000000000000000000'
        },
        eip712Domain: {
          name: 'Updated Chain',
          version: '2.0'
        }
      }
    };

    const mergedConfigs = mergeChainConfigs(providedConfig);

    expect(mergedConfigs[existingChainId].contracts.atlas).toBe('0x7000000000000000000000000000000000000000');
    expect(mergedConfigs[existingChainId].contracts.sorter).toBe('0x8000000000000000000000000000000000000000');
    expect(mergedConfigs[existingChainId].eip712Domain.name).toBe('Updated Chain');
    expect(mergedConfigs[existingChainId].eip712Domain.version).toBe('2.0');
  });

  it('should not modify the original chainConfig when merging', () => {
    const originalConfig = JSON.parse(JSON.stringify(chainConfig));
    const existingChainId = Object.keys(chainConfig)[0];
    const providedConfig = {
      [existingChainId]: {
        contracts: {
          atlas: '0x7000000000000000000000000000000000000000'
        }
      }
    };

    mergeChainConfigs(providedConfig);

    expect(chainConfig).toEqual(originalConfig);
  });
});
