import "./SupportCounter.scss";

import * as React from "react";
import { Idea, IdeaStatus } from "@fider/models";
import { actions, device, classSet } from "@fider/services";
import { SignInModal } from "@fider/components";

interface SupportCounterProps {
  idea: Idea;
}

interface SupportCounterState {
  supported: boolean;
  total: number;
  showSignIn: boolean;
}

export class SupportCounter extends React.Component<SupportCounterProps, SupportCounterState> {
  constructor(props: SupportCounterProps) {
    super(props);
    this.state = {
      supported: props.idea.viewerSupported,
      total: props.idea.totalSupporters,
      showSignIn: false
    };
  }

  private supportOrUndo = async () => {
    if (!Fider.session.isAuthenticated) {
      this.setState({ showSignIn: true });
      return;
    }

    const action = this.state.supported ? actions.removeSupport : actions.addSupport;

    const response = await action(this.props.idea.number);
    if (response.ok) {
      this.setState(state => ({
        supported: !state.supported,
        total: state.total + (state.supported ? -1 : 1)
      }));
    }
  };

  public render() {
    const status = IdeaStatus.Get(this.props.idea.status);

    const className = classSet({
      "m-supported": !status.closed && this.state.supported,
      "m-disabled": status.closed,
      "no-touch": !device.isTouch()
    });

    const vote = (
      <button className={className} onClick={this.supportOrUndo}>
        <i className="caret up icon" />
        {this.state.total}
      </button>
    );

    const disabled = (
      <button className={className}>
        <i className="caret up icon" />
        {this.state.total}
      </button>
    );

    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} />
        <div className="c-support-counter">{status.closed ? disabled : vote}</div>
      </>
    );
  }
}
