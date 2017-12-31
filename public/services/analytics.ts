export const analytics = {
  event: (evtName: string): void => {
    if (window.ga) {
      // TODO: implement correct API and use this on every action
      window.ga('send', evtName);
    }
  },
  error: (err: Error): void => {
    if (window.ga) {
      // TODO: send to GoogleAnalytics
    }
  }
};
