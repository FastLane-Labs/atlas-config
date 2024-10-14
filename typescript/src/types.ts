export type ChainConfig = {
  contracts: {
    atlas: { address: string };
    atlasVerification: { address: string };
    sorter: { address: string };
    simulator: { address: string };
    multicall3: { address: string };
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
