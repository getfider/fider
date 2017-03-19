export interface Tenant {
  id: number;
  name: string;
  domain: string;
}

export interface User {
  id: number;
  name: string;
  email: string;
}

export interface Idea {
  id: number;
  title: string;
  description: string;
  createdOn: string;
}
