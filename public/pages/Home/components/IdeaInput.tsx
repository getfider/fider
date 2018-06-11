import * as React from "react";
import { DisplayError, Button, ButtonClickEvent, Input, Form, TextArea } from "@fider/components";
import { SignInModal } from "@fider/components";
import { cache, actions, Failure } from "@fider/services";
import { CurrentUser } from "@fider/models";

interface IdeaInputProps {
  placeholder: string;
  onTitleChanged: (title: string) => void;
}

interface IdeaInputState {
  title: string;
  description: string;
  focused: boolean;
  showSignIn: boolean;
  error?: Failure;
}

const CACHE_TITLE_KEY = "IdeaInput-Title";
const CACHE_DESCRIPTION_KEY = "IdeaInput-Description";

export class IdeaInput extends React.Component<IdeaInputProps, IdeaInputState> {
  private title?: HTMLInputElement;

  constructor(props: IdeaInputProps) {
    super(props);
    this.state = {
      title: (!!page.user && cache.get(CACHE_TITLE_KEY)) || "",
      description: (!!page.user && cache.get(CACHE_DESCRIPTION_KEY)) || "",
      focused: false,
      showSignIn: false
    };

    if (this.state.title) {
      this.props.onTitleChanged(this.state.title);
    }
  }

  public componentDidMount() {
    if (page.user && this.title) {
      this.title.focus();
    }
  }

  private handleTitleFocus = () => {
    if (!page.user && this.title) {
      this.title.blur();
      this.setState({ showSignIn: true });
    }
  };

  private setTitle = (title: string) => {
    cache.set(CACHE_TITLE_KEY, title);
    this.setState({ title });
    this.props.onTitleChanged(title);
  };

  private setDescription = (description: string) => {
    cache.set(CACHE_DESCRIPTION_KEY, description);
    this.setState({ description });
  };

  private submit = async (event: ButtonClickEvent) => {
    if (this.state.title) {
      const result = await actions.createIdea(this.state.title, this.state.description);
      if (result.ok) {
        this.setState({ error: undefined });
        cache.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY);
        location.href = `/ideas/${result.data.number}/${result.data.slug}`;
        event.preventEnable();
      } else if (result.error) {
        this.setState({ error: result.error });
      }
    }
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
          placeholder="Describe your idea (optional)"
        />
        <Button color="positive" onClick={this.submit}>
          Submit
        </Button>
      </>
    );

    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} />
        <Form error={this.state.error}>
          <Input
            field="title"
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
