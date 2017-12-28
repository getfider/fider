import { User, CurrentUser, AppSettings, Tenant } from '@fider/models';
import { injectable } from '@fider/di';

export interface Session {
  props(): any;
  set<T>(key: string, value: T): void;
  get<T>(key: string): T;
  getAppSettings(): AppSettings;
  isSingleHostMode(): boolean;
  setCache(key: string, value: string): void;
  getCache(key: string): string | undefined;
  removeCache(...key: string[]): void;
}

@injectable()
export class BrowserSession implements Session {
  private w: any;

  constructor(window: Window) {
    this.w = window;
    this.w.props = {};
    this.w.set = this.set.bind(this);
    this.w.get = this.get.bind(this);
  }

  public props(): any {
    return this.w.props;
  }

  public set<T>(key: string, value: T): void {
    this.w.props[`${key}`] = value;
  }

  public get<T>(key: string): T {
    return this.w.props[`${key}`];
  }

  public getAppSettings(): AppSettings {
    return this.get<AppSettings>('settings');
  }

  public isSingleHostMode(): boolean {
    return this.getAppSettings().mode.toLowerCase() === 'single';
  }

  public setCache(key: string, value: string): void {
    if (this.w.sessionStorage) {
      this.w.sessionStorage.setItem(key, value);
    }
  }

  public getCache(key: string): string | undefined {
    if (this.w.sessionStorage) {
      return this.w.sessionStorage.getItem(key);
    }
    return undefined;
  }

  public removeCache(...keys: string[]): void {
    if (this.w.sessionStorage && keys) {
      for (const key of keys) {
        this.w.sessionStorage.removeItem(key);
      }
    }
  }
}
