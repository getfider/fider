/**
 * Normalizes a URL, adding protocol if needed for regular URLs
 * @param url The URL to normalize
 * @returns The normalized URL
 */
export const normalizeUrl = (url: string) => {
  // Parse the scheme from the URL (e.g. https, mailto, monero, bitcoin, etc.)
  const trimmedUrl = url.trim()
  const scheme = trimmedUrl.match(/^([a-zA-Z][a-zA-Z0-9+.-]*):/)?.[1]?.toLowerCase()
  // For URLs without a scheme, add https://, otherwise return the URL as is
  return scheme ? trimmedUrl : `https://${trimmedUrl}`
}

/**
 * URL validation with allowed schemes
 * @param url The URL to validate
 * @param allowedSchemesRegex Array of RegExp patterns for allowed schemes
 * @returns true if valid, false otherwise
 */
export const isValidUrl = (url: string, allowedSchemesRegex?: RegExp[]): boolean => {
  // Return false if the URL is empty
  if (!url || !url.trim()) {
    return false
  }
  // Normalize the URL
  const normalizedUrl = normalizeUrl(url)
  // Check URL against allowed schemes, always including https?://
  const allowed = [new RegExp("^https?://", "i"), ...(allowedSchemesRegex || [])]
  return allowed.some((pattern) => pattern.test(normalizedUrl))
}
