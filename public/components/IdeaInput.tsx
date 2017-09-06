import * as React from 'react';
import Textarea from 'react-textarea-autosize';
import { DisplayError, Button, Form, LogInControl } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';

interface IdeaInputProps {
    placeholder: string;
}

interface IdeaInputState {
    title: string;
    description: string;
    focused: boolean;
}

export class IdeaInput extends React.Component<IdeaInputProps, IdeaInputState> {
    private title: HTMLInputElement;
    private form: Form;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.IdeaService)
    public service: IdeaService;

    constructor() {
        super();
        this.state = {
          title: '',
          description: '',
          focused: false
        };
    }

    public componentDidMount() {
      this.title.focus();
    }

    public async submit() {
      if (this.state.title) {
        this.form.clearFailure();
        const result = await this.service.addIdea(this.state.title, this.state.description);
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
        const buttonCss = this.state.title === '' ? 'primary disabled' : 'primary';
        const details = user  ?
                          <div>
                            <div className="field">
                              <Textarea onChange={(e) => this.setState({ description: e.currentTarget.value })} placeholder="Describe your idea" />
                            </div>
                            <Button className={buttonCss} onClick={() => this.form.submit() }>
                              Submit
                            </Button>
                          </div> :
                          <div className="ui message login-message">
                            <div className="header">
                              Log in to raise your voice.
                            </div>
                            <LogInControl />
                          </div>;

        return <Form ref={(f) => { this.form = f!; } } onSubmit={() => this.submit()}>
                    <input id="new-idea-input"
                           type="text"
                           ref={(e) => this.title = e! }
                           maxLength={100}
                           onKeyUp={(e) => { this.setState({ title: e.currentTarget.value }); }}
                           placeholder={this.props.placeholder} />
                { this.state.title && details }
               </Form>;
    }
}
