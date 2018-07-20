export interface OAuthProviderOption {
  provider: string;
  displayName: string;
  clientId: string;
  url: string;
  callbackUrl: string;
  logoId: number;
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
  assetsURL: string;
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
  clientId: string;
  clientSecret: string;
  authorizeUrl: string;
  tokenUrl: string;
  profileUrl: string;
  logoId: number;
  scope: string;
  jsonUserIdPath: string;
  jsonUserNamePath: string;
  jsonUserEmailPath: string;
}
