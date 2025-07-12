import { useContext } from "react"
import { FiderContext } from "@fider/services"

export const useFider = () => useContext(FiderContext)
// useContext() doesn't work in tests
export const fiderAllowedSchemes = { get: () => useFider().session.tenant.allowedSchemes }
