import * as React from 'react';
import * as ReactDOM from 'react-dom';
import Textarea from 'react-textarea-autosize';

import { Idea, CurrentUser } from '@fider/models';
import { Gravatar, UserName, Button, DisplayError, SignInControl } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';
import { showSignIn } from '@fider/utils/page';

interface CommentInputProps {
    idea: Idea;
}

interface CommentInputState {
    content: string;
    error?: Failure;
}

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
    private input: HTMLTextAreaElement;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.IdeaService)
    public service: IdeaService;

    private user?: CurrentUser;

    constructor() {
        super();
        this.user = this.session.getCurrentUser();

        this.state = {
          content: ''
        };
    }

    private onTextFocused() {
      if (!this.user) {
        this.input.blur();
        showSignIn();
      }
    }

    public async submit() {
      this.setState({
        error: undefined
      });

      const result = await this.service.addComment(this.props.idea.number, this.state.content);
      if (result.ok) {
        location.reload();
      } else {
        this.setState({
          error: result.error,
        });
      }
    }

    public render() {
        const user = this.session.getCurrentUser();

        return <div className={`comment-input ${user && 'authenticated' }`}>
                  { this.user && <Gravatar user={ this.user } /> }
                  <div className="ui form">
                    { this.user && <UserName user={ this.user } /> }
                    <DisplayError error={this.state.error} />
                    <div className="field">
                      <Textarea onChange={(e) => { this.setState({ content: e.currentTarget.value }); }}
                                onFocus={() => this.onTextFocused()}
                                inputRef={(e) => this.input = e!}
                                rows={1}
                                placeholder="Write a comment..."/>
                    </div>
                    { this.state.content &&
                      <Button className="primary" onClick={ () => this.submit() }>
                        Submit
                      </Button>
                    }
                  </div>
                </div>;
    }
}
