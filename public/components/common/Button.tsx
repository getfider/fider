import * as React from "react";
import { classSet } from "@fider/services";

interface ButtonProps {
  className?: string;
  disabled?: boolean;
  href?: string;
  color?: "positive" | "danger" | "default";
  fluid?: boolean;
  size?: "mini" | "tiny" | "small" | "normal" | "large";
  onClick?: (event: ButtonClickEvent) => Promise<any>;
}

interface ButtonState {
  clicked: boolean;
}

import "./Button.scss";

export class ButtonClickEvent {
  private shouldEnable = true;
  public preventEnable(): void {
    this.shouldEnable = false;
  }
  public canEnable(): boolean {
    return this.shouldEnable;
  }
}

export class Button extends React.Component<ButtonProps, ButtonState> {
  private unmounted: boolean = false;

  public static defaultProps: Partial<ButtonProps> = {
    size: "small",
    fluid: false,
    color: "default"
  };

  public constructor(props: ButtonProps) {
    super(props);
    this.state = {
      clicked: false
    };
  }

  public componentWillUnmount() {
    this.unmounted = true;
  }

  public async click(e?: React.MouseEvent<HTMLButtonElement>) {
    if (e) {
      e.preventDefault();
      e.stopPropagation();
    }

    if (this.state.clicked) {
      return;
    }

    const event = new ButtonClickEvent();
    this.setState({ clicked: true });
    if (this.props.onClick) {
      await this.props.onClick(event);
      if (!this.unmounted && event.canEnable()) {
        this.setState({ clicked: false });
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
      [this.props.className!]: this.props.className
    });

    if (this.props.href) {
      return (
        <a href={this.props.href} className={className}>
          {this.props.children}
        </a>
      );
    } else {
      return (
        <button type="button" className={className} onClick={e => this.click(e)}>
          {this.props.children}
        </button>
      );
    }
  }
}
