export interface Tenant {
  id: number;
  name: string;
  cname: string;
  subdomain: string;
  invitation: string;
  welcomeMessage: string;
  isPrivate: boolean;
  logoID: number;
}

export interface User {
  id: number;
  name: string;
  role: UserRole;
  status: UserStatus;
}

export enum UserStatus {
  Active = "active",
  Deleted = "deleted",
  Blocked = "blocked"
}

export enum UserRole {
  Visitor = "visitor",
  Collaborator = "collaborator",
  Administrator = "administrator"
}

export const isCollaborator = (role: UserRole): boolean => {
  return role === UserRole.Collaborator || role === UserRole.Administrator;
};

export interface CurrentUser {
  id: number;
  name: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  isAdministrator: boolean;
  isCollaborator: boolean;
}
