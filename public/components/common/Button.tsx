import * as React from 'react';
import { classSet } from '@fider/services';

interface ButtonProps {
  className?: string;
  simple?: boolean;
  disabled?: boolean;
  href?: string;
  color?: 'green' | 'red';
  fluid?: boolean;
  size?: 'mini' | 'tiny' | 'small' | 'large';
  onClick?: (event: ButtonClickEvent) => Promise<any>;
}

interface ButtonState {
  clicked: boolean;
}

import './Button.scss';

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
    size: 'small',
    fluid: false,
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
    // TODO: as-link if simple === true

    const className = classSet({
      'c-button': true,
      [this.props.size!]: this.props.size,
      [this.props.color!]: this.props.color,
      'loading': this.state.clicked,
      'disabled': this.state.clicked || this.props.disabled,
      [this.props.className!]: this.props.className,
    });

    if (this.props.href) {
      return (
        <a href={this.props.href} className={className} onClick={() => this.click()}>
          {this.props.children}
        </a>
      );
    } else {
      return (
        <button type="button" className={className} onClick={(e) => this.click(e)}>
          {this.props.children}
        </button>
      );
    }
  }
}
