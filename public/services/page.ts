export const refresh = (): void => {
  document.location.reload();
};

export const setTitle = (title: string): void => {
  document.title = title;
};

export const getBaseUrl = (): string => {
  return (window as any).props.baseURL;
};

export const isSingleHostMode = (): boolean => {
  return (window as any).props.system.mode === 'single';
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

export interface QueryString {
  [key: string]: string | string[];
}

export const toQueryString = (object: QueryString): string => {
  if (!object) {
    return '';
  }

  let qs = '';

  for (const key of Object.keys(object)) {
    const symbol = qs ? '&' : '?';
    const value = object[key];
    if (value instanceof Array) {
      if (value.length > 0) {
        qs += `${symbol}${key}=${value.join(',')}`;
      }
    } else if (value) {
      qs += `${symbol}${key}=${encodeURIComponent(value).replace(/%20/g, '+')}`;
    }
  }

  return qs;
};

export const replaceState = (path: string): void => {
  if (history.replaceState) {
    const newUrl = getBaseUrl() + path;
    window.history.replaceState({ path: newUrl }, '', newUrl);
  }
};
