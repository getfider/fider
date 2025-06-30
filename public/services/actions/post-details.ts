import { http, Result } from "@fider/services"
import { Post, Comment, Tag, Vote } from "@fider/models"

export interface PostDetailsResponse {
  post: Post
  comments: Comment[]
  tags: Tag[]
  votes: Vote[]
  subscribed: boolean
  attachments: string[]
}

export const getPostDetails = async (postNumber: number): Promise<Result<PostDetailsResponse>> => {
  return await http.get<PostDetailsResponse>(`/api/v1/posts/${postNumber}/details`)
}
