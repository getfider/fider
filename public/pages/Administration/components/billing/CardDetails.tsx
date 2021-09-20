import React from "react"
import { Icon } from "@fider/components"
import IconVisa from "@fider/assets/images/cc-visa.svg"
import IconAMEX from "@fider/assets/images/cc-amex.svg"
import IconDiners from "@fider/assets/images/cc-diners.svg"
import IconDiscover from "@fider/assets/images/cc-discover.svg"
import IconJCB from "@fider/assets/images/cc-jcb.svg"
import IconMaestro from "@fider/assets/images/cc-maestro.svg"
import IconMasterCard from "@fider/assets/images/cc-mastercard.svg"
import IconUnionPay from "@fider/assets/images/cc-unionpay.svg"
import IconGeneric from "@fider/assets/images/cc-generic.svg"
import { HStack } from "@fider/components/layout"

interface CardDetailsProps {
  cardType: string
  lastFourDigits: string
  expiryDate: string
}

const brands: { [key: string]: SpriteSymbol } = {
  visa: IconVisa,
  master: IconMasterCard,
  american_express: IconAMEX,
  discover: IconDiscover,
  jcb: IconJCB,
  maestro: IconMaestro,
  diners_club: IconDiners,
  unionpay: IconUnionPay,
}

export const CardDetails = (props: CardDetailsProps) => {
  const icon = brands[props.cardType] || IconGeneric

  return (
    <HStack>
      <Icon sprite={icon} className="h-6" />
      <span>{props.lastFourDigits}</span>
      <span>Exp. {props.expiryDate}</span>
    </HStack>
  )
}
