export interface OAuthProviderOption {
  provider: string;
  displayName: string;
  url: string;
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

export interface OAuthConfig {
  id: number;
  provider: string;
  displayName: string;
  status: number;
  clientId: string;
  authorizeUrl: string;
  tokenUrl: string;
  profileUrl: string;
  scope: string;
  jsonUserIdPath: string;
  jsonUserNamePath: string;
  jsonUserEmailPath: string;
}
