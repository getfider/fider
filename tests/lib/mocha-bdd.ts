import { ISuiteCallbackContext } from 'mocha';

export const specification = (name: string, callback: (this: ISuiteCallbackContext) => void) => {
  return describe(`Specification: ${name}`, callback);
};

export const when = (condition: string, callback: (this: ISuiteCallbackContext) => void) => {
  return describe(`when ${condition}`, callback);
};

export const then = (condition: string, callback: (this: ISuiteCallbackContext) => void) => {
  return it(`then ${condition}`, callback);
};

export const action = before;
