import config from '../../configs/chain-config.json';

export interface ChainConfig {
  contracts: {
    atlas: {
      address: string;
    };
    atlasVerification: {
      address: string;
    };
    sorter: {
      address: string;
    };
    simulator: {
      address: string;
    };
    multicall3: {
      address: string;
    };
  };
  eip712Domain: {
    name: string;
    version: string;
    chainId: number;
    verifyingContract: string;
  };
}

export const chainConfig: { [chainId: string]: ChainConfig } = config;

export function getChainConfig(chainId: number): ChainConfig | undefined {
  return chainConfig[chainId.toString()];
}