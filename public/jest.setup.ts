import { configure } from "enzyme";
import * as Adapter from "enzyme-adapter-react-16";

configure({ adapter: new Adapter() });

let storage: {
  [key: string]: string | undefined;
};

beforeEach(() => {
  storage = {};
});

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
