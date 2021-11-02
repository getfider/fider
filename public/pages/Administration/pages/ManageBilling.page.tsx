import { Button, Moment, Money } from "@fider/components"
import { VStack } from "@fider/components/layout"
import { useFider } from "@fider/hooks"
import { BillingStatus } from "@fider/models"
import { actions } from "@fider/services"
import React from "react"
import { AdminPageContainer } from "../components/AdminBasePage"
import { CardDetails } from "../components/billing/CardDetails"

interface ManageBillingPageProps {
  status: BillingStatus
  trialEndsAt: string
  subscriptionEndsAt: string
  subscription: {
    signupAt: string
    updateURL: string
    cancelURL: string
    paymentInformation: {
      paymentMethod: string
      cardType: string
      lastFourDigits: string
      expiryDate: string
    }
    lastPayment: {
      amount: number
      currency: string
      date: string
    }
  }
}

const subscribe = async () => {
  const result = await actions.generateCheckoutLink()
  if (result.ok) {
    location.href = result.data.url
  }
}

const ActiveSubscriptionInformation = (props: ManageBillingPageProps) => {
  const fider = useFider()

  return (
    <VStack>
      <h3 className="text-display">Your subscription is Active</h3>
      <CardDetails {...props.subscription.paymentInformation} />
      <p>
        Last payment was{" "}
        <strong>
          <Money amount={props.subscription.lastPayment.amount} currency={props.subscription.lastPayment.currency} locale={fider.currentLocale} />
        </strong>{" "}
        on{" "}
        <strong>
          <Moment locale={fider.currentLocale} format="date" date={props.subscription.lastPayment.date} />
        </strong>
        .
      </p>
      <p>
        You can{" "}
        <a rel="noopener" className="text-link" href={props.subscription.updateURL}>
          update
        </a>{" "}
        your payment information or{" "}
        <a rel="noopener" className="text-link" href={props.subscription.cancelURL}>
          cancel
        </a>{" "}
        your subscription.
      </p>
    </VStack>
  )
}

const CancelledSubscriptionInformation = (props: ManageBillingPageProps) => {
  const fider = useFider()

  const isExpired = new Date(props.trialEndsAt) <= new Date()

  return (
    <>
      {isExpired ? (
        <p>
          Your subscription expired on{" "}
          <strong>
            <Moment locale={fider.currentLocale} format="date" date={props.subscriptionEndsAt} />
          </strong>
          . Resubscribe to remove the read-only constraint from this site.
        </p>
      ) : (
        <p>
          Your subscription is currently cancelled. This site will stay active until{" "}
          <strong>
            <Moment locale={fider.currentLocale} format="date" date={props.subscriptionEndsAt} />
          </strong>
          . <br /> Resubscribe to avoid a service interruption.
        </p>
      )}

      <Button variant="primary" onClick={subscribe}>
        Resubscribe for $30/mo + Tax
      </Button>
    </>
  )
}

const TrialInformation = (props: ManageBillingPageProps) => {
  const fider = useFider()

  const isExpired = new Date(props.trialEndsAt) <= new Date()

  return (
    <>
      {isExpired ? (
        <p>
          Your trial expired on{" "}
          <strong>
            <Moment locale={fider.currentLocale} format="date" date={props.trialEndsAt} />
          </strong>
          . Subscribe to remove the read-only constraint from this site.
        </p>
      ) : (
        <p>
          Your account is currently on a trial until{" "}
          <strong>
            <Moment locale={fider.currentLocale} format="date" date={props.trialEndsAt} />
          </strong>
          . <br />
          Subscribe before the end of your trial to avoid a service interruption.
        </p>
      )}

      <Button variant="primary" onClick={subscribe}>
        Subscribe for $30/mo + Tax
      </Button>
    </>
  )
}

const ManageBillingPage = (props: ManageBillingPageProps) => {
  return (
    <AdminPageContainer id="p-admin-billing" name="billing" title="Billing" subtitle="Manage your billing settings">
      {props.status === BillingStatus.Trial && <TrialInformation {...props} />}
      {props.status === BillingStatus.Active && <ActiveSubscriptionInformation {...props} />}
      {props.status === BillingStatus.Cancelled && <CancelledSubscriptionInformation {...props} />}

      <p className="text-muted mt-4">
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
