import React, { useState } from "react"
import { Button } from "@fider/components"
import { VStack } from "@fider/components/layout"
import { AdminPageContainer } from "../components/AdminBasePage"
import { http } from "@fider/services"

interface ManageBilling2PageProps {
  stripeCustomerID: string
  stripeSubscriptionID: string
}

const ManageBilling2Page = (props: ManageBilling2PageProps) => {
  const [isLoading, setIsLoading] = useState(false)
  const hasSubscription = !!props.stripeSubscriptionID

  const openPortal = async () => {
    setIsLoading(true)
    const result = await http.post<{ url: string }>("/_api/admin/billing2/portal")
    if (result.ok) {
      window.location.href = result.data.url
    } else {
      setIsLoading(false)
    }
  }

  const startCheckout = async () => {
    setIsLoading(true)
    const result = await http.post<{ url: string }>("/_api/admin/billing2/checkout")
    if (result.ok) {
      window.location.href = result.data.url
    } else {
      setIsLoading(false)
    }
  }

  return (
    <AdminPageContainer id="p-admin-billing2" name="billing" title="Billing" subtitle="Manage your billing settings">
      {hasSubscription ? (
        <VStack spacing={4}>
          <h3 className="text-display">Pro Plan</h3>
          <p>You are currently on the Pro plan.</p>
          <p>Click the button below to manage your subscription, update payment methods, or view billing history.</p>
          <Button variant="primary" onClick={openPortal} disabled={isLoading}>
            {isLoading ? "Loading..." : "Manage Billing"}
          </Button>
        </VStack>
      ) : (
        <VStack spacing={4}>
          <h3 className="text-display">Free Plan</h3>
          <p>You are currently on the Free plan.</p>
          <p>Upgrade to Pro to unlock additional features and support the development of Fider.</p>
          <Button variant="primary" onClick={startCheckout} disabled={isLoading}>
            {isLoading ? "Loading..." : "Upgrade to Pro"}
          </Button>
        </VStack>
      )}
    </AdminPageContainer>
  )
}

export default ManageBilling2Page
