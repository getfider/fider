export const setTitle = (title: string): void => {
  document.title = title;
};

export const getBaseUrl = (): string => {
  return (window as any).props.baseUrl;
};

export const isSingleHostMode = (): boolean => {
  return (window as any).props.settings.mode === 'single';
};

export interface ModalOptions {
  closable: boolean;
}

export const showModal = (selector: string, options?: ModalOptions): void => {
  const opts = Object.assign({ blurring: true }, options || { });
  $(selector).modal(opts).modal('show');
};

export const hideModal = (selector: string): void => {
  $(selector).modal('hide');
};

export const showSignIn = (): void => {
  showModal('#signin-modal');
};

export const hideSignIn = (): void => {
  hideModal('#signin-modal');
};

export const getQueryString = (name: string): string => {
  const url = window.location.href;
  name = name.replace(/[\[\]]/g, '\\$&');
  const regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)');
  const results = regex.exec(url);

  if (!results || !results[2]) {
    return '';
  }

  return decodeURIComponent(results[2].replace(/\+/g, ' '));
};

export const getQueryStringArray = (name: string): string[] => {
  const qs = getQueryString(name);
  if (qs) {
    return qs.split(',').filter((i) => i);
  }

  return [];
};
