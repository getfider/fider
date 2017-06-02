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

export const IdeaNew = 0;
export const IdeaStarted = 1;
export const IdeaCompleted = 2;
export const IdeaDeclined = 3;

export const IdeaStatusMetadata: { [key: number]: any} = { };

IdeaStatusMetadata[IdeaNew] = {
  title: "New",
  showStatus: false,
  closed: false,
  color: "black",
};

IdeaStatusMetadata[IdeaStarted] = {
  title: "Started",
  showStatus: true,
  closed: false,
  color: "blue",
};

IdeaStatusMetadata[IdeaCompleted] = {
  title: "Completed",
  showStatus: true,
  closed: true,
  color: "green",
};

IdeaStatusMetadata[IdeaDeclined] = {
  title: "Declined",
  showStatus: true,
  closed: true,
  color: "red",
};

export interface IdeaResponse {
  user: User;
  text: string;
  respondedOn: Date;
}

export interface Comment {
  id: number;
  content: string;
  createdOn: string;
  user: User;
}
