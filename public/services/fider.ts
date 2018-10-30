import { CurrentUser, SystemSettings, Tenant } from "@fider/models";

export class FiderSession {
  private pContextID: string;
  private pTenant: Tenant;
  private pUser: CurrentUser | undefined;
  private pProps: { [key: string]: any } = {};

  constructor() {
    this.pContextID = window.__contextID;
    this.pProps = window.__props;
    this.pUser = window.__user;
    this.pTenant = window.__tenant;
  }

  public get contextID(): string {
    return this.pContextID;
  }

  public get user(): CurrentUser {
    return this.pUser!;
  }

  public get tenant(): Tenant {
    return this.pTenant;
  }

  public get props(): { [key: string]: any } {
    return this.pProps;
  }

  public get isAuthenticated(): boolean {
    return !!this.pUser;
  }
}

export class FiderImpl {
  private pSettings!: SystemSettings;
  private pSession!: FiderSession;

  public initialize = (): FiderImpl => {
    this.pSettings = window.__settings;
    this.pSession = new FiderSession();
    return this;
  };

  public get session(): FiderSession {
    return this.pSession;
  }

  public get settings(): SystemSettings {
    return this.pSettings;
  }

  public isProduction(): boolean {
    return this.pSettings.environment === "production";
  }

  public isSingleHostMode(): boolean {
    return this.pSettings.mode === "single";
  }
}

export let Fider = new FiderImpl();
