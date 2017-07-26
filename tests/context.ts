import { Browser } from './lib';
import { AllPages } from './pages';

export let pages: AllPages;
export let browser: Browser;
export let tenant: string;

export const setPages = (value: AllPages) => {
  pages = value;
};

export const setBrowser = (value: Browser) => {
  browser = value;
};

export const setTenant = (value: string) => {
  tenant = value;
};
