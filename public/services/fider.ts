import { createContext } from "react"
import { CurrentUser, SystemSettings, Tenant, TenantStatus } from "@fider/models"

export class FiderSession {
  private pPage: string
  private pContextID: string
  private pTenant: Tenant
  private pUser: CurrentUser | undefined
  private pProps: { [key: string]: any } = {}

  constructor(data: any) {
    this.pPage = data.page
    this.pContextID = data.contextID
    this.pProps = data.props
    this.pUser = data.user
    this.pTenant = data.tenant
  }

  public get page(): string {
    return this.pPage
  }

  public get contextID(): string {
    return this.pContextID
  }

  public get user(): CurrentUser {
    if (!this.pUser) throw new Error("User is undefined")
    return this.pUser
  }

  public get tenant(): Tenant {
    return this.pTenant
  }

  public get props(): { [key: string]: any } {
    return this.pProps
  }

  public get isAuthenticated(): boolean {
    return !!this.pUser
  }
}

export class FiderImpl {
  private pSettings!: SystemSettings
  private pSession!: FiderSession

  public initialize = (initData?: any): FiderImpl => {
    if (initData) {
      this.pSettings = initData.settings
      this.pSession = new FiderSession(initData)
      return this
    }

    const el = document.getElementById("server-data")
    const data = el ? JSON.parse(el.textContent || el.innerText) : {}
    this.pSettings = data.settings
    this.pSession = new FiderSession(data)
    return this
  }

  public get currentLocale(): string {
    if (this.session.tenant) {
      return this.session.tenant.locale
    }
    return this.settings.locale
  }

  public get session(): FiderSession {
    return this.pSession
  }

  public get settings(): SystemSettings {
    return this.pSettings
  }

  public get isReadOnly(): boolean {
    return this.session.tenant && this.session.tenant.status === TenantStatus.Locked
  }

  public isProduction(): boolean {
    return this.pSettings.environment === "production"
  }

  public isSingleHostMode(): boolean {
    return this.pSettings.mode === "single"
  }
}

export const Fider = new FiderImpl()

export const FiderContext = createContext<FiderImpl>(Fider)
