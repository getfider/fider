import * as React from 'react';
import { DisplayError, Button, Form, SocialSignInList } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';

interface IdeaInputState {
    title: string;
}

export class IdeaInput extends React.Component<{}, IdeaInputState> {
    private description: HTMLTextAreaElement;
    private form: Form;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.IdeaService)
    public service: IdeaService;

    constructor() {
        super();
        this.state = {
          title: '',
        };
    }

    public async submit() {
      this.form.clearFailure();
      const result = await this.service.addIdea(this.state.title, this.description.value);
      if (result.ok) {
        location.href = `/ideas/${result.data.number}/${result.data.slug}`;
      } else if (result.error) {
        this.form.setFailure(result.error);
      }
    }

    public render() {
        const user = this.session.getCurrentUser();
        const details = user ?
                          <div>
                            <div className="field">
                              <textarea ref={(ref) => this.description = ref! }
                                        rows={6}
                                        placeholder="Describe your idea (optional)"></textarea>
                            </div>
                            <Button classes="primary" onClick={() => this.form.submit() }>
                              Submit Idea
                            </Button>
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

        return <Form ref={(f) => { this.form = f!; } } onSubmit={() => this.submit()}>
                <div className="ui fluid input">
                    <input  id="new-idea-input"
                            type="text"
                            maxLength={100}
                            onKeyUp={(e) => { this.setState({ title: e.currentTarget.value }); }}
                            placeholder="Enter your idea, new feature or suggestion here ..." />
                </div>

                { this.state.title && details }
               </Form>;
    }
}
