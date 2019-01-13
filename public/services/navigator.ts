import { Fider, http } from "@fider/services";
import { cache } from "./cache";

const navigator = {
  userAgent: () => {
    return window.navigator.userAgent || window.navigator.appVersion;
  },
  url: () => {
    return window.location.href;
  },
  isBrowserSupported() {
    const ua = this.userAgent();
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
  },
  getCountryCode: (): Promise<string> => {
    const countryCode = cache.session.get("geolocation_countrycode");
    if (countryCode) {
      return Promise.resolve(countryCode);
    }

    return http.get<any>("https://ipinfo.io/geo").then(res => {
      if (res.ok) {
        cache.session.set("geolocation_countrycode", res.data.country);
        return res.data.country;
      }
    });
  }
};

export default navigator;
