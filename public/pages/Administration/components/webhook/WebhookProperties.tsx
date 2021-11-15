import "./WebhookProperties.scss"

import React from "react"

import { HStack, VStack } from "@fider/components/layout"
import { StringObject } from "@fider/services"

interface WebhookPropertiesProps {
  properties: StringObject
  propsName: string
  valueName: string
}

interface PropertyProps {
  value: any
}

const Property = (props: PropertyProps) => {
  const grayValText = (txt: string) => <span className="c-webhook-properties__val--gray c-webhook-properties__val--italic">&lt;{txt}&gt;</span>

  if (Array.isArray(props.value))
    return (
      <VStack spacing={2} divide>
        {props.value.map((val, i) => (
          <Property key={i} value={val} />
        ))}
      </VStack>
    )

  if (props.value === "") return grayValText("empty")
  if (props.value === null) return grayValText("null")
  if (props.value === undefined) return grayValText("undefined")
  if (props.value === true) return <span className="c-webhook-properties__val--green">true</span>
  if (props.value === false) return <span className="c-webhook-properties__val--red">false</span>

  const type = typeof props.value
  switch (type) {
    case "string":
      return <span>{props.value}</span>
    case "number":
    case "bigint":
      return <span className="c-webhook-properties__val--blue">{props.value}</span>
    case "object":
      return <WebhookProperties properties={props.value} propsName="Name" valueName="Value" />
    default:
      return props.value
  }
}

export const WebhookProperties = (props: WebhookPropertiesProps) => {
  return (
    <VStack className="c-webhook-properties" spacing={2} divide>
      <HStack className="flex-wrap" spacing={0}>
        <div className="c-webhook-properties__header">{props.propsName}</div>
        <div className="c-webhook-properties__header">{props.valueName}</div>
      </HStack>
      {Object.entries(props.properties).map(([prop, val]) => (
        <HStack key={prop} className="flex-wrap" spacing={0}>
          <div className="c-webhook-properties__prop">{prop}</div>
          <div className="c-webhook-properties__val">
            <Property value={val} />
          </div>
        </HStack>
      ))}
    </VStack>
  )
}
