import { User } from "./identity"

export interface Post {
  id: number
  number: number
  slug: string
  title: string
  description: string
  createdAt: string
  status: string
  user: User
  hasVoted: boolean
  response: PostResponse | null
  votesCount: number
  commentsCount: number
  tags: string[]
}

export class PostStatus {
  constructor(public title: string, public value: string, public show: boolean, public closed: boolean, public filterable: boolean) {}

  public static Open = new PostStatus("Open", "open", false, false, false)
  public static Planned = new PostStatus("Planned", "planned", true, false, true)
  public static Started = new PostStatus("Started", "started", true, false, true)
  public static Completed = new PostStatus("Completed", "completed", true, true, true)
  public static Declined = new PostStatus("Declined", "declined", true, true, true)
  public static Duplicate = new PostStatus("Duplicate", "duplicate", true, true, false)
  public static Deleted = new PostStatus("Deleted", "deleted", false, true, false)

  public static Get(value: string): PostStatus {
    for (const status of PostStatus.All) {
      if (status.value === value) {
        return status
      }
    }
    throw new Error(`PostStatus not found for value ${value}.`)
  }

  public static All = [PostStatus.Open, PostStatus.Planned, PostStatus.Started, PostStatus.Completed, PostStatus.Duplicate, PostStatus.Declined]
}

export interface PostResponse {
  user: User
  text: string
  respondedAt: Date
  original?: {
    number: number
    title: string
    slug: string
    status: string
  }
}

export interface Comment {
  id: number
  content: string
  createdAt: string
  user: User
  attachments?: string[]
  editedAt?: string
  editedBy?: User
}

export interface Tag {
  id: number
  slug: string
  name: string
  color: string
  isPublic: boolean
}

export interface Vote {
  createdAt: Date
  user: {
    id: number
    name: string
    email: string
    avatarURL: string
  }
}
