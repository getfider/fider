export function decode(token: string): any {
    const segments = token.split('.');
    return JSON.parse(window.atob(segments[1]));
}
