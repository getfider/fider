import './SupportCounter.scss';

import * as React from 'react';
import { Idea, User, IdeaStatus } from '@fider/models';
import { page, actions } from '@fider/services';

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
    const noTouch = !('ontouchstart' in window);
    const status = IdeaStatus.Get(this.props.idea.status);

    const vote = (
      <button
        className={`button ${noTouch ? 'no-touch' : ''} ${this.state.supported ? 'supported' : ''} `}
        onClick={async () => await this.supportOrUndo()}
      >
        <i className="medium caret up icon" />
        {this.state.total}
      </button>
    );

    const disabled = (
      <div className="button disabled">
        <i className="medium caret up icon" />
        {this.state.total}
      </div>
    );

    return  (
      <div className="c-support-counter">
        {status.closed ? disabled : vote}
      </div>
    );
  }
}
