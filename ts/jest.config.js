module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src'],
  testMatch: ['**/__tests__/**/*.ts', '**/?(*.)+(spec|test).ts'],
  transform: {
    '^.+\\.ts$': 'ts-jest',
  },
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  moduleNameMapper: {
    // These are mock files that replace the WASI module imports during testing
    '^wasi:otel/logs@0.2.0-rc.2$': '<rootDir>/__mocks__/wasi-otel-logs.ts',
    '^wasi:otel/types@0.2.0-rc.2$': '<rootDir>/__mocks__/wasi-otel-types.ts',
  },
};
