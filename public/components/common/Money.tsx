import React from "react"

interface MomentProps {
  locale: string
  amount: number
  currency: string
}

export const Money = (props: MomentProps) => {
  const formatter = new Intl.NumberFormat(props.locale, {
    style: "currency",
    currency: props.currency,
  })

  return <span>{formatter.format(props.amount)}</span>
}
