import { http, Result } from "@fider/services/http"
import { Status, StatusKind } from "@fider/models"

export interface CreateStatusInput {
  slug: string
  label: string
  kind: StatusKind
  color: string
  icon: string
  showOnHome: boolean
  filterable: boolean
  sortOrder: number
}

export interface UpdateStatusInput {
  label: string
  color: string
  icon: string
  showOnHome: boolean
  filterable: boolean
  sortOrder: number
  isActive: boolean
}

export const listStatuses = async (): Promise<Result<Status[]>> => {
  return http.get<Status[]>(`/_api/admin/statuses`)
}

export const createStatus = async (input: CreateStatusInput): Promise<Result<Status>> => {
  return http.post<Status>(`/_api/admin/statuses`, input).then(http.event("status", "create"))
}

export const updateStatus = async (id: number, input: UpdateStatusInput): Promise<Result> => {
  return http.put(`/_api/admin/statuses/${id}`, input).then(http.event("status", "update"))
}

export const deleteStatus = async (id: number): Promise<Result> => {
  return http.delete(`/_api/admin/statuses/${id}`).then(http.event("status", "delete"))
}
