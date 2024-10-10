import { chainConfig, getChainConfig, ChainConfig } from '../dist/index';

describe('Chain Config', () => {
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
          atlas: expect.any(Object),
          atlasVerification: expect.any(Object),
          sorter: expect.any(Object),
          simulator: expect.any(Object),
          multicall3: expect.any(Object)
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
});