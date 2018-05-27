export interface Tenant {
  id: number;
  name: string;
  cname: string;
  subdomain: string;
  invitation: string;
  welcomeMessage: string;
  isPrivate: boolean;
  logoId: number;
}

export interface User {
  id: number;
  name: string;
  role: UserRole;
  status: UserStatus;
}

export enum UserStatus {
  Active = 1,
  Deleted = 2
}

export enum UserRole {
  Visitor = 1,
  Collaborator = 2,
  Administrator = 3
}

export interface CurrentUser {
  id: number;
  name: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  isAdministrator: boolean;
  isCollaborator: boolean;
}
