import "./PostDetailsOverlay.scss"

import React, { ReactNode, useEffect } from "react"

interface PostDetailsOverlayProps {
  children: ReactNode
  onClose: () => void
}

export const PostDetailsOverlay: React.FC<PostDetailsOverlayProps> = ({ children, onClose }) => {
  // Handle escape key to close
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose()
      }
    }

    document.addEventListener("keydown", handleEscape)
    return () => document.removeEventListener("keydown", handleEscape)
  }, [onClose])

  // Prevent body scroll when overlay is open
  useEffect(() => {
    document.body.style.overflow = "hidden"
    return () => {
      document.body.style.overflow = ""
    }
  }, [])

  return (
    <div className="post-details-overlay">
      <div className="post-details-overlay__panel">
        <div className="post-details-overlay__content">{children}</div>
      </div>
    </div>
  )
}
