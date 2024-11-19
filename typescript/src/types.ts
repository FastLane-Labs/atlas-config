export const AtlasV100 = "1.0.0";
export const AtlasV101 = "1.0.1";
export const AtlasV110 = "1.1.0";
export const AtlasVLatest = AtlasV110;

export type ChainConfig = {
  contracts: {
    atlas: string;
    atlasVerification: string;
    sorter: string;
    simulator: string;
    multicall3: string;
  };
  eip712Domain: {
    name: string;
    version: string;
    chainId: number;
    verifyingContract: string;
  };
};

export type PartialChainConfig = {
  contracts?: Partial<ChainConfig['contracts']>;
  eip712Domain?: Partial<ChainConfig['eip712Domain']>;
};
