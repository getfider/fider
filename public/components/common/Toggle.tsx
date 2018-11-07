import "./Toggle.scss";

import React from "react";
import { classSet } from "@fider/services";

interface ToggleProps {
  label?: string;
  active: boolean;
  disabled?: boolean;
  onToggle?: (active: boolean) => Promise<any>;
}

interface ToggleState {
  active: boolean;
}

export class Toggle extends React.Component<ToggleProps, ToggleState> {
  public constructor(props: ToggleProps) {
    super(props);
    this.state = {
      active: props.active
    };
  }

  public toggle = async () => {
    if (!!this.props.disabled) {
      return;
    }

    this.setState(
      state => ({
        active: !state.active
      }),
      () => {
        if (this.props.onToggle) {
          this.props.onToggle(this.state.active);
        }
      }
    );
  };

  public render() {
    const className = classSet({
      "c-toggle": true,
      "m-disabled": !!this.props.disabled
    });

    return (
      <span className={className} onClick={this.toggle}>
        <input type="checkbox" checked={this.state.active} readOnly={true} />
        <label>
          <span className="switch" />
        </label>
        <span className="text">{!!this.props.label && this.props.label}</span>
      </span>
    );
  }
}
