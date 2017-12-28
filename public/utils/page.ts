export function setTitle(title: string) {
  document.title = title;
}

export function getBaseUrl(): string {
  return (window as any).props.baseUrl;
}

export interface ModalOptions {
  closable: boolean;
}

export function showModal(selector: string, options?: ModalOptions): void {
  const opts = Object.assign({ blurring: true }, options || { });

  $(selector).modal(opts).modal('show');
}

export function hideModal(selector: string): void {
  $(selector).modal('hide');
}

export function showSignIn(): void {
  showModal('#signin-modal');
}

export function hideSignIn(): void {
  hideModal('#signin-modal');
}

export function getQueryString(name: string): string {
  const url = window.location.href;
  name = name.replace(/[\[\]]/g, '\\$&');
  const regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)');
  const results = regex.exec(url);

  if (!results || !results[2]) {
    return '';
  }

  return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

export function getQueryStringArray(name: string): string[] {
  const qs = getQueryString(name);
  if (qs) {
    return qs.split(',').filter((i) => i);
  }

  return [];
}
