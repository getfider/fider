import { Fider } from "@fider/services";

const navigator = {
  url: () => {
    return window.location.href;
  },
  goHome: () => {
    window.location.href = "/";
  },
  goTo: (url: string) => {
    window.location.href = url;
  },
  replaceState: (path: string): void => {
    if (history.replaceState) {
      const newURL = Fider.settings.baseURL + path;
      window.history.replaceState({ path: newURL }, "", newURL);
    }
  }
};

export default navigator;
