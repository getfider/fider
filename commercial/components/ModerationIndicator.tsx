import React, { useState, useEffect } from "react"
import { Icon } from "@fider/components"
import { useFider } from "@fider/hooks"
import ThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import ThumbsDown from "@fider/assets/images/heroicons-thumbsdown.svg"
import { HStack } from "@fider/components/layout"

export const ModerationIndicator = () => {
  const fider = useFider()
  const [count, setCount] = useState(0)
  const [loading, setLoading] = useState(true)

  // Check if commercial license is available
  // TODO: This will be replaced with proper license validation in Phase 4
  const hasCommercialLicense = true // Placeholder for development

  useEffect(() => {
    const fetchCount = async () => {
      try {
        const response = await fetch("/_api/admin/moderation/count")
        if (response.ok) {
          const data = await response.json()
          setCount(data.count || 0)
        }
      } catch (error) {
        console.error("Failed to fetch moderation count:", error)
      } finally {
        setLoading(false)
      }
    }

    // Only fetch if user is admin/collaborator and moderation is enabled
    if ((fider.session.user.isAdministrator || fider.session.user.isCollaborator) && fider.session.tenant.isModerationEnabled) {
      fetchCount()
    } else {
      setLoading(false)
    }
  }, [fider.session.user, fider.session.tenant.isModerationEnabled])

  // Don't show the indicator if commercial license is not available
  console.log("hasCommercialLicense:", hasCommercialLicense)
  console.log("isModerationEnabled:", fider.session.tenant.isModerationEnabled)
  if (!hasCommercialLicense) {
    return null
  }

  // Don't show the indicator if user is not admin/collaborator or moderation is disabled
  if (!fider.session.user.isAdministrator && !fider.session.user.isCollaborator) {
    return null
  }

  if (!fider.session.tenant.isModerationEnabled) {
    return null
  }

  if (loading) {
    return null
  }

  if (count > 0) {
    return (
      <a href="/admin/moderation">
        <HStack className="bg-green-200 rounded-full px-4">
          <Icon width="18" height="18" sprite={ThumbsUp} />
          <Icon width="18" height="18" sprite={ThumbsDown} />
          <span className="py-2">New ideas and comments waiting</span>
        </HStack>
      </a>
    )
  } else {
    return <></>
  }
}
