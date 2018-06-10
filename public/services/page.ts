import { SystemSettings } from "@fider/models";

export const goHome = (): void => {
  document.location.href = "/";
};

export const refresh = (): void => {
  document.location.reload();
};

export const isSingleHostMode = (): boolean => {
  return window.props.settings.mode === "single";
};

export interface ModalOptions {
  closable: boolean;
}

export const getQueryStringAsNumber = (name: string): number | undefined => {
  return parseInt(getQueryString(name), 10) || undefined;
};

export const getQueryString = (name: string): string => {
  const url = window.location.href;
  name = name.replace(/[\[\]]/g, "\\$&");
  const regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)");
  const results = regex.exec(url);

  if (!results || !results[2]) {
    return "";
  }

  return decodeURIComponent(results[2].replace(/\+/g, " "));
};

export const getQueryStringArray = (name: string): string[] => {
  const qs = getQueryString(name);
  if (qs) {
    return qs.split(",").filter(i => i);
  }

  return [];
};

export interface QueryString {
  [key: string]: string | string[] | number | undefined;
}

export const toQueryString = (object: QueryString): string => {
  if (!object) {
    return "";
  }

  let qs = "";

  for (const key of Object.keys(object)) {
    const symbol = qs ? "&" : "?";
    const value = object[key];
    if (value instanceof Array) {
      if (value.length > 0) {
        qs += `${symbol}${key}=${value.join(",")}`;
      }
    } else if (value) {
      qs += `${symbol}${key}=${encodeURIComponent(value.toString()).replace(/%20/g, "+")}`;
    }
  }

  return qs;
};

export const replaceState = (path: string): void => {
  if (history.replaceState) {
    const newUrl = window.props.settings.baseURL + path;
    window.history.replaceState({ path: newUrl }, "", newUrl);
  }
};
