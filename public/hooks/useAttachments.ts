import { useState, useEffect, useCallback } from "react"
import { ImageUpload } from "@fider/models"
import { cache } from "@fider/services"

interface UseAttachmentsOptions {
  cacheKey?: string
  maxAttachments?: number
}

interface UseAttachmentsReturn {
  attachments: ImageUpload[]
  setAttachments: React.Dispatch<React.SetStateAction<ImageUpload[]>>
  handleImageUploaded: (upload: ImageUpload) => void
  getImageSrc: (bkey: string) => string
  clearAttachments: () => void
}

export const useAttachments = (options: UseAttachmentsOptions = {}): UseAttachmentsReturn => {
  const { cacheKey, maxAttachments } = options

  const getAttachmentsFromCache = useCallback(() => {
    if (!cacheKey) return []
    const cachedAttachments = cache.session.get(cacheKey)
    return cachedAttachments ? JSON.parse(cachedAttachments) : []
  }, [cacheKey])

  const [attachments, setAttachments] = useState<ImageUpload[]>(getAttachmentsFromCache())

  // Cache attachments whenever they change
  useEffect(() => {
    if (cacheKey) {
      cache.session.set(cacheKey, JSON.stringify(attachments))
    }
  }, [attachments, cacheKey])

  const handleImageUploaded = useCallback(
    (upload: ImageUpload) => {
      setAttachments((prev) => {
        // If this is a removal request, find and remove the attachment
        // if (upload.remove && upload.bkey) {
        //   return prev.filter((att) => att.bkey !== upload.bkey)
        // }
        if (upload.remove) {
          // Rules are:
          // 1.If the bkey exists in attachments:
          // 1.1 If there is something in upload.content, then it's a new upload that's not yet saved, so just remove the entry from attachments.
          // 1.2 If there is nothing in upload.content, it's been previously uploaded, so set the remove flag to true.
          // 2. If the bkey doesn't exist, then you're editing, and attachments hasn't been "primed" so add it as a new one but with with remove flag set to true

          const existing = prev.find((att) => att.bkey === upload.bkey)

          if (existing) {
            if (existing.upload?.content && existing.upload.content.length > 0) {
              return prev.filter((att) => att.bkey !== upload.bkey)
            } else {
              return prev.map((att) => (att.bkey === upload.bkey ? { ...att, remove: true } : att))
            }
          } else {
            return [...prev, { ...upload, remove: true }]
          }
        }

        // Check max attachments limit
        if (maxAttachments && prev.length >= maxAttachments) {
          console.warn(`Maximum ${maxAttachments} attachments allowed`)
          return prev
        }

        // Otherwise add the new upload
        return [...prev, upload]
      })
    },
    [maxAttachments]
  )

  const getImageSrc = useCallback(
    (bkey: string): string => {
      return attachments.find((att) => att.bkey === bkey)?.upload?.content?.toString() ?? ""
    },
    [attachments]
  )

  const clearAttachments = useCallback(() => {
    setAttachments([])
    if (cacheKey) {
      cache.session.remove(cacheKey)
    }
  }, [cacheKey])

  return {
    attachments,
    setAttachments,
    handleImageUploaded,
    getImageSrc,
    clearAttachments,
  }
}
