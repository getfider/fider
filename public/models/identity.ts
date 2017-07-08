export interface Tenant {
  id: number;
  name: string;
  domain: string;
}

export interface User {
  id: number;
  name: string;
  gravatar: string;
  role: number;
}
