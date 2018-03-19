import * as React from "react";
import * as ReactDOM from "react-dom";

import { Idea, CurrentUser } from "@fider/models";
import { Gravatar, UserName, Button, Textarea, DisplayError, SignInControl } from "@fider/components/common";
import { SignInModal } from "@fider/components";

import { page, actions, Failure } from "@fider/services";

interface CommentInputProps {
  user?: CurrentUser;
  idea: Idea;
}

interface CommentInputState {
  content: string;
  error?: Failure;
  showSignIn: boolean;
}

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
  private input!: HTMLTextAreaElement;

  constructor(props: CommentInputProps) {
    super(props);

    this.state = {
      content: "",
      showSignIn: false
    };
  }

  private onTextFocused() {
    if (!this.props.user) {
      this.input.blur();
      this.setState({ showSignIn: true });
    }
  }

  public async submit() {
    this.setState({
      error: undefined
    });

    const result = await actions.createComment(this.props.idea.number, this.state.content);
    if (result.ok) {
      location.reload();
    } else {
      this.setState({
        error: result.error
      });
    }
  }

  public render() {
    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} />
        <div className={`comment-input ${this.props.user && "authenticated"}`}>
          {this.props.user && <Gravatar user={this.props.user} />}
          <div className="ui form">
            {this.props.user && <UserName user={this.props.user} />}
            <DisplayError error={this.state.error} />
            <div className="field">
              <Textarea
                onChange={e => {
                  this.setState({ content: e.currentTarget.value });
                }}
                onFocus={() => this.onTextFocused()}
                inputRef={e => (this.input = e!)}
                rows={1}
                placeholder="Write a comment..."
              />
            </div>
            {this.state.content && (
              <Button color="green" onClick={() => this.submit()}>
                Submit
              </Button>
            )}
          </div>
        </div>
      </>
    );
  }
}
