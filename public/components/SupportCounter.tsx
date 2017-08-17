import * as React from 'react';
import { Idea, User, IdeaStatus } from '@fider/models';
import { SocialSignInList } from '@fider/components/common';

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
    private elem: HTMLElement;
    private list: HTMLElement;

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

    public componentDidMount() {
        if (!this.props.user) {
            $(this.elem).popup({
                inline: true,
                hoverable: true,
                popup: this.list,
                on: 'click',
                position : 'bottom left'
            });
            return;
        }
    }

    public async supportOrUndo() {
        if (!this.props.user) {
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

        const vote = <button ref={(e) => { this.elem = e!; } }
                        className={`ui button ${noTouch ? 'no-touch' : ''} ${this.state.supported ? 'supported' : ''} `}
                        onClick={async () => await this.supportOrUndo()}>
                        { this.state.total }
                        <i className="medium thumbs up icon"></i>
                     </button>;

        const disabled = <div className="ui button disabled">
                            { this.state.total }
                            <i className="medium thumbs up icon"></i>
                         </div>;

        return <div>
                    <div className="support-counter ui">
                        { status.closed ? disabled : vote }
                    </div>
                    <div ref={(e) => { this.list = e!; } } className="ui popup transition hidden signin-message">
                        <p className="header">Hey! We need to know who you are</p>
                        <SocialSignInList orientation="horizontal" size="small" />
                    </div>
                </div>;
    }
}
