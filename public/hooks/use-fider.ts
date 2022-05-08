import { useContext } from "react"
import { FiderContext } from "@fider/services"

export const useFider = () => useContext(FiderContext)
