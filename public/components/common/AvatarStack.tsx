import "./AvatarStack.scss"

import React, { useState, useEffect } from "react"
import { UserRole } from "@fider/models"
import { classSet } from "@fider/services"

interface AvatarStackProps {
  overlap?: boolean
  users: Array<{
    role?: UserRole
    avatarURL: string
    name: string
  }>
}

export const AvatarStack = (props: AvatarStackProps) => {
  const [loadedCount, setLoadedCount] = useState(0)
  const [errorCount, setErrorCount] = useState(0)
  const [allLoaded, setAllLoaded] = useState(false)
  const totalCount = props.users.length

  useEffect(() => {
    // Reset when users change
    setLoadedCount(0)
    setErrorCount(0)
    setAllLoaded(false)

    // Fallback: show avatars after 2 seconds regardless of load state
    const fallbackTimer = setTimeout(() => {
      setAllLoaded(true)
    }, 2000)

    return () => clearTimeout(fallbackTimer)
  }, [props.users])

  useEffect(() => {
    // Check if all images are loaded or errored
    if (loadedCount + errorCount === totalCount && totalCount > 0) {
      setAllLoaded(true)
    }
  }, [loadedCount, errorCount, totalCount])

  const handleImageLoad = () => {
    setLoadedCount((prev) => prev + 1)
  }

  const handleImageError = () => {
    setErrorCount((prev) => prev + 1)
  }

  const classes = classSet({
    "c-avatar-stack": true,
    "c-avatar-stack--overlap": props.overlap ?? true,
    "c-avatar-stack--loaded": allLoaded,
  })

  return (
    <div className={classes}>
      {props.users.map((x, i) => (
        <img
          key={i}
          className="c-avatar h-8 w-8"
          alt={x.name}
          src={`${x.avatarURL}?size=50`}
          onLoad={handleImageLoad}
          onError={handleImageError}
        />
      ))}
    </div>
  )
}
