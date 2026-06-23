import { User } from "./identity"

// Status mirrors entity.Status from the Go server. Provided per-tenant on
// session bootstrap as fider.session.tenant.statuses. The admin UI lets
// tenants extend this catalogue beyond the built-in 6 (feedback.fider.io/111).
export type StatusKind = "open" | "active" | "closed-completed" | "closed-declined" | "duplicate"

export interface Status {
  id: number
  slug: string
  label: string
  kind: StatusKind
  color: string
  icon: string
  showOnHome: boolean
  filterable: boolean
  sortOrder: number
  isSystem: boolean
  isActive: boolean
}

export interface Post {
  id: number
  number: number
  slug: string
  title: string
  description: string
  createdAt: string
  // Tenant-defined slug. Built-in or admin-added; matches a row in
  // tenant.statuses for color/label resolution.
  status: string
  // Semantic kind joined from the statuses table. Empty for the special
  // "deleted" sentinel which has no row in the catalogue.
  statusKind?: string
  user: User
  hasVoted: boolean
  response: PostResponse | null
  votesCount: number
  commentsCount: number
  tags: string[]
  isApproved: boolean
}

// Returns the effective status slug for a post. Kept as a small helper so
// callers stay forward-compatible if the shape evolves again.
export const postStatusValue = (p: Pick<Post, "status">): string => p.status

export class PostStatus {
  constructor(public title: string, public value: string, public show: boolean, public closed: boolean, public filterable: boolean) {}

  public static Open = new PostStatus("Open", "open", false, false, true)
  public static Planned = new PostStatus("Planned", "planned", true, false, true)
  public static Started = new PostStatus("Started", "started", true, false, true)
  public static Completed = new PostStatus("Completed", "completed", true, true, true)
  public static Declined = new PostStatus("Declined", "declined", true, true, true)
  public static Duplicate = new PostStatus("Duplicate", "duplicate", true, true, true)
  public static Deleted = new PostStatus("Deleted", "deleted", false, true, true)

  public static Get(value: string): PostStatus {
    for (const status of PostStatus.All) {
      if (status.value === value) {
        return status
      }
    }
    // Tenant-defined custom statuses (feedback.fider.io/111) won't have a
    // matching constant. Return a synthetic so the post page can still render;
    // callers that need richer info (color/icon/label) should look up the
    // status in fider.session.tenant.statuses via the helpers in
    // status-helpers.ts.
    const synthesizedTitle = value.charAt(0).toUpperCase() + value.slice(1).replace(/-/g, " ")
    return new PostStatus(synthesizedTitle, value, true, false, true)
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

export interface ReactionCount {
  emoji: string
  count: number
  includesMe: boolean
}

export interface Comment {
  id: number
  content: string
  createdAt: string
  user: User
  attachments?: string[]
  reactionCounts?: ReactionCount[]
  editedAt?: string
  editedBy?: User
  isApproved: boolean
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

export interface InlineImage {
  bkey: string
  remove: boolean
}
