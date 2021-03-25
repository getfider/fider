import React from "react"
import { classSet } from "@fider/services"

interface ButtonProps {
  className?: string
  disabled?: boolean
  href?: string
  rel?: "nofollow"
  type?: "button" | "submit"
  color?: "positive" | "danger" | "default" | "cancel"
  fluid?: boolean
  size?: "mini" | "tiny" | "small" | "normal" | "large"
  onClick?: (event: ButtonClickEvent) => Promise<any>
}

interface ButtonState {
  clicked: boolean
}

import "./Button.scss"

export class ButtonClickEvent {
  private shouldEnable = true
  public preventEnable(): void {
    this.shouldEnable = false
  }
  public canEnable(): boolean {
    return this.shouldEnable
  }
}

export class Button extends React.Component<ButtonProps, ButtonState> {
  private unmounted = false

  public static defaultProps: Partial<ButtonProps> = {
    size: "small",
    fluid: false,
    color: "default",
    type: "button",
  }

  public constructor(props: ButtonProps) {
    super(props)
    this.state = {
      clicked: false,
    }
  }

  public componentWillUnmount() {
    this.unmounted = true
  }

  public click = async (e?: React.SyntheticEvent<HTMLElement>) => {
    if (e) {
      e.preventDefault()
      e.stopPropagation()
    }

    if (this.state.clicked) {
      return
    }

    const event = new ButtonClickEvent()
    this.setState({ clicked: true })
    if (this.props.onClick) {
      await this.props.onClick(event)
      if (!this.unmounted && event.canEnable()) {
        this.setState({ clicked: false })
      }
    }
  }

  public render() {
    const className = classSet({
      "c-button": true,
      "m-fluid": this.props.fluid,
      [`m-${this.props.size}`]: this.props.size,
      [`m-${this.props.color}`]: this.props.color,
      "m-loading": this.state.clicked,
      "m-disabled": this.state.clicked || this.props.disabled,
      [this.props.className || ""]: this.props.className,
    })

    if (this.props.href) {
      return (
        <a href={this.props.href} rel={this.props.rel} className={className}>
          {this.props.children}
        </a>
      )
    } else if (this.props.onClick) {
      return (
        <button type={this.props.type} className={className} onClick={this.click}>
          {this.props.children}
        </button>
      )
    } else {
      return (
        <button type={this.props.type} className={className}>
          {this.props.children}
        </button>
      )
    }
  }
}
