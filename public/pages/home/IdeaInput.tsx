import * as React from 'react';
import { DisplayError, Button, Form, SocialSignInList } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';

interface IdeaInputState {
    title: string;
    focused: boolean;
}

export class IdeaInput extends React.Component<{}, IdeaInputState> {
    private title: HTMLInputElement;
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
          focused: false
        };
    }

    public componentDidMount() {
      this.title.focus();
    }

    public async submit() {
      if (this.state.title) {
        this.form.clearFailure();
        const result = await this.service.addIdea(this.state.title, this.description.value);
        if (result.ok) {
          location.href = `/ideas/${result.data.number}/${result.data.slug}`;
        } else if (result.error) {
          this.form.setFailure(result.error);
        }
      }
    }

    public titleChanged(title: string) {
      this.setState({ title });
    }

    public titleFocused() {
      this.setState({ focused: true });
    }

    public titleUnfocused() {
      this.setState({ focused: false });
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
                            <Button className="primary" onClick={() => this.form.submit() }>
                              Submit Idea
                            </Button>
                          </div> :
                          <div>
                            <div className="ui message">
                              <div className="header">
                                Hey! We need to know who you are
                              </div>
                              <p>This only takes a second and you'll be good to go!</p>
                              <SocialSignInList orientation="horizontal" size="small" />
                            </div>
                          </div>;

        return <Form ref={(f) => { this.form = f!; } } onSubmit={() => this.submit()}>
                    <input id="new-idea-input"
                           type="text"
                           ref={(e) => this.title = e! }
                           maxLength={100}
                           onKeyUp={(e) => { this.setState({ title: e.currentTarget.value }); }}
                           placeholder="Tell us your ideas" />
                { this.state.title && details }
               </Form>;
    }
}
