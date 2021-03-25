/* eslint-disable */

module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  parserOptions: {
    ecmaFeatures: { jsx: true },
  },
  settings: {
    react: {
      version: "detect",
    },
  },
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/eslint-recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react/recommended",
    "plugin:jsx-a11y/recommended",
    "plugin:prettier/recommended",
  ],
  rules: {
    // Include .prettierrc.js rules
    "prettier/prettier": ["error", {}, { usePrettierrc: true }],
    "react/react-in-jsx-scope": "off",
    "react/prop-types": "off",
    "@typescript-eslint/explicit-module-boundary-types": "off", // turn this on in future
  },
};
