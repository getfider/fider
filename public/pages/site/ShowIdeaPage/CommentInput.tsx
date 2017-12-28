import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { Idea, CurrentUser } from '@fider/models';
import { Gravatar, UserName, Button, Textarea, DisplayError, SignInControl } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';
import { showSignIn } from '@fider/utils/page';

interface CommentInputProps {
  user?: CurrentUser;
  idea: Idea;
}

interface CommentInputState {
  content: string;
  error?: Failure;
}

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
  private input: HTMLTextAreaElement;

  @inject(injectables.IdeaService)
  public service: IdeaService;

  constructor(props: CommentInputProps) {
      super(props);

      this.state = {
        content: ''
      };
  }

  private onTextFocused() {
    if (!this.props.user) {
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

    return (
      <div className={`comment-input ${this.props.user && 'authenticated' }`}>
        {this.props.user && <Gravatar user={this.props.user} />}
        <div className="ui form">
          {this.props.user && <UserName user={this.props.user} />}
          <DisplayError error={this.state.error} />
          <div className="field">
            <Textarea
              onChange={(e) => { this.setState({ content: e.currentTarget.value }); }}
              onFocus={() => this.onTextFocused()}
              inputRef={(e) => this.input = e!}
              rows={1}
              placeholder="Write a comment..."
            />
          </div>
          {
            this.state.content &&
            <Button className="primary" onClick={() => this.submit()}>
              Submit
            </Button>
          }
        </div>
      </div>
    );
  }
}
