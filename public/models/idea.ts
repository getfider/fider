import { User } from "./identity";

export interface Idea {
  id: number;
  number: number;
  slug: string;
  title: string;
  description: string;
  createdOn: string;
  status: number;
  user: User;
  viewerSupported: boolean;
  response: IdeaResponse;
  totalSupporters: number;
  totalComments: number;
  tags: string[];
}

export class IdeaStatus {
  constructor(
    public value: number,
    public title: string,
    public slug: string,
    public show: boolean,
    public closed: boolean,
    public filterable: boolean,
    public color: string
  ) {}

  public static Open = new IdeaStatus(0, "Open", "open", false, false, false, "");
  public static Planned = new IdeaStatus(4, "Planned", "planned", true, false, true, "violet");
  public static Started = new IdeaStatus(1, "Started", "started", true, false, true, "blue");
  public static Completed = new IdeaStatus(2, "Completed", "completed", true, true, true, "green");
  public static Declined = new IdeaStatus(3, "Declined", "declined", true, true, true, "red");
  public static Duplicate = new IdeaStatus(5, "Duplicate", "duplicate", true, true, false, "yellow");

  public static Get(value: number): IdeaStatus {
    for (const status of IdeaStatus.All) {
      if (status.value === value) {
        return status;
      }
    }
    throw new Error(`IdeaStatus not found for value ${value}.`);
  }

  public static All = [
    IdeaStatus.Open,
    IdeaStatus.Planned,
    IdeaStatus.Started,
    IdeaStatus.Completed,
    IdeaStatus.Duplicate,
    IdeaStatus.Declined
  ];
}

export interface IdeaResponse {
  user: User;
  text: string;
  respondedOn: Date;
  original?: {
    number: number;
    title: string;
    slug: string;
    status: number;
  };
}

export interface Comment {
  id: number;
  content: string;
  createdOn: string;
  user: User;
  editedOn?: string;
  editedBy?: User;
}

export interface Tag {
  id: number;
  slug: string;
  name: string;
  color: string;
  isPublic: boolean;
}
