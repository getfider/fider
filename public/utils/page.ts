export function setTitle(title: string) {
  document.title = title;
}

export function getBaseUrl(): string {
  return (window as any)._baseUrl;
}

export function showSignIn(): void {
  $('#signin-modal').modal({
    blurring: true
  }).modal('show');
}

export function hideSignIn(): void {
  $('#signin-modal').modal('hide');
}

export function getQueryString(name: string) {
  const url = window.location.href;
  name = name.replace(/[\[\]]/g, '\\$&');
  const regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)');
  const results = regex.exec(url);

  if (!results) {
    return null;
  }

  if (!results[2]) {
    return '';
  }

  return decodeURIComponent(results[2].replace(/\+/g, ' '));
}
