import { Fider } from "@fider/services";

const navigator = {
  url: () => {
    return window.location.href;
  },
  isBrowserSupported: () => {
    const ua = window.navigator.userAgent || window.navigator.appVersion;
    const isIE = ua.indexOf("MSIE") >= 0;
    if (isIE) {
      // On IE, the only supported version is IE 11
      return window.navigator.appVersion.indexOf("MSIE 11") >= 0;
    }
    return true;
  },
  goHome: () => {
    window.location.href = "/";
  },
  goTo: (url: string) => {
    const isEqual = window.location.href === url || window.location.pathname === url;
    if (!isEqual) {
      window.location.href = url;
    }
  },
  replaceState: (path: string): void => {
    if (history.replaceState) {
      const newURL = Fider.settings.baseURL + path;
      window.history.replaceState({ path: newURL }, "", newURL);
    }
  }
};

export default navigator;
