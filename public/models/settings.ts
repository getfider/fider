export interface OAuthProviderOption {
  provider: string;
  displayName: string;
  clientID: string;
  url: string;
  callbackURL: string;
  logoID: number;
  isCustomProvider: boolean;
  isEnabled: boolean;
}

export interface SystemSettings {
  mode: string;
  buildTime: string;
  version: string;
  environment: string;
  compiler: string;
  domain: string;
  hasLegal: boolean;
  baseURL: string;
  tenantAssetsURL: string;
  globalAssetsURL: string;
  oauth: OAuthProviderOption[];
}

export interface UserSettings {
  [key: string]: string;
}

export const OAuthConfigStatus = {
  Disabled: 1,
  Enabled: 2
};

export interface OAuthConfig {
  provider: string;
  displayName: string;
  status: number;
  clientID: string;
  clientSecret: string;
  authorizeURL: string;
  tokenURL: string;
  profileURL: string;
  logoID: number;
  scope: string;
  jsonUserIDPath: string;
  jsonUserNamePath: string;
  jsonUserEmailPath: string;
}
