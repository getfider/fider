import { CurrentUser, SystemSettings, Tenant } from ".";

export class FiderPage {
  public user?: CurrentUser;
  public tenant!: Tenant;
  public settings!: SystemSettings;
  public props: { [key: string]: any } = {};

  public set(key: string, value: any) {
    if (key === "user") {
      this.user = value;
    } else if (key === "tenant") {
      this.tenant = value;
    } else if (key === "settings") {
      this.settings = value;
    } else {
      this.props[key] = value;
    }
  }

  public isProduction(): boolean {
    return this.settings.environment === "production";
  }

  public isSingleHostMode(): boolean {
    return this.settings.mode === "single";
  }
}
