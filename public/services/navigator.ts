const navigator = {
  url: () => {
    return window.location.href;
  },
  goHome: () => {
    window.location.href = "/";
  },
  replaceState: (path: string): void => {
    if (history.replaceState) {
      const newUrl = Fider.settings.baseURL + path;
      window.history.replaceState({ path: newUrl }, "", newUrl);
    }
  }
};

export default navigator;
