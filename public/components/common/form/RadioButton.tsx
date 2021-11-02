import { HStack, VStack } from "@fider/components/layout"
import React, { useState } from "react"

import "./RadioButton.scss"

interface RadioButtonOption {
  value: string
  label: string
}

interface RadioButtonProps {
  label: string
  field: string
  defaultOption: RadioButtonOption
  options: RadioButtonOption[]
  onSelect?: (value: RadioButtonOption) => void
}

export const RadioButton = (props: RadioButtonProps) => {
  const [selected, setSelected] = useState(props.defaultOption)

  const onChange = (option: RadioButtonOption) => () => {
    setSelected(option)
    props.onSelect?.(option)
  }

  const inputs = props.options.map((option) => (
    <HStack key={option.value} className="text-sm">
      <input id={`visibility-${option.value}`} type="radio" name={`input-${props.field}`} checked={selected === option} onChange={onChange(option)} />
      <label htmlFor={`visibility-${option.value}`}>{option.label}</label>
    </HStack>
  ))

  return (
    <div className="c-form-field">
      <label htmlFor={`input-${props.field}`}>{props.label}</label>
      <VStack className="c-radiobutton">{inputs}</VStack>
    </div>
  )
}
