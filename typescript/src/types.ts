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
