import { Browser } from "./lib";
import { AllPages } from "./pages";

let pages: AllPages;
let browser: Browser;
let tenant: string;

export const setPages = (value: AllPages) => {
  pages = value;
};

export const setBrowser = (value: Browser) => {
  browser = value;
};

export const setTenant = (value: string) => {
  if (process.env.HOST_MODE === "single") {
    tenant = "login";
  } else {
    tenant = value;
  }
};

export const getTenant = (): string => {
  return tenant;
};
