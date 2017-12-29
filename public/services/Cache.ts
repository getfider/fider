import { injectable } from '@fider/di';

export interface Cache {
  set(key: string, value: string): void;
  get(key: string): string | null;
  remove(...key: string[]): void;
}

@injectable()
export class BrowserCache implements Cache {
  private w: Window;

  constructor(window: Window) {
    this.w = window;
  }

  public set(key: string, value: string): void {
    if (this.w.sessionStorage) {
      this.w.sessionStorage.setItem(key, value);
    }
  }

  public get(key: string): string | null {
    if (this.w.sessionStorage) {
      return this.w.sessionStorage.getItem(key);
    }
    return null;
  }

  public remove(...keys: string[]): void {
    if (this.w.sessionStorage && keys) {
      for (const key of keys) {
        this.w.sessionStorage.removeItem(key);
      }
    }
  }
}
