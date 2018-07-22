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
      const newUrl = Fider.settings.baseURL + path;
      window.history.replaceState({ path: newUrl }, "", newUrl);
    }
  }
};

export default navigator;
