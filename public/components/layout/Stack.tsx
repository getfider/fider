import React from "react"
import { classSet } from "@fider/services"

interface StackProps {
  className?: string
  children: React.ReactNode
  onClick?: () => void
  center?: boolean
  divide?: boolean
  justify?: "between" | "evenly"
  spacing?: 0 | 1 | 2 | 4 | 6 | 8
}

const Stack = (props: StackProps, dir: "x" | "y") => {
  const spacing = props.spacing === undefined ? 1 : props.spacing
  const className = classSet({
    [`${props.className}`]: props.className,
    flex: true,
    "flex flex-x": dir === "x",
    "flex flex-y": dir === "y",
    [`flex--spacing-${spacing}`]: spacing > 0 && !props.divide,
    [`flex--divide-${spacing}`]: spacing > 0 && !!props.divide,
    "flex-items-center": dir === "x" && props.center !== false,
    "justify-between": props.justify === "between",
    "justify-evenly": props.justify === "evenly",
  })

  return (
    <div onClick={props.onClick} className={className}>
      {props.children}
    </div>
  )
}

export const HStack = (props: StackProps) => {
  return Stack(props, "x")
}

export const VStack = (props: StackProps) => {
  return Stack(props, "y")
}
