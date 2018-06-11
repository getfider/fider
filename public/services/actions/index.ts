export * from "./user";
export * from "./tag";
export * from "./idea";
export * from "./tenant";
export * from "./notification";
export * from "./invite";
export { Failure } from "@fider/services/http";

export const goHome = () => {
  window.location.href = "/";
};

export const replaceState = (path: string): void => {
  if (history.replaceState) {
    const newUrl = page.settings.baseURL + path;
    window.history.replaceState({ path: newUrl }, "", newUrl);
  }
};
