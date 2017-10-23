export interface Tenant {
  id: number;
  name: string;
  domain: string;
  invitation: string;
  welcomeMessage: string;
}

export interface User {
  id: number;
  name: string;
  gravatar: string;
  role: number;
}

export enum UserRole {
  Visitor = 1,
  Collaborator = 2,
  Administrator = 3,
}

export interface CurrentUser {
  id: number;
  name: string;
  email: string;
  gravatar: string;
  role: number;
}
