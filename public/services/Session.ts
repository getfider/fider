import { User, AppSettings } from '../models';
import { injectable } from '../di';

export interface Session {
    getCurrentUser(): User;
    setUser(user: User): void;
    isStaff(): boolean;
    set<T>(key: string, value: T): void;
    get<T>(key: string): T;
    getArray<T>(key: string): T[];
    getAppSettings(): AppSettings;
    isSingleHostMode(): boolean;
    isProduction(): boolean;
}

@injectable()
export class BrowserSession implements Session {
    private w: any;

    constructor(window: Window) {
        this.w = window;
        this.w.setUser = this.setUser.bind(this);
        this.w.getCurrentUser = this.getCurrentUser.bind(this);
        this.w.set = this.set.bind(this);
        this.w.get = this.get.bind(this);
    }

    public getCurrentUser(): User {
        const w: any = window;
        return w[`_user`] as User;
    }

    public setUser(user: User): void {
        this.set<User>('user', user);
    }

    public isStaff(): boolean {
        const user = this.getCurrentUser();
        return user && user.role >= 2;
    }

    public set<T>(key: string, value: T): void {
        this.w[`_${key}`] = value;
    }

    public get<T>(key: string): T {
        return this.w[`_${key}`];
    }

    public getArray<T>(key: string): T[] {
        return this.w[`_${key}`] || [];
    }

    public getAppSettings(): AppSettings {
        return this.get<AppSettings>('settings');
    }

    public isSingleHostMode(): boolean {
        return this.getAppSettings().mode.toLowerCase() === 'single';
    }

    public isProduction(): boolean {
        return this.getAppSettings().environment.toLowerCase() === 'production';
    }
}
