import React from "react";
import { Button, ButtonClickEvent, Input, Form, TextArea, MultiImageUploader } from "@fider/components";
import { SignInModal } from "@fider/components";
import { cache, actions, Failure, Fider } from "@fider/services";
import { ImageUpload } from "@fider/models";

interface PostInputProps {
  placeholder: string;
  onTitleChanged: (title: string) => void;
}

interface PostInputState {
  title: string;
  description: string;
  attachments: ImageUpload[];
  focused: boolean;
  showSignIn: boolean;
  error?: Failure;
}

const CACHE_TITLE_KEY = "PostInput-Title";
const CACHE_DESCRIPTION_KEY = "PostInput-Description";

export class PostInput extends React.Component<PostInputProps, PostInputState> {
  private title?: HTMLInputElement;

  constructor(props: PostInputProps) {
    super(props);
    this.state = {
      title: (Fider.session.isAuthenticated && cache.session.get(CACHE_TITLE_KEY)) || "",
      description: (Fider.session.isAuthenticated && cache.session.get(CACHE_DESCRIPTION_KEY)) || "",
      focused: false,
      showSignIn: false,
      attachments: []
    };

    if (this.state.title) {
      this.props.onTitleChanged(this.state.title);
    }
  }

  private handleTitleFocus = () => {
    if (!Fider.session.isAuthenticated && this.title) {
      this.title.blur();
      this.setState({ showSignIn: true });
    }
  };

  private setTitle = (title: string) => {
    cache.session.set(CACHE_TITLE_KEY, title);
    this.setState({ title });
    this.props.onTitleChanged(title);
  };

  private setDescription = (description: string) => {
    cache.session.set(CACHE_DESCRIPTION_KEY, description);
    this.setState({ description });
  };

  private setAttachments = (attachments: ImageUpload[]) => {
    this.setState({ attachments });
  };

  private submit = async (event: ButtonClickEvent) => {
    if (this.state.title) {
      const result = await actions.createPost(this.state.title, this.state.description, this.state.attachments);
      if (result.ok) {
        this.setState({ error: undefined });
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY);
        location.href = `/posts/${result.data.number}/${result.data.slug}`;
        event.preventEnable();
      } else if (result.error) {
        this.setState({ error: result.error });
      }
    }
  };

  private hideModal = () => {
    this.setState({ showSignIn: false });
  };

  private setInputRef = (e: HTMLInputElement) => {
    this.title = e;
  };

  public render() {
    const details = (
      <>
        <TextArea
          field="description"
          onChange={this.setDescription}
          value={this.state.description}
          minRows={5}
          placeholder="Describe your suggestion (optional)"
        />
        <MultiImageUploader field="attachments" maxUploads={3} previewMaxWidth={100} onChange={this.setAttachments} />
        <Button type="submit" color="positive" onClick={this.submit}>
          Submit
        </Button>
      </>
    );

    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} onClose={this.hideModal} />
        <Form error={this.state.error}>
          <Input
            field="title"
            noTabFocus={!Fider.session.isAuthenticated}
            inputRef={this.setInputRef}
            onFocus={this.handleTitleFocus}
            maxLength={100}
            value={this.state.title}
            onChange={this.setTitle}
            placeholder={this.props.placeholder}
          />
          {this.state.title && details}
        </Form>
      </>
    );
  }
}
