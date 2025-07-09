import { dirname } from "path";
import { fileURLToPath } from "url";
import rootConfig from "../../eslint.config.js";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default [
  ...rootConfig,
  {
    files: ["**/*.{ts,tsx}"],
    rules: {
      // Add any app-specific overrides here
    },
  },
];
