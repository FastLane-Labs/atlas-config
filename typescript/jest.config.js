const path = require('path');

const configPath = path.resolve(__dirname, '../configs/chain-config.json');

module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src', '<rootDir>/test'],
  moduleNameMapper: {
    '^chain-config.json$': configPath,
    '^./chain-config.json$': configPath,
  },
  moduleDirectories: ['node_modules', 'src', path.resolve(__dirname, '../../configs')],
  modulePaths: ['.', path.resolve(__dirname, '../../configs')],
  transform: {
    '^.+\\.tsx?$': ['ts-jest', {
      tsconfig: 'tsconfig.json',
    }],
  },
};
