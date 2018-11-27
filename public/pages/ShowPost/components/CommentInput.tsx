import React from "react";

import { Post, CurrentUser } from "@fider/models";
import {
  Gravatar,
  UserName,
  Button,
  DisplayError,
  SignInControl,
  Form,
  MentionableTextArea
} from "@fider/components/common";
import { SignInModal, SuggestionBox } from "@fider/components";

import { cache, actions, Failure, Fider } from "@fider/services";

interface CommentInputProps {
  post: Post;
}

interface CommentInputState {
  content: string;
  error?: Failure;
  showSignIn: boolean;
}

const CACHE_TITLE_KEY = "CommentInput-Comment-";

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
  private input!: React.RefObject<HTMLTextAreaElement>;

  constructor(props: CommentInputProps) {
    super(props);

    this.state = {
      content: (Fider.session.isAuthenticated && cache.get(this.getCacheKey())) || "",
      showSignIn: false
    };
  }

  private getCacheKey(): string {
    return `${CACHE_TITLE_KEY}${this.props.post.id}`;
  }

  private commentChanged = (content: string) => {
    console.log("CommentInput "+ content);
    cache.set(this.getCacheKey(), content);
    this.setState({ content });
  };

  public submit = async () => {
    this.setState({
      error: undefined
    });

    const result = await actions.createComment(this.props.post.number, this.state.content);
    if (result.ok) {
      cache.remove(this.getCacheKey());
      location.reload();
    } else {
      this.setState({
        error: result.error
      });
    }
  };

  private handleOnFocus = () => {
    if (!Fider.session.isAuthenticated && this.input && this.input.current) {
      this.input.current.blur();
      this.setState({ showSignIn: true });
    }
  };

  public render() {
    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} />
        <div className={`c-comment-input ${Fider.session.isAuthenticated && "m-authenticated"}`}>
          {Fider.session.isAuthenticated && <Gravatar user={Fider.session.user} />}
          <Form error={this.state.error}>
            {Fider.session.isAuthenticated && <UserName user={Fider.session.user} />}
            <MentionableTextArea
              placeholder="Write a comment..."
              field="content"
              value={this.state.content}
              minRows={1} //
              onChange={this.commentChanged}
              onFocus={this.handleOnFocus} //
              inputRef={this.input}
            />
            {this.state.content && (
              <Button color="positive" onClick={this.submit}>
                Submit
              </Button>
            )}
          </Form>
        </div>
      </>
    );
  }
}
