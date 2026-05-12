import { http, Result } from "@fider/services"
import { ImageUpload } from "@fider/models"

export interface UploadImageResponse {
  bkey: string
}

/**
 * Uploads an image to the server and returns the bkey
 * This doesn't associate the image with any post or comment yet
 */
export const uploadImage = async (image: ImageUpload): Promise<Result<UploadImageResponse>> => {
  return http.post<UploadImageResponse>("/api/v1/images", { image })
}

export const deleteImage = async (bkey: string): Promise<Result> => {
  return http.delete(`/api/v1/images/${bkey}`)
}
