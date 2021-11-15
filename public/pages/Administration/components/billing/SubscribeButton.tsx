import React, { useEffect, useState } from "react"
import { Button } from "@fider/components"
import { getCountryCode } from "@fider/services"

interface SubscribeButtonProps {
  onClick: () => void
}

interface SubscriptionPrice {
  currency: string
  symbol: string
  countries: string[]
  amount: number
}

const prices = [
  { currency: "USD", symbol: "$", amount: 30, countries: ["US"] },
  {
    currency: "EUR",
    symbol: "€",
    amount: 28,
    countries: ["AD", "AT", "BE", "FI", "FR", "GF", "TF", "DE", "GR", "GP", "VA", "IE", "IT", "LU", "MQ", "YT", "MC", "NL", "PT", "RE", "PM", "SM", "CS", "ES"],
  },
  { currency: "GBP", symbol: "£", amount: 26, countries: ["GB", "UK"] },
]

const getPrice = async (): Promise<SubscriptionPrice> => {
  const countryCode = await getCountryCode()

  if (countryCode) {
    for (const price of prices) {
      if (price.countries.includes(countryCode)) {
        return price
      }
    }
  }

  return prices[0]
}

export const SubscribeButton = (props: SubscribeButtonProps) => {
  const [price, setPrice] = useState<SubscriptionPrice | undefined>()

  useEffect(() => {
    getPrice().then(setPrice)
  }, [])

  if (!price) {
    return null
  }

  return (
    <p>
      <Button variant="primary" onClick={props.onClick}>
        Resubscribe for {price.symbol} {price.amount}/mo
      </Button>
      <span className="block text-muted">Taxes may be added during checkout.</span>
    </p>
  )
}
