import * as React from 'react';
import { DisplayError } from '../shared/Common';
import { SocialSignInList } from '../shared/SocialSignInList';

import { inject, injectables } from '../di';
import { Session, IdeaService, Failure } from '../services';

interface IdeaInputState {
    title: string;
    clicked: boolean;
    error?: Failure;
}

export class IdeaInput extends React.Component<{}, IdeaInputState> {
    private description: HTMLTextAreaElement;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.IdeaService)
    public service: IdeaService;

    constructor() {
        super();
        this.state = {
          title: '',
          clicked: false
        };
    }

    public async submit() {
      this.setState({
        clicked: true,
        error: undefined
      });

      const result = await this.service.addIdea(this.state.title, this.description.value);
      if (result.ok) {
        location.href = `/ideas/${result.data.number}/${result.data.slug}`;
      } else if (result.error) {
        this.setState({
          clicked: false,
          error: result.error
        });
      }
    }

    public render() {
        const user = this.session.getCurrentUser();
        const buttonClasses = `ui primary button ${this.state.clicked && 'loading disabled'}`;
        const details = user ?
                        <div>
                          <div className="field">
                            <textarea ref={(ref) => this.description = ref! }
                                      rows={6}
                                      placeholder="Describe your idea (optional)"></textarea>
                            </div>
                              <button className={buttonClasses} onClick={async () => { await this.submit(); } }>
                                Submit Idea
                              </button>
                            </div> :
                          <div>
                          <div className="ui message">
                            <div className="header">
                              Please log in before posting an idea
                            </div>
                            <p>This only takes a second and you'll be good to go!</p>
                              <SocialSignInList orientation="horizontal" size="small" />
                          </div>
                        </div>;

        return <div className="ui form">
                <DisplayError error={this.state.error} />
                <div className="ui fluid input">
                    <input  id="new-idea-input"
                            type="text"
                            maxLength={100}
                            onKeyUp={(e) => { this.setState({ title: e.currentTarget.value }); }}
                            placeholder="Enter your idea, new feature or suggestion here ..." />
                </div>

                { this.state.title.length > 0 && details }
               </div>;
    }
}
