import React from "react"
import { Button, Moment, Money } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { useFider } from "@fider/hooks"
import { BillingStatus } from "@fider/models"
import { AdminPageContainer } from "../components/AdminBasePage"
import { CardDetails } from "../components/billing/CardDetails"
import { usePaddle } from "../hooks/use-paddle"

interface ManageBillingPageProps {
  paddle: {
    isSandbox: boolean
    vendorId: string
    monthlyPlanId: string
    yearlyPlanId: string
  }
  status: BillingStatus
  trialEndsAt: string
  subscriptionEndsAt: string
  subscription: {
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
    nextPayment: {
      amount: number
      currency: string
      date: string
    }
  }
}

const SubscribePanel = (props: { monthlyPrice: string; subscribeMonthly: () => void; yearlyPrice: string; subscribeYearly: () => void }) => {
  return (
    <div>
      <HStack spacing={4}>
        <VStack spacing={4} className="py-2 px-4 shadow rounded text-center">
          <div>
            <span className="block text-xs p-1 rounded mb-2">&nbsp;</span>
            <span className="text-category">Monthy Subscription</span>
          </div>
          <span className="text-display2 block">
            {props.monthlyPrice}
            <span className="text-title">/month</span>
          </span>
          <Button variant="secondary" size="small" className="mx-auto" onClick={props.subscribeMonthly}>
            Subscribe
          </Button>
        </VStack>
        <VStack spacing={4} className="py-2 px-4 shadow rounded text-center">
          <div>
            <span className="block text-xs bg-yellow-100 p-1 rounded mb-2">
              <strong>2 months free!</strong>
            </span>
            <span className="text-category">Yearly Subscription</span>
          </div>
          <span className="text-display2 block">
            {props.yearlyPrice}
            <span className="text-title">/year</span>
          </span>
          <Button variant="secondary" size="small" className="mx-auto" onClick={props.subscribeYearly}>
            Subscribe
          </Button>
        </VStack>
      </HStack>

      <span className="block mt-4 text-muted">VAT/Tax may be added during checkout.</span>
    </div>
  )
}

const ActiveSubscriptionInformation = (props: ManageBillingPageProps) => {
  const fider = useFider()
  const { isReady, openUrl } = usePaddle({ ...props.paddle })

  const open = (url: string) => () => {
    if (isReady) {
      openUrl(url)
    }
  }

  return (
    <VStack>
      <h3 className="text-display">Your subscription is Active</h3>
      <CardDetails {...props.subscription.paymentInformation} />
      <p>
        Your next payment is{" "}
        <strong>
          <Money amount={props.subscription.nextPayment.amount} currency={props.subscription.nextPayment.currency} locale={fider.currentLocale} />
        </strong>{" "}
        on{" "}
        <strong>
          <Moment locale={fider.currentLocale} format="date" date={props.subscription.nextPayment.date} />
        </strong>
        .
      </p>
      <p>
        You can{" "}
        <a href="#" rel="noopener" className="text-link" onClick={open(props.subscription.updateURL)}>
          update
        </a>{" "}
        your payment information or{" "}
        <a href="#" rel="noopener" className="text-link" onClick={open(props.subscription.cancelURL)}>
          cancel
        </a>{" "}
        your subscription at any time.
      </p>
      <p>
        To change your billing interval from monthly to yearly or vice-versa, please contact us at{" "}
        <a className="text-link" href="mailto:billing@fider.io">
          billing@fider.io
        </a>
        .
      </p>
    </VStack>
  )
}

const CancelledSubscriptionInformation = (props: ManageBillingPageProps) => {
  const fider = useFider()
  const paddle = usePaddle({ ...props.paddle })

  const isExpired = new Date(props.subscriptionEndsAt) <= new Date()

  return (
    <VStack>
      <h3 className="text-display">Your subscription was Cancelled</h3>
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
      <SubscribePanel {...paddle} />
    </VStack>
  )
}

const TrialInformation = (props: ManageBillingPageProps) => {
  const fider = useFider()
  const paddle = usePaddle({ ...props.paddle })

  const isExpired = new Date(props.trialEndsAt) <= new Date()

  return (
    <VStack>
      <h3 className="text-display">Trial</h3>
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

      <SubscribePanel {...paddle} />
    </VStack>
  )
}

const FreeForeverInformation = () => {
  return (
    <VStack>
      <h3 className="text-display">Free!</h3>
      <p>
        This site is on a <strong>Forever Free</strong> subscription, enjoy it! 🎉
      </p>
      <p className="text-muted">
        You can still help us fund the development of Fider by contribution to our{" "}
        <a rel="noopener" target="_blank" className="text-link" href="https://opencollective.com/fider">
          OpenCollective
        </a>
        .
      </p>
    </VStack>
  )
}

const OpenCollectiveInformation = () => {
  return (
    <VStack>
      <h3 className="text-display">Open Source Subscription</h3>
      <p>
        This site is linked to a monthly{" "}
        <a rel="noopener" target="_blank" className="text-link" href="https://opencollective.com/fider">
          OpenCollective
        </a>{" "}
        donation.
      </p>
      <p className="text-muted">Thanks for being a financial support! Keep your monthly donation active to avoid a service interruption.</p>
    </VStack>
  )
}

const ManageBillingPage = (props: ManageBillingPageProps) => {
  const showPaddleFooter = [BillingStatus.Active, BillingStatus.Cancelled, BillingStatus.Trial].includes(props.status)

  return (
    <AdminPageContainer id="p-admin-billing" name="billing" title="Billing" subtitle="Manage your billing settings">
      {props.status === BillingStatus.Trial && <TrialInformation {...props} />}
      {props.status === BillingStatus.Active && <ActiveSubscriptionInformation {...props} />}
      {props.status === BillingStatus.Cancelled && <CancelledSubscriptionInformation {...props} />}
      {props.status === BillingStatus.FreeForever && <FreeForeverInformation />}
      {props.status === BillingStatus.OpenCollective && <OpenCollectiveInformation />}

      {showPaddleFooter && (
        <p className="text-muted mt-4">
          <strong>
            <a href="https://paddle.com" target="_blank" rel="noopener" className="text-link">
              Paddle
            </a>
          </strong>{" "}
          is our billing partner. You may see {'"PADDLE.NET* FIDER"'} on your credit card.
        </p>
      )}
    </AdminPageContainer>
  )
}

export default ManageBillingPage
