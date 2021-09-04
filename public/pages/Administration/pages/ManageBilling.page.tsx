import { Button } from "@fider/components"
import { actions } from "@fider/services"
import React from "react"
import { AdminPageContainer } from "../components/AdminBasePage"

const ManageBillingPage = () => {
  const subscribe = async () => {
    const result = await actions.generateCheckoutLink()
    if (result.ok) {
      location.href = result.data.url
    }
  }

  return (
    <AdminPageContainer id="p-admin-billing" name="billing" title="Billing" subtitle="Manage your billing settings">
      <Button variant="primary" onClick={subscribe}>
        Subscribe
      </Button>
    </AdminPageContainer>
  )
}

export default ManageBillingPage
