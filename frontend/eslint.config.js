import js from "@eslint/js";
import globals from "globals";
import pluginReact from "eslint-plugin-react";
import { defineConfig } from "eslint/config";

export default defineConfig([
  {
    files: ["**/*.{js,mjs,cjs,jsx}"],
    plugins: { js },
    extends: ["js/recommended"],
    languageOptions: {
      globals: globals.browser,
    },
  },
  pluginReact.configs.flat.recommended,

  // Enable Jest and Node globals in test files
  {
    files: ["**/*.test.{js,jsx}"],
    languageOptions: {
      globals: { ...globals.jest, ...globals.node },
    },
  },

  // React and general linting rules for JS and JSX files
  {
    files: ["**/*.{js,jsx}"],
    settings: {
      react: {
        version: "detect",  // Automatically detect the React version
      },
    },
    rules: {
      "react/react-in-jsx-scope": "off",
      "react/prop-types": "warn",
      "no-unused-vars": ["warn", { argsIgnorePattern: "^_" }],
    },
  },
]);
