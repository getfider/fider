import React, { useState, useEffect } from "react"
import { Icon } from "@fider/components"
import { useFider } from "@fider/hooks"
import IconShield from "@fider/assets/images/heroicons-shieldcheck.svg"

export const ModerationIndicator = () => {
  const fider = useFider()
  const [count, setCount] = useState(0)
  const [loading, setLoading] = useState(true)

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
    if ((fider.session.user.isAdministrator || fider.session.user.isCollaborator) && 
        fider.session.tenant.isModerationEnabled) {
      fetchCount()
    } else {
      setLoading(false)
    }
  }, [fider.session.user, fider.session.tenant.isModerationEnabled])

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

  return (
    <a 
      href="/admin/moderation" 
      className="relative c-themeswitcher"
      title={`${count} item(s) awaiting moderation`}
    >
      <Icon sprite={IconShield} className="h-6 text-gray-500 hover:text-gray-700" />
      {count > 0 && (
        <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
          {count > 99 ? "99+" : count}
        </span>
      )}
    </a>
  )
}