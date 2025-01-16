import React from "react"
import { classSet } from "@fider/services"

interface StackProps {
  className?: string
  children: React.ReactNode
  onClick?: () => void
  divide?: boolean
  justify?: "between" | "evenly" | "full" | "center"
  align?: "start" | "center" | "end"
  spacing?: 0 | 1 | 2 | 4 | 6 | 8
}

const Stack = (props: StackProps, dir: "x" | "y") => {
  const spacing = props.spacing === undefined ? 2 : props.spacing
  const className = classSet({
    [`${props.className}`]: props.className,
    flex: true,
    "flex-x": dir === "x",
    "flex-y": dir === "y",
    [`flex--spacing-${spacing}`]: spacing > 0 && !props.divide,
    [`flex--divide-${spacing}`]: spacing > 0 && !!props.divide,
    "justify-between": props.justify === "between",
    "justify-evenly": props.justify === "evenly",
    "justify-full": props.justify === "full",
    "justify-center": props.justify === "center",
    "flex-items-start": props.align === "start",
    "flex-items-center": props.align === "center" || (dir === "x" && props.align === undefined),
    "flex-items-end": props.align === "end",
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
