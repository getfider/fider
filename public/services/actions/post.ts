import { http, Result, querystring } from "@fider/services"
import { Post, Vote, ImageUpload } from "@fider/models"

export const getAllPosts = async (): Promise<Result<Post[]>> => {
  return await http.get<Post[]>("/api/v1/posts")
}

export interface SearchPostsParams {
  query?: string
  view?: string
  limit?: number
  tags?: string[]
}

export const searchPosts = async (params: SearchPostsParams): Promise<Result<Post[]>> => {
  return await http.get<Post[]>(
    `/api/v1/posts${querystring.stringify({
      tags: params.tags,
      query: params.query,
      view: params.view,
      limit: params.limit,
    })}`
  )
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

export const subscribe = async (postNumber: number): Promise<Result> => {
  return http.post(`/api/v1/posts/${postNumber}/subscription`).then(http.event("post", "subscribe"))
}

export const unsubscribe = async (postNumber: number): Promise<Result> => {
  return http.delete(`/api/v1/posts/${postNumber}/subscription`).then(http.event("post", "unsubscribe"))
}

export const listVotes = async (postNumber: number): Promise<Result<Vote[]>> => {
  return http.get<Vote[]>(`/api/v1/posts/${postNumber}/votes`)
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

export const createPost = async (title: string, description: string, attachments: ImageUpload[]): Promise<Result<CreatePostResponse>> => {
  return http.post<CreatePostResponse>(`/api/v1/posts`, { title, description, attachments }).then(http.event("post", "create"))
}

export const updatePost = async (postNumber: number, title: string, description: string, attachments: ImageUpload[]): Promise<Result> => {
  return http.put(`/api/v1/posts/${postNumber}`, { title, description, attachments }).then(http.event("post", "update"))
}
