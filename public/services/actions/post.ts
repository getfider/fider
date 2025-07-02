import { http, Result, querystring } from "@fider/services"
import { Post, Vote, ImageUpload, UserNames } from "@fider/models"

export const getAllPosts = async (): Promise<Result<Post[]>> => {
  return await http.get<Post[]>("/api/v1/posts")
}

export interface SearchPostsParams {
  query?: string
  view?: string
  limit?: number
  tags?: string[]
  myVotes?: boolean
  statuses?: string[]
}

export const searchPosts = async (params: SearchPostsParams): Promise<Result<Post[]>> => {
  let qsParams = querystring.stringify({
    tags: params.tags,
    statuses: params.statuses,
    query: params.query,
    view: params.view,
    limit: params.limit,
  })
  if (params.myVotes) {
    qsParams += `&myvotes=true`
  }
  return await http.get<Post[]>(`/api/v1/posts${qsParams}`)
}

export const findSimilarPosts = async (query: string): Promise<Result<Post[]>> => {
  const params = querystring.stringify({ query: query })
  return await http.get<Post[]>(`/api/v1/similarposts${params}`)
}

export const deletePost = async (postNumber: number, text: string): Promise<Result> => {
  return http
    .delete(`/api/v1/posts/${postNumber}`, {
      text,
    })
    .then(http.event("post", "delete"))
}

export const addVote = async (postNumber: number): Promise<Result> => {
  return http.post(`/api/v1/posts/${postNumber}/votes`).then(http.event("post", "vote"))
}

export const removeVote = async (postNumber: number): Promise<Result> => {
  return http.delete(`/api/v1/posts/${postNumber}/votes`).then(http.event("post", "unvote"))
}

export const toggleVote = async (postNumber: number): Promise<Result<{ voted: boolean }>> => {
  return http.post<{ voted: boolean }>(`/api/v1/posts/${postNumber}/votes/toggle`).then(http.event("post", "toggle-vote"))
}

export const subscribe = async (postNumber: number): Promise<Result> => {
  return http.post(`/api/v1/posts/${postNumber}/subscription`).then(http.event("post", "subscribe"))
}

export const unsubscribe = async (postNumber: number): Promise<Result> => {
  return http.delete(`/api/v1/posts/${postNumber}/subscription`).then(http.event("post", "unsubscribe"))
}

export const listVotes = async (postNumber: number): Promise<Result<Vote[]>> => {
  return http.get<Vote[]>(`/api/v1/posts/${postNumber}/votes`)
}

export const getTaggableUsers = async (userFilter: string): Promise<Result<UserNames[]>> => {
  return http.get<UserNames[]>(`/api/v1/taggable-users${querystring.stringify({ query: userFilter })}`)
}

export const createComment = async (postNumber: number, content: string, attachments: ImageUpload[]): Promise<Result> => {
  return http.post(`/api/v1/posts/${postNumber}/comments`, { content, attachments }).then(http.event("comment", "create"))
}

export const updateComment = async (postNumber: number, commentID: number, content: string, attachments: ImageUpload[]): Promise<Result> => {
  return http.put(`/api/v1/posts/${postNumber}/comments/${commentID}`, { content, attachments }).then(http.event("comment", "update"))
}

export const deleteComment = async (postNumber: number, commentID: number): Promise<Result> => {
  return http.delete(`/api/v1/posts/${postNumber}/comments/${commentID}`).then(http.event("comment", "delete"))
}
interface ToggleReactionResponse {
  added: boolean
}

export const toggleCommentReaction = async (postNumber: number, commentID: number, emoji: string): Promise<Result<ToggleReactionResponse>> => {
  return http.post<ToggleReactionResponse>(`/api/v1/posts/${postNumber}/comments/${commentID}/reactions/${emoji}`)
}

interface SetResponseInput {
  status: string
  text: string
  originalNumber: number
}

export const respond = async (postNumber: number, input: SetResponseInput): Promise<Result> => {
  return http
    .put(`/api/v1/posts/${postNumber}/status`, {
      status: input.status,
      text: input.text,
      originalNumber: input.originalNumber,
    })
    .then(http.event("post", "respond"))
}

interface CreatePostResponse {
  id: number
  number: number
  title: string
  slug: string
}

export const createPost = async (title: string, description: string, attachments: ImageUpload[], tags: string[]): Promise<Result<CreatePostResponse>> => {
  return http.post<CreatePostResponse>(`/api/v1/posts`, { title, description, attachments, tags }).then(http.event("post", "create"))
}

export const updatePost = async (postNumber: number, title: string, description: string, attachments: ImageUpload[]): Promise<Result> => {
  return http.put(`/api/v1/posts/${postNumber}`, { title, description, attachments }).then(http.event("post", "update"))
}

export const approvePost = async (postID: number): Promise<Result> => {
  return http.post(`/api/v1/admin/moderation/posts/${postID}/approve`).then(http.event("post", "approve"))
}

export const declinePost = async (postID: number): Promise<Result> => {
  return http.post(`/api/v1/admin/moderation/posts/${postID}/decline`).then(http.event("post", "decline"))
}

export const approveComment = async (commentID: number): Promise<Result> => {
  return http.post(`/api/v1/admin/moderation/comments/${commentID}/approve`).then(http.event("comment", "approve"))
}

export const declineComment = async (commentID: number): Promise<Result> => {
  return http.post(`/api/v1/admin/moderation/comments/${commentID}/decline`).then(http.event("comment", "decline"))
}
