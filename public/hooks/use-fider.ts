import { useContext } from "react"
import { FiderContext } from "@fider/services"

export const useFider = () => useContext(FiderContext)

// Parse allowed schemes and extract all derived data in one go
const parseAllowedSchemesData = (allowedSchemesString: string) => {
  const schemes = allowedSchemesString
    .split("\n")
    .filter((s) => s.trim())
    .map((s) => s.trim())

  const regexArray = schemes.map((s) => new RegExp(s, "i"))

  const protocols = new Set<string>(["http", "https"])
  schemes.forEach((scheme) => {
    if (scheme.startsWith("^")) {
      const match = scheme.match(/^\^([a-zA-Z][a-zA-Z0-9+.-]*):/)
      if (match) {
        protocols.add(match[1])
      }
    }
  })

  return {
    raw: allowedSchemesString,
    regexArray,
    protocols: Array.from(protocols),
  }
}

// Hook that provides all parsed allowed schemes data
const useAllowedSchemesData = () => {
  const fider = useFider()
  const allowedSchemes = fider.session.tenant.allowedSchemes
  return parseAllowedSchemesData(allowedSchemes)
}

// Hook version - get parsed allowed schemes as RegExp array
export const useAllowedSchemesRegex = (): RegExp[] => {
  return useAllowedSchemesData().regexArray
}

// Extract protocol names from allowed schemes for Tiptap Link extension
export const useAllowedProtocols = (): string[] => {
  return useAllowedSchemesData().protocols
}

// For non-React contexts, we need to pass the schemes as a parameter
export const fiderAllowedSchemes = { get: () => useFider().session.tenant.allowedSchemes }
