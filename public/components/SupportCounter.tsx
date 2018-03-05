import './SupportCounter.scss';

import * as React from 'react';
import { Idea, User, IdeaStatus } from '@fider/models';
import { page, actions, classSet } from '@fider/services';

interface SupportCounterProps {
  user?: User;
  idea: Idea;
}

interface SupportCounterState {
  supported: boolean;
  total: number;
}

export class SupportCounter extends React.Component<SupportCounterProps, SupportCounterState> {

  constructor(props: SupportCounterProps) {
    super(props);
    this.state = {
      supported: props.idea.viewerSupported,
      total: props.idea.totalSupporters
    };
  }

  public async supportOrUndo() {
    if (!this.props.user) {
      page.showSignIn();
      return;
    }

    const action = this.state.supported ? actions.removeSupport : actions.addSupport;

    const response = await action(this.props.idea.number);
    if (response.ok) {
      this.setState((state) => ({
        supported: !state.supported,
        total: state.total + (state.supported ? -1 : 1)
      }));
    }
  }

  public render() {
    const status = IdeaStatus.Get(this.props.idea.status);

    const className = classSet({
      'supported': !status.closed && this.state.supported,
      'disabled': status.closed,
      'no-touch': !('ontouchstart' in window),
    });

    const vote = (
      <button
        className={className}
        onClick={async () => await this.supportOrUndo()}
      >
        <i className="medium caret up icon" />
        {this.state.total}
      </button>
    );

    const disabled = (
      <button className={className}>
        <i className="medium caret up icon" />
        {this.state.total}
      </button>
    );

    return  (
      <div className="c-support-counter">
        {status.closed ? disabled : vote}
      </div>
    );
  }
}
