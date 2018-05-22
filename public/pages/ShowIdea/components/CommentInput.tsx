import * as React from "react";
import * as ReactDOM from "react-dom";

import { Idea, CurrentUser } from "@fider/models";
import { Gravatar, UserName, Button, DisplayError, SignInControl, TextArea, Form } from "@fider/components/common";
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

  private commentChanged = (content: string) => {
    cache.set(this.getCacheKey(), content);
    this.setState({ content });
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
          <Form error={this.state.error}>
            {this.props.user && <UserName user={this.props.user} />}
            <TextArea
              placeholder="Write a comment..."
              field="content"
              value={this.state.content}
              minRows={1}
              onChange={this.commentChanged}
              onFocus={() => this.onTextFocused()}
              inputRef={e => (this.input = e!)}
            />
            {this.state.content && (
              <Button color="positive" onClick={() => this.submit()}>
                Submit
              </Button>
            )}
          </Form>
        </div>
      </>
    );
  }
}
