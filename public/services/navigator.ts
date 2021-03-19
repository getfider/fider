import { Fider, http } from "@fider/services";
import { cache } from "./cache";

const navigator = {
  url: () => {
    return window.location.href;
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
    if (history.replaceState !== undefined) {
      const newURL = Fider.settings.baseURL + path;
      window.history.replaceState({ path: newURL }, "", newURL);
    }
  },
  getCountryCode: (): Promise<string> => {
    const countryCode = cache.session.get("geolocation_countrycode");
    if (countryCode) {
      return Promise.resolve(countryCode);
    }

    return http.get<any>("https://ipinfo.io/geo").then((res) => {
      if (res.ok) {
        cache.session.set("geolocation_countrycode", res.data.country);
        return res.data.country;
      }
    });
  },
};

export default navigator;
