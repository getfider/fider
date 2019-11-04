import React from "react";

import { Post, ImageUpload } from "@fider/models";
import { Avatar, UserName, Button, TextArea, Form, MultiImageUploader } from "@fider/components/common";
import { SignInModal } from "@fider/components";

import { cache, actions, Failure, Fider } from "@fider/services";

interface CommentInputProps {
  post: Post;
}

interface CommentInputState {
  content: string;
  error?: Failure;
  showSignIn: boolean;
  attachments: ImageUpload[];
}

const CACHE_TITLE_KEY = "CommentInput-Comment-";

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
  private input!: HTMLTextAreaElement;

  constructor(props: CommentInputProps) {
    super(props);

    this.state = {
      content: (Fider.session.isAuthenticated && cache.session.get(this.getCacheKey())) || "",
      showSignIn: false,
      attachments: []
    };
  }

  private getCacheKey(): string {
    return `${CACHE_TITLE_KEY}${this.props.post.id}`;
  }

  private commentChanged = (content: string) => {
    cache.session.set(this.getCacheKey(), content);
    this.setState({ content });
  };

  private setAttachments = (attachments: ImageUpload[]) => {
    this.setState({ attachments });
  };

  private hideModal = () => {
    this.setState({ showSignIn: false });
  };

  public submit = async () => {
    this.setState({
      error: undefined
    });

    const result = await actions.createComment(this.props.post.number, this.state.content, this.state.attachments);
    if (result.ok) {
      cache.session.remove(this.getCacheKey());
      location.reload();
    } else {
      this.setState({
        error: result.error
      });
    }
  };

  private handleOnFocus = () => {
    if (!Fider.session.isAuthenticated) {
      this.input.blur();
      this.setState({ showSignIn: true });
    }
  };

  private setInputRef = (e: HTMLTextAreaElement) => {
    this.input = e;
  };

  public render() {
    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} onClose={this.hideModal} />
        <div className={`c-comment-input ${Fider.session.isAuthenticated && "m-authenticated"}`}>
          {Fider.session.isAuthenticated && <Avatar user={Fider.session.user} />}
          <Form error={this.state.error}>
            {Fider.session.isAuthenticated && <UserName user={Fider.session.user} />}
            <TextArea
              placeholder="Write a comment..."
              field="content"
              value={this.state.content}
              minRows={1}
              onChange={this.commentChanged}
              onFocus={this.handleOnFocus}
              inputRef={this.setInputRef}
            />
            {this.state.content && (
              <MultiImageUploader
                field="attachments"
                maxUploads={2}
                previewMaxWidth={100}
                onChange={this.setAttachments}
              />
            )}
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
