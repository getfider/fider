import * as React from 'react';
import Textarea from 'react-textarea-autosize';
import { DisplayError, Button, Form } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';
import { User } from '@fider/models';
import { showLogin } from '@fider/utils/page';

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
    private user: User;

    @inject(injectables.Session)
    private session: Session;

    @inject(injectables.IdeaService)
    private service: IdeaService;

    constructor() {
      super();
      this.user = this.session.getCurrentUser();
      this.state = {
        title: '',
        description: '',
        focused: false
      };
    }

    public componentDidMount() {
      if (this.user) {
        this.title.focus();
      }
    }

    private onTitleFocused() {
      if (!this.user) {
        this.title.blur();
        showLogin();
      }
    }

    private async submit() {
      if (this.state.title) {
        const result = await this.service.addIdea(this.state.title, this.state.description);
        if (result.ok) {
          this.form.clearFailure();
          location.href = `/ideas/${result.data.number}/${result.data.slug}`;
        } else if (result.error) {
          this.form.setFailure(result.error);
        }
      }
    }

    public render() {
      const buttonCss = this.state.title === '' ? 'primary disabled' : 'primary';
      const details = <div>
                        <div className="field">
                          <Textarea onChange={(e) => this.setState({ description: e.currentTarget.value })} placeholder="Describe your idea" />
                        </div>
                        <Button className={buttonCss} onClick={() => this.form.submit() }>
                          Submit
                        </Button>
                      </div>;

      return <Form ref={(f) => { this.form = f!; } } onSubmit={() => this.submit()}>
                  <input id="new-idea-input"
                          type="text"
                          ref={(e) => this.title = e! }
                          onFocus={() => this.onTitleFocused()}
                          maxLength={100}
                          onKeyUp={(e) => { this.setState({ title: e.currentTarget.value }); }}
                          placeholder={this.props.placeholder} />
              { this.state.title && details }
              </Form>;
    }
}
