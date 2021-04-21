import { http, Result } from "@fider/services/http"
import { Tag } from "@fider/models"

export const createTag = async (name: string, color: string, isPublic: boolean): Promise<Result<Tag>> => {
  return http
    .post<Tag>(`/api/v1/tags`, { name, color, isPublic })
    .then(http.event("tag", "create"))
}

export const updateTag = async (slug: string, name: string, color: string, isPublic: boolean): Promise<Result<Tag>> => {
  return http
    .put<Tag>(`/api/v1/tags/${slug}`, { name, color, isPublic })
    .then(http.event("tag", "update"))
}

export const deleteTag = async (slug: string): Promise<Result> => {
  return http.delete(`/api/v1/tags/${slug}`).then(http.event("tag", "delete"))
}

export const assignTag = async (slug: string, postNumber: number): Promise<Result> => {
  return http.post(`/api/v1/posts/${postNumber}/tags/${slug}`).then(http.event("tag", "assign"))
}

export const unassignTag = async (slug: string, postNumber: number): Promise<Result> => {
  return http.delete(`/api/v1/posts/${postNumber}/tags/${slug}`).then(http.event("tag", "unassign"))
}
