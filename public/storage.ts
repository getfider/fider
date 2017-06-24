import { User, AppSettings } from './models';

const w: any = window;

export function setup() {
  w.setUser = setUser;
  w.getCurrentUser = getCurrentUser;
  w.set = set;
  w.get = get;
}

export function setUser(user: User) {
  set<User>('user', user);
}

export function getCurrentUser(): User {
  return get<User>('user');
}

export function isStaff() {
  const user = getCurrentUser();
  return user && user.role >= 2;
}

export function set<T>(key: string, value: T) {
  w[`_${key}`] = value;
}

export function get<T>(key: string): T {
  return w[`_${key}`];
}

export function getArray<T>(key: string): T[] {
  return w[`_${key}`] || [];
}

export function getAppSettings(): AppSettings {
  return get<AppSettings>('settings');
}

export function isSingleHostMode() {
  return getAppSettings().mode.toLowerCase() === 'single';
}

export function isProduction() {
  return getAppSettings().environment.toLowerCase() === 'production';
}
