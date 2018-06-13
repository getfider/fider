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
  auth: {
    endpoint: string;
    providers: {
      google: boolean;
      facebook: boolean;
      github: boolean;
    };
  };
}

export interface UserSettings {
  [key: string]: string;
}
