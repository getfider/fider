import React, { useState } from "react"
import { Button, Icon } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { AdminPageContainer } from "../components/AdminBasePage"
import { http } from "@fider/services"
import IconCheck from "@fider/assets/images/heroicons-check.svg"

import "./ManageBilling.page.scss"

interface ManageBillingPageProps {
  stripeCustomerID: string
  stripeSubscriptionID: string
}

interface PlanCardProps {
  name: string
  price: string
  period: string
  description: string
  features: string[]
  isCurrent: boolean
  isHighlighted?: boolean
  buttonText: string
  buttonVariant: "primary" | "secondary"
  onButtonClick?: () => void
  isLoading: boolean
}

const PlanCard = (props: PlanCardProps) => {
  const showButton = !!props.onButtonClick
  const showCurrentLabel = props.isCurrent && !props.onButtonClick

  const cardClasses = ["c-plan-card p-6", props.isHighlighted ? "c-plan-card--highlighted" : "bg-gray-100", props.isCurrent ? "c-plan-card--current" : ""].join(
    " "
  )

  const textColor = props.isHighlighted ? "text-white" : "text-gray-900"

  return (
    <div className={cardClasses}>
      <VStack spacing={4}>
        <HStack justify="between" align="center">
          <span className={`text-title ${textColor}`}>{props.name}</span>
          {props.isCurrent && <span className="text-xs text-semibold px-2 py-1 rounded-full bg-green-100 text-green-700">CURRENT</span>}
        </HStack>

        <div className="flex flex-items-baseline">
          <span className={`text-2xl text-bold ${textColor}`}>{props.price}</span>
          {props.period && <span className={`text-sm c-plan-card__muted ${props.isHighlighted ? "" : "text-gray-500"}`}>/{props.period}</span>}
        </div>

        <p className={`text-sm c-plan-card__muted ${props.isHighlighted ? "" : "text-gray-600"}`}>{props.description}</p>

        {showButton && (
          <Button variant={props.buttonVariant} onClick={props.onButtonClick} disabled={props.isLoading}>
            {props.isLoading ? "Loading..." : props.buttonText}
          </Button>
        )}
        {showCurrentLabel && <div className="text-center py-2 px-4 text-sm text-medium text-gray-500 bg-gray-200 rounded-md">Current Plan</div>}

        <VStack spacing={2} className={`pt-4 border-t c-plan-card__light ${props.isHighlighted ? "border-gray-700" : "border-gray-200 text-gray-700"}`}>
          {props.features.map((feature, index) => (
            <HStack key={index} spacing={2} align="start">
              <Icon sprite={IconCheck} className="text-green-500" height="16" />
              <span className="text-sm">{feature}</span>
            </HStack>
          ))}
        </VStack>
      </VStack>
    </div>
  )
}

const ManageBillingPage = (props: ManageBillingPageProps) => {
  const [isLoading, setIsLoading] = useState(false)
  const hasSubscription = !!props.stripeSubscriptionID

  const openPortal = async () => {
    setIsLoading(true)
    const result = await http.post<{ url: string }>("/_api/admin/billing/portal")
    if (result.ok) {
      window.location.href = result.data.url
    } else {
      setIsLoading(false)
    }
  }

  const startCheckout = async () => {
    setIsLoading(true)
    const result = await http.post<{ url: string }>("/_api/admin/billing/checkout")
    if (result.ok) {
      window.location.href = result.data.url
    } else {
      setIsLoading(false)
    }
  }

  const freeFeatures = ["Unlimited feedback posts", "Unlimited voters", "Your own subdomain or custom domain", "All core functionality", "Email notifications"]

  const proFeatures = ["Everything in Free", "Content Moderation", "Responsive email support", "Priority support"]

  return (
    <AdminPageContainer id="p-admin-billing" name="billing" title="Billing" subtitle="Manage your subscription and billing">
      <p>Fider is free forever. But if you need advanced features and support, consider upgrading to Pro.</p>
      <div className="c-billing-plans">
        <PlanCard
          name="Free"
          price="$0"
          period="month"
          description="Perfect for getting started with feedback collection."
          features={freeFeatures}
          isCurrent={!hasSubscription}
          buttonText="Downgrade"
          buttonVariant="secondary"
          onButtonClick={hasSubscription ? openPortal : undefined}
          isLoading={isLoading && hasSubscription}
        />

        <PlanCard
          name="Pro"
          price="$25"
          period="month"
          description="For teams that need advanced features and support."
          features={proFeatures}
          isCurrent={hasSubscription}
          isHighlighted={true}
          buttonText={hasSubscription ? "Manage Billing" : "Upgrade to Pro"}
          buttonVariant="primary"
          onButtonClick={hasSubscription ? openPortal : startCheckout}
          isLoading={isLoading}
        />
      </div>
    </AdminPageContainer>
  )
}

export default ManageBillingPage
