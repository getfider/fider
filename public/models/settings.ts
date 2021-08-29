export interface OAuthProviderOption {
  provider: string
  displayName: string
  clientID: string
  url: string
  callbackURL: string
  logoBlobKey: string
  isCustomProvider: boolean
  isEnabled: boolean
}

export interface SystemSettings {
  mode: string
  locale: string
  version: string
  environment: string
  domain: string
  hasLegal: boolean
  isBillingEnabled: boolean
  baseURL: string
  assetsURL: string
  oauth: OAuthProviderOption[]
}

export interface UserSettings {
  [key: string]: string
}

export const OAuthConfigStatus = {
  Disabled: 1,
  Enabled: 2,
}

export interface OAuthConfig {
  provider: string
  displayName: string
  status: number
  clientID: string
  clientSecret: string
  authorizeURL: string
  tokenURL: string
  profileURL: string
  logoBlobKey: string
  scope: string
  jsonUserIDPath: string
  jsonUserNamePath: string
  jsonUserEmailPath: string
}

export interface ImageUpload {
  bkey?: string
  upload?: {
    fileName?: string
    content?: string
    contentType?: string
  }
  remove: boolean
}
