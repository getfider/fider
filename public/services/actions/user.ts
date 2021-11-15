import { http, Result } from "@fider/services/http"
import { UserSettings, UserAvatarType, ImageUpload } from "@fider/models"

interface UpdateUserSettings {
  name: string
  avatar?: ImageUpload
  avatarType: UserAvatarType
  settings: UserSettings
}

export const updateUserSettings = async (request: UpdateUserSettings): Promise<Result> => {
  return await http.post("/_api/user/settings", request)
}

export const changeUserEmail = async (email: string): Promise<Result> => {
  return await http.post("/_api/user/change-email", {
    email,
  })
}

export const deleteCurrentAccount = async (): Promise<Result> => {
  return await http.delete("/_api/user")
}

export const regenerateAPIKey = async (): Promise<Result<{ apiKey: string }>> => {
  return await http.post<{ apiKey: string }>("/_api/user/regenerate-apikey")
}
