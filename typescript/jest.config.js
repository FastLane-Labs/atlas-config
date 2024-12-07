const path = require('path');

module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  testMatch: ['**/*.test.ts'],
  moduleFileExtensions: ['ts', 'js', 'json'],
  moduleNameMapper: {
    'chain-configs-multi-version.json': path.resolve(__dirname, '../configs/chain-configs-multi-version.json'),
    './chain-configs-multi-version.json': path.resolve(__dirname, '../configs/chain-configs-multi-version.json')
  },
  moduleDirectories: ['node_modules', 'src', path.resolve(__dirname, '../configs')]
};
