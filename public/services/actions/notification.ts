import { http, Result } from "@fider/services"
import { Notification } from "@fider/models"

export const getTotalUnreadNotifications = async (): Promise<Result<number>> => {
  return http.get<{ total: number }>("/_api/notifications/unread/total").then((result) => {
    return {
      ok: result.ok,
      error: result.error,
      data: result.data ? result.data.total : 0,
    }
  })
}

export const getAllNotifications = async (): Promise<Result<Notification[]>> => {
  return http.get<Notification[]>("/_api/notifications/unread")
}

export const markAllAsRead = async (): Promise<Result> => {
  return await http.post("/_api/notifications/read-all")
}
