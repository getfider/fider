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
    "plugin:prettier/recommended",
  ],
  rules: {
    // Include .prettierrc.js rules
    "prettier/prettier": ["error", {}, { usePrettierrc: true }],
    "react/react-in-jsx-scope": "off",
    "@typescript-eslint/no-unused-vars": ["error", { varsIgnorePattern: "_" }],
    "react/prop-types": "off",
    "react/jsx-no-target-blank": ["error", { allowReferrer: true }],
    "@typescript-eslint/no-explicit-any": "off", // TODO: turn this on in future
    "@typescript-eslint/explicit-module-boundary-types": "off", // TODO: turn this on in future
  },
}
