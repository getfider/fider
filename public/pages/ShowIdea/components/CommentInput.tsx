import * as React from "react";
import * as ReactDOM from "react-dom";

import { Idea, CurrentUser } from "@fider/models";
import { Gravatar, UserName, Button, Textarea, DisplayError, SignInControl } from "@fider/components/common";
import { SignInModal } from "@fider/components";

import { cache, page, actions, Failure } from "@fider/services";

interface CommentInputProps {
  user?: CurrentUser;
  idea: Idea;
}

interface CommentInputState {
  content: string;
  error?: Failure;
  showSignIn: boolean;
}

const CACHE_TITLE_KEY = "CommentInput-Comment-";

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
  private input!: HTMLTextAreaElement;

  constructor(props: CommentInputProps) {
    super(props);

    this.state = {
      content: (!!this.props.user && cache.get(this.getCacheKey())) || "",
      showSignIn: false
    };
  }

  private getCacheKey(): string {
    return `${CACHE_TITLE_KEY}${this.props.idea.id}`;
  }

  private commentChanged = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    cache.set(this.getCacheKey(), e.currentTarget.value);
    this.setState({ content: e.currentTarget.value });
  };

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
      cache.remove(this.getCacheKey());
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
        <div className={`c-comment-input ${this.props.user && "m-authenticated"}`}>
          {this.props.user && <Gravatar user={this.props.user} />}
          <div className="ui form">
            {this.props.user && <UserName user={this.props.user} />}
            <DisplayError error={this.state.error} />
            <div className="field">
              <Textarea
                onChange={this.commentChanged}
                defaultValue={this.state.content}
                onFocus={() => this.onTextFocused()}
                inputRef={e => (this.input = e!)}
                rows={1}
                placeholder="Write a comment..."
              />
            </div>
            {this.state.content && (
              <Button color="positive" onClick={() => this.submit()}>
                Submit
              </Button>
            )}
          </div>
        </div>
      </>
    );
  }
}
