import { http } from "../http"

export async function updateProfanityWords(profanityWords: string) {
  const result = await http.post("/api/admin/profanity-words", { profanityWords })
  return result
}