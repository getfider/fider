import { Button, Moment } from "@fider/components"
import { useFider } from "@fider/hooks"
import { BillingStatus } from "@fider/models"
import { actions } from "@fider/services"
import React from "react"
import { AdminPageContainer } from "../components/AdminBasePage"

interface ManageBillingPageProps {
  status: BillingStatus
  trialEndsAt: string
  subscriptionEndsAt: string
}

const ManageBillingPage = (props: ManageBillingPageProps) => {
  const fider = useFider()

  const subscribe = async () => {
    const result = await actions.generateCheckoutLink()
    if (result.ok) {
      location.href = result.data.url
    }
  }

  return (
    <AdminPageContainer id="p-admin-billing" name="billing" title="Billing" subtitle="Manage your billing settings">
      {props.status === BillingStatus.Trial && (
        <>
          <p>
            Your account is currently on a trial until{" "}
            <strong>
              <Moment locale={fider.currentLocale} format="full" date={props.trialEndsAt} />
            </strong>
            . <br />
            Subscribe before the end of your trial to avoid a service interruption.
          </p>

          <Button variant="primary" onClick={subscribe}>
            Subscribe for $30/mo + Tax
          </Button>
        </>
      )}

      {props.status === BillingStatus.Active && (
        <>
          <p>Your Fider subscription is currently active.</p>
        </>
      )}

      {props.status === BillingStatus.Cancelled && (
        <>
          <p>
            Your Fider subscription is currently cancelled. This site will stay active until{" "}
            <strong>
              <Moment locale={fider.currentLocale} format="full" date={props.subscriptionEndsAt} />
            </strong>
            . <br /> Resubscribe to avoid a service interruption.
          </p>

          <Button variant="primary" onClick={subscribe}>
            Resubscribe for $30/mo + Tax
          </Button>
        </>
      )}
      <p className="text-muted">
        We use{" "}
        <strong>
          <a href="https://paddle.com" target="_blank" rel="noopener" className="text-link">
            Paddle
          </a>
        </strong>{" "}
        as our billing partner. You may see {'"PADDLE.COM FIDER"'} on your credit card.
      </p>
    </AdminPageContainer>
  )
}

export default ManageBillingPage
