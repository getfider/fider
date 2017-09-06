import * as React from 'react';
import { Idea, User, IdeaStatus } from '@fider/models';
import { LogInControl } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService } from '@fider/services';

interface SupportCounterProps {
    user: User;
    idea: Idea;
}

interface SupportCounterState {
    supported: boolean;
    total: number;
}

export class SupportCounter extends React.Component<SupportCounterProps, SupportCounterState> {
    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.IdeaService)
    public service: IdeaService;

    constructor(props: SupportCounterProps) {
        super(props);
        const supportedIdeas = this.session.getArray<number>('supportedIdeas');

        this.state = {
          supported: props.user && supportedIdeas && supportedIdeas.indexOf(props.idea.id) >= 0,
          total: props.idea.totalSupporters
        };
    }

    public async supportOrUndo() {
        if (!this.props.user) {
            $('#login-modal').modal({
                blurring: true
            }).modal('show');
            return;
        }

        const action = this.state.supported ? this.service.removeSupport : this.service.addSupport;

        const response = await action(this.props.idea.number);
        if (response.ok) {
            this.setState({
                supported: !this.state.supported,
                total: this.state.total + (this.state.supported ? -1 : 1)
            });
        } else {
            // TODO: handle this. we should have a global alert box
        }
    }

    public render() {

        const noTouch = !('ontouchstart' in window);
        const status = IdeaStatus.Get(this.props.idea.status);

        const vote = <button
                        className={`ui button ${noTouch ? 'no-touch' : ''} ${this.state.supported ? 'supported' : ''} `}
                        onClick={async () => await this.supportOrUndo()}>
                        <i className="medium caret up icon"></i>
                        { this.state.total }
                     </button>;

        const disabled = <div className="ui button disabled">
                            <i className="medium caret up icon"></i>
                            { this.state.total }
                         </div>;

        return  <div className="support-counter ui">
                    { status.closed ? disabled : vote }
                </div>;
    }
}
