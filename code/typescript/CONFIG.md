# Shared TypeScript and ESLint Configuration

This project uses shared TypeScript and ESLint configurations to ensure consistency across all applications in the monorepo.

## TypeScript Configuration

The base TypeScript configuration is defined in `typescript/tsconfig.base.json`. All applications in the `apps/` directory extend this base configuration.

### How it works

1. The base configuration (`tsconfig.base.json`) contains common compiler options that are shared across all apps.
2. Each app has its own `tsconfig.app.json` that extends the base configuration using:
   ```json
   {
     "extends": "../../tsconfig.base.json",
     "compilerOptions": {
       // App-specific configurations...
     }
   }
   ```
3. App-specific configurations should only include settings that differ from the base configuration.

## ESLint Configuration

ESLint configuration is shared using the root configuration defined in `typescript/eslint.config.js`.

### How it works

1. The root ESLint configuration (`eslint.config.js`) defines common rules and plugins.
2. Each app has its own `eslint.config.js` that imports and extends the root configuration:

   ```javascript
   import rootConfig from '../../eslint.config.js';

   export default [
     ...rootConfig,
     {
       // App-specific configurations...
     },
   ];
   ```

3. App-specific configurations can override rules or add additional plugins as needed.

## Adding a New App

When adding a new app to the monorepo:

1. Create a basic `tsconfig.json` that references your app-specific configurations:

   ```json
   {
     "files": [],
     "references": [{ "path": "./tsconfig.app.json" }, { "path": "./tsconfig.node.json" }]
   }
   ```

2. Create a `tsconfig.app.json` that extends the base configuration:

   ```json
   {
     "extends": "../../tsconfig.base.json",
     "compilerOptions": {
       // App-specific path aliases and configurations
     },
     "include": ["src"]
   }
   ```

3. Create an `eslint.config.js` that extends the root configuration:

   ```javascript
   import rootConfig from '../../eslint.config.js';

   export default [
     ...rootConfig,
     {
       // App-specific rules and configurations
     },
   ];
   ```
