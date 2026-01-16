import React from "react"
import { Header } from "@fider/components"
import { Button } from "@fider/components/common"

// Dynamically import the commercial component
const CommercialContentModerationPage = React.lazy(() => import("@commercial/pages/Administration/ContentModeration.page"))

const ContentModerationPage = () => {
  // Check if commercial license is available
  // TODO: This will be replaced with proper license validation in Phase 4
  const hasCommercialLicense = true // Placeholder for development

  // If no commercial license, show upgrade message
  if (!hasCommercialLicense) {
    return (
      <>
        <Header />
        <div id="p-admin-moderation" className="page container">
          <h1 className="text-large">Content Moderation</h1>
          <div className="text-center p-8 border rounded mt-4">
            <h2 className="text-medium mb-4">Commercial Feature</h2>
            <p className="text-body mb-4">Content moderation is a commercial feature. Please upgrade to access the moderation queue.</p>
            <Button size="small" variant="primary" href="/upgrade">
              Upgrade to Commercial
            </Button>
          </div>
        </div>
      </>
    )
  }

  // If licensed, load the commercial component
  return (
    <React.Suspense fallback={<div>Loading...</div>}>
      <CommercialContentModerationPage />
    </React.Suspense>
  )
}

export default ContentModerationPage
