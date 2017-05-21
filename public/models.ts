export interface Tenant {
  id: number;
  name: string;
  domain: string;
}

export interface User {
  id: number;
  name: string;
  email: string;
  supportedIdeas: number[];
  role: number;
}

export interface Idea {
  id: number;
  number: number;
  slug: string;
  title: string;
  description: string;
  createdOn: string;
  status: number;
  user: User;
  response: IdeaResponse;
  totalSupporters: number;
}

export interface IdeaResponse {
  userId: number;
  text: string;
  createdOn: Date;
}

export interface Comment {
  id: number;
  content: string;
  createdOn: string;
  user: User;
}
