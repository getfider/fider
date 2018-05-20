import * as React from "react";
import { DisplayError, Button, ButtonClickEvent, Form, Textarea, Input, Form2, TextArea } from "@fider/components";
import { SignInModal } from "@fider/components";
import { page, cache, actions, Failure } from "@fider/services";
import { CurrentUser } from "@fider/models";

interface IdeaInputProps {
  user?: CurrentUser;
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
      title: (!!this.props.user && cache.get(CACHE_TITLE_KEY)) || "",
      description: (!!this.props.user && cache.get(CACHE_DESCRIPTION_KEY)) || "",
      focused: false,
      showSignIn: false
    };
    if (this.state) {
      this.props.onTitleChanged(this.state.title);
    }
  }

  public componentDidMount() {
    if (this.props.user && this.title) {
      this.title.focus();
    }
  }

  private onTitleFocused() {
    if (!this.props.user && this.title) {
      this.title.blur();
      this.setState({ showSignIn: true });
    }
  }

  private onTitleChanged(title: string) {
    cache.set(CACHE_TITLE_KEY, title);
    this.setState({ title });
    this.props.onTitleChanged(title);
  }

  private onDescriptionChanged(description: string) {
    cache.set(CACHE_DESCRIPTION_KEY, description);
    this.setState({ description });
  }

  private async submit(event: ButtonClickEvent) {
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
  }

  public render() {
    const details = (
      <>
        <TextArea
          field="description"
          onChange={description => this.onDescriptionChanged(description)}
          value={this.state.description}
          minRows={5}
          placeholder="Describe your idea"
        />
        <Button color="positive" onClick={e => this.submit(e)}>
          Submit
        </Button>
      </>
    );

    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} />
        <Form2 error={this.state.error}>
          <Input
            field="title"
            inputRef={e => (this.title = e!)}
            onFocus={() => this.onTitleFocused()}
            maxLength={100}
            value={this.state.title}
            onChange={title => this.onTitleChanged(title)}
            placeholder={this.props.placeholder}
          />
          {this.state.title && details}
        </Form2>
      </>
    );
  }
}
