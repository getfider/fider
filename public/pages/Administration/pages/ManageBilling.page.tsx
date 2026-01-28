import React, { useState } from "react"
import { Button, Icon } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { AdminPageContainer } from "../components/AdminBasePage"
import { http } from "@fider/services"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconInfo from "@fider/assets/images/heroicons-information-circle.svg"

import "./ManageBilling.page.scss"

interface ManageBillingPageProps {
  stripeCustomerID: string
  stripeSubscriptionID: string
  licenseKey: string
  paddleSubscriptionID: string
  hasCommercialFeatures: boolean
}

interface PlanFeature {
  text: string
  isNegative?: boolean
}

interface PlanCardProps {
  name: string
  price?: string
  period?: string
  description: string
  features: (string | PlanFeature)[]
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
          {props.price ? (
            <>
              <span className={`text-2xl text-bold ${textColor}`}>{props.price}</span>
              {props.period && <span className={`text-sm c-plan-card__muted ${props.isHighlighted ? "" : "text-gray-500"}`}>/{props.period}</span>}
            </>
          ) : (
            <span className="text-2xl">&nbsp;</span>
          )}
        </div>

        <p className={`text-sm c-plan-card__muted ${props.isHighlighted ? "" : "text-gray-600"}`}>{props.description}</p>

        {showButton && (
          <Button variant={props.buttonVariant} onClick={props.onButtonClick} disabled={props.isLoading}>
            {props.isLoading ? "Loading..." : props.buttonText}
          </Button>
        )}
        {showCurrentLabel && <div className="text-center py-2 px-4 text-sm text-medium text-gray-500 bg-gray-200 rounded-md">Current Plan</div>}
        {!showButton && !showCurrentLabel && <div className="py-2 px-4 text-sm">&nbsp;</div>}

        <VStack spacing={2} className={`pt-4 border-t c-plan-card__light ${props.isHighlighted ? "border-gray-700" : "border-gray-200 text-gray-700"}`}>
          {props.features.map((feature, index) => {
            const featureText = typeof feature === "string" ? feature : feature.text
            const isNegative = typeof feature === "object" && feature.isNegative
            return (
              <HStack key={index} spacing={2} align="start">
                <Icon sprite={isNegative ? IconX : IconCheck} className={isNegative ? "text-red-500" : "text-green-500"} height="16" />
                <span className="text-sm">{featureText}</span>
              </HStack>
            )
          })}
        </VStack>
      </VStack>
    </div>
  )
}

const PaddleMigrationBanner = () => {
  return (
    <div className="bg-blue-50 p-4 rounded mb-6 border border-blue-200">
      <HStack spacing={2} align="start">
        <Icon sprite={IconInfo} className="text-blue-600 flex-shrink-0 mt-0.5" height="20" />
        <VStack spacing={1}>
          <p className="text-sm text-gray-900 text-medium">Migration to Stripe Billing</p>
          <p className="text-sm text-gray-700">
            You&apos;re currently entitled to pro features because of your existing subscription. Switch to our new Stripe billing to manage your subscription
            and save money.
          </p>
        </VStack>
      </HStack>
    </div>
  )
}

const ManageBillingPage = (props: ManageBillingPageProps) => {
  const [isLoading, setIsLoading] = useState(false)

  // Detect Paddle customers who need to migrate
  const isPaddleCustomer = Boolean(props.paddleSubscriptionID && !props.stripeSubscriptionID)

  // Display as commercial only if they're truly a Stripe customer
  const displayAsCommercial = props.hasCommercialFeatures && !isPaddleCustomer

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

  const freeFeatures = ["250 suggestions", "Unlimited voters", "Your own subdomain or custom domain", "All core functionality"]

  const proFeatures = ["Everything in free", "Unlimited suggestions", "Content moderation", "Billing month to month", "Responsive email support"]

  const legacyProFeatures: PlanFeature[] = [
    { text: "Same features as Pro" },
    { text: "More expensive", isNegative: true },
    { text: "Billing management not supported", isNegative: true },
  ]

  return (
    <AdminPageContainer id="p-admin-billing" name="billing" title="Billing" subtitle="Manage your subscription and billing">
      <p>Fider is free forever. But if you need advanced features and support, consider upgrading to Pro.</p>

      {isPaddleCustomer && <PaddleMigrationBanner />}

      {displayAsCommercial && props.licenseKey && (
        <div className="bg-blue-50 p-4 rounded mb-4 border border-blue-200">
          <p className="text-sm text-gray-700">
            If you want to run Fider self-hosted with commercial features,{" "}
            <a href="/admin/advanced" className="text-blue-600 hover:text-blue-800 underline">
              see advanced
            </a>
            .
          </p>
        </div>
      )}

      <div className="c-billing-plans">
        <PlanCard
          name="Free"
          price="$0"
          period="month"
          description="Perfect for getting started with feedback collection."
          features={freeFeatures}
          isCurrent={!displayAsCommercial && !isPaddleCustomer}
          buttonText="Downgrade"
          buttonVariant="secondary"
          onButtonClick={displayAsCommercial ? openPortal : undefined}
          isLoading={isLoading && displayAsCommercial}
        />

        {isPaddleCustomer && (
          <PlanCard
            name="Legacy Pro"
            description="Your current plan from our previous billing system."
            features={legacyProFeatures}
            isCurrent={true}
            buttonText="Current Plan"
            buttonVariant="secondary"
            isLoading={false}
          />
        )}

        <PlanCard
          name="Pro"
          price="$25"
          period="month"
          description="For teams that need advanced features and support."
          features={proFeatures}
          isCurrent={displayAsCommercial}
          isHighlighted={true}
          buttonText={displayAsCommercial ? "Manage Billing" : isPaddleCustomer ? "Switch to new Pro Plan" : "Upgrade to Pro"}
          buttonVariant="primary"
          onButtonClick={displayAsCommercial ? openPortal : startCheckout}
          isLoading={isLoading}
        />
      </div>
    </AdminPageContainer>
  )
}

export default ManageBillingPage
