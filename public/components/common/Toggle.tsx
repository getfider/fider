import * as React from "react";
import { classSet } from "@fider/services";

interface ToggleProps {
  active: boolean;
  onToggle: (active: boolean) => Promise<any>;
}

interface ToggleState {
  active: boolean;
}

import "./Toggle.scss";

export class Toggle extends React.Component<ToggleProps, ToggleState> {
  public constructor(props: ToggleProps) {
    super(props);
    this.state = {
      active: props.active
    };
  }

  public toggle = async () => {
    this.setState(
      state => ({
        active: !state.active
      }),
      () => {
        this.props.onToggle(this.state.active);
      }
    );
  };

  public render() {
    return (
      <div className="c-toggle" onClick={this.toggle}>
        <input type="checkbox" checked={this.state.active} readOnly={true} />
        <label>
          <div className="switch" />
        </label>
      </div>
    );
  }
}
