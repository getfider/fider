import { User } from "./identity";

export interface Post {
  id: number;
  number: number;
  slug: string;
  title: string;
  description: string;
  createdOn: string;
  status: number;
  user: User;
  viewerSupported: boolean;
  response: PostResponse | null;
  totalSupporters: number;
  totalComments: number;
  tags: string[];
}

export class PostStatus {
  constructor(
    public value: number,
    public title: string,
    public slug: string,
    public show: boolean,
    public closed: boolean,
    public filterable: boolean
  ) {}

  public static Open = new PostStatus(0, "Open", "open", false, false, false);
  public static Planned = new PostStatus(4, "Planned", "planned", true, false, true);
  public static Started = new PostStatus(1, "Started", "started", true, false, true);
  public static Completed = new PostStatus(2, "Completed", "completed", true, true, true);
  public static Declined = new PostStatus(3, "Declined", "declined", true, true, true);
  public static Duplicate = new PostStatus(5, "Duplicate", "duplicate", true, true, false);
  public static Deleted = new PostStatus(6, "Deleted", "deleted", false, true, false);

  public static Get(value: number): PostStatus {
    for (const status of PostStatus.All) {
      if (status.value === value) {
        return status;
      }
    }
    throw new Error(`PostStatus not found for value ${value}.`);
  }

  public static All = [
    PostStatus.Open,
    PostStatus.Planned,
    PostStatus.Started,
    PostStatus.Completed,
    PostStatus.Duplicate,
    PostStatus.Declined
  ];
}

export interface PostResponse {
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
