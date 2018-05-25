export interface SystemSettings {
  mode: string;
  buildTime: string;
  version: string;
  authEndpoint: string;
  environment: string;
  googleAnalytics: string;
  compiler: string;
  domain: string;
  hasLegal: boolean;
}

export interface AuthSettings {
  endpoint: string;
  providers: {
    google: boolean;
    facebook: boolean;
    github: boolean;
  };
}

export interface UserSettings {
  [key: string]: string;
}
