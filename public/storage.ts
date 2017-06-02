import { User } from './models';

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
