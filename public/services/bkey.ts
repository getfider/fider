/**
 * Converts a string to a URL-friendly slug
 * @param text The text to convert to a slug
 * @returns A URL-friendly slug
 */
export const slugify = (text: string): string => {
  return text
    .toString()
    .toLowerCase()
    .trim()
    .replace(/\s+/g, "-") // Replace spaces with -
    .replace(/[^\w-]+/g, "") // Remove all non-word chars
    .replace(/--+/g, "-") // Replace multiple - with single -
    .replace(/^-+/, "") // Trim - from start of text
    .replace(/-+$/, "") // Trim - from end of text
}

/**
 * Generates a random string of the specified length
 * @param length The length of the random string to generate
 * @returns A random string
 */
export const generateRandomString = (length: number): string => {
  const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  let result = ""
  const randomValues = new Uint8Array(length)

  // Use crypto.getRandomValues for secure random generation if available
  if (window.crypto && window.crypto.getRandomValues) {
    window.crypto.getRandomValues(randomValues)

    for (let i = 0; i < length; i++) {
      result += chars[randomValues[i] % chars.length]
    }
  } else {
    // Fallback to Math.random (less secure)
    for (let i = 0; i < length; i++) {
      result += chars[Math.floor(Math.random() * chars.length)]
    }
  }

  return result
}

/**
 * Sanitizes a filename to be used in a URL
 * @param fileName The filename to sanitize
 * @returns A sanitized filename
 */
export const sanitizeFileName = (fileName: string): string => {
  fileName = fileName.trim()
  const lastDotIndex = fileName.lastIndexOf(".")

  if (lastDotIndex !== -1) {
    const name = fileName.substring(0, lastDotIndex)
    const ext = fileName.substring(lastDotIndex)
    return slugify(name) + ext
  }

  return slugify(fileName)
}

/**
 * Generates a bkey for an image upload
 * @param fileName The name of the file being uploaded
 * @returns A bkey in the format used by the server
 */
export const generateBkey = (fileName: string): string => {
  const folder = "attachments"
  const randomString = generateRandomString(64)
  const sanitizedFileName = sanitizeFileName(fileName)

  return `${folder}/${randomString}-${sanitizedFileName}`
}
