import React from "react"
import { classSet } from "@fider/services"

interface StackProps {
  className?: string
  children: React.ReactNode
  onClick?: () => void
  center?: boolean
  divide?: boolean
  justify?: "between" | "evenly" | "full" | "center"
  spacing?: 0 | 1 | 2 | 4 | 6 | 8
}

const Stack = (props: StackProps, dir: "x" | "y") => {
  const spacing = props.spacing === undefined ? 1 : props.spacing
  const className = classSet({
    [`${props.className}`]: props.className,
    flex: true,
    "flex-x": dir === "x",
    "flex-y": dir === "y",
    [`flex--spacing-${spacing}`]: spacing > 0 && !props.divide,
    [`flex--divide-${spacing}`]: spacing > 0 && !!props.divide,
    "flex-items-center": dir === "x" && props.center !== false,
    [`justify-${props.justify}`]: props.justify,
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
