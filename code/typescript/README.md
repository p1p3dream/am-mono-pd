# TypeScript Project with ESLint and Prettier

This project demonstrates how to use ESLint and Prettier for TypeScript code quality and formatting.

## Setup

The project is already set up with:

- ESLint for TypeScript linting (using the new eslint.config.js format for ESLint v9+)
- Prettier for code formatting
- Configuration files for both tools

## Running the Tools

### Format Code with Prettier

To format your TypeScript code with Prettier:

```bash
# Format all TypeScript files
npx prettier --write "**/*.{ts,tsx}"

# Check formatting without changing files
npx prettier --check "**/*.{ts,tsx}"
```

### Lint Code with ESLint

This project uses ESLint v9+ with the new configuration format in `eslint.config.js`. The configuration includes:

- TypeScript-specific rules
- Integration with Prettier
- Ignore patterns for node_modules, dist, etc.

```bash
# Check for linting issues
npx eslint "**/*.{ts,tsx}"

# Fix automatically fixable issues
npx eslint "**/*.{ts,tsx}" --fix
```

### Quick Fix Script

A custom script has been provided to quickly format and fix common TypeScript issues:

```bash
node lint-example.js
```

This script:

1. Runs Prettier on the example TypeScript file
2. Fixes common TypeScript issues:
   - Adds missing return types
   - Replaces `any` with `unknown`
   - Adds type annotations
   - Handles unused variables

## Customizing the Configuration

- ESLint: Edit `eslint.config.js` to change linting rules and ignore patterns
- Prettier: Edit `.prettierrc.json` to change formatting rules

```sh
make install
```
