export type ContractConfig = {
  atlas: string;
  atlasVerification: string;
  sorter: string;
  simulator: string;
  multicall3: string;
};

export type EIP712Domain = {
  name: string;
  version: string;
  chainId: number;
  verifyingContract: string;
};

export type VersionConfig = {
  contracts: ContractConfig;
  eip712Domain: EIP712Domain;
};

export type ChainConfig = {
  [version: string]: VersionConfig;
};

export type PartialContractConfig = Partial<ContractConfig>;
export type PartialEIP712Domain = Partial<EIP712Domain>;

export type PartialVersionConfig = {
  contracts?: PartialContractConfig;
  eip712Domain?: PartialEIP712Domain;
};

export type PartialChainConfig = {
  [version: string]: PartialVersionConfig;
};
