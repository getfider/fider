// Fider's proprietary markdown syntax. Shared by the display renderer (services/markdown.ts)
// and the comment editor (CustomMention / CustomImage) so both read stored content the same way.

// Mentions: @[Display Name]. The name is non-empty, single-line and ends at the first "]"
// (so "@[]" is plain text and never swallows a later "]").
const MENTION_SOURCE = "@\\[([^\\]\\n]+)\\]"
export const MENTION_START = "@["
// Anchored: matches a mention at the start of the input (editor tokenizer).
export const MENTION_AT_START = new RegExp(`^${MENTION_SOURCE}`)
// Global: finds every mention in a string (display renderer). String.replace resets lastIndex.
export const MENTION_GLOBAL = new RegExp(MENTION_SOURCE, "g")

// Inline images: ![](fider-image:<bkey>). Only bkeys in the shape the server/editor generate
// (e.g. attachments/<random>-<name>.png) are accepted, and never with ".." path traversal,
// because the bkey is appended to /static/images/.
const FIDER_IMAGE_PREFIX = "fider-image:"
const BKEY_PATTERN = /^[a-zA-Z0-9_/.-]+$/

export const fiderImageBkey = (href: string | null | undefined): string | undefined => {
  if (!href || !href.startsWith(FIDER_IMAGE_PREFIX)) return undefined
  const bkey = href.substring(FIDER_IMAGE_PREFIX.length)
  if (!BKEY_PATTERN.test(bkey) || bkey.includes("..")) return undefined
  return bkey
}

export const fiderImageMarkdown = (bkey: string): string => `![](${FIDER_IMAGE_PREFIX}${bkey})`
