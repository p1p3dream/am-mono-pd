import { dirname } from "path";
import { fileURLToPath } from "url";
import query from "eslint-plugin-react-query";
import rootConfig from "../../eslint.config.js";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default [
  ...rootConfig,
  {
    files: ["**/*.{ts,tsx}"],
    plugins: {
      "react-query": query,
    },
    rules: {
      "react-query/exhaustive-deps": "warn",
      // Add any other app-specific overrides here
    },
  },
];
