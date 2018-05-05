import { JSDOM } from "jsdom";

const jsdom = new JSDOM("<!doctype html><html><body></body></html>");
const { window } = jsdom;

function copyProps(src: any, target: any) {
  const props = Object.getOwnPropertyNames(src)
    .filter(prop => typeof target[prop] === "undefined")
    .reduce(
      (result, prop) => ({
        ...result,
        [prop]: Object.getOwnPropertyDescriptor(src, prop)
      }),
      {}
    );
  Object.defineProperties(target, props);
}

(global as any).window = window;
(global as any).document = window.document;
(global as any).navigator = {
  userAgent: "node.js"
};

let storage: { [key: string]: string | undefined } = {};
(window as any).sessionStorage = {
  getItem: (key: string) => {
    const value = storage[key];
    return typeof value === "undefined" ? null : value;
  },
  setItem: (key: string, value: string) => {
    storage[key] = value;
  },
  removeItem: (key: string) => {
    return delete storage[key];
  }
};

beforeEach(() => {
  storage = {};
});

copyProps(window, global);
