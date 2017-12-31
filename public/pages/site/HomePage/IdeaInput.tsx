import * as React from 'react';
import { DisplayError, Button, ButtonClickEvent, Form, Textarea } from '@fider/components/common';
import { page, cache, actions, Failure } from '@fider/services';
import { CurrentUser } from '@fider/models';

interface IdeaInputProps {
  user?: CurrentUser;
  placeholder: string;
}

interface IdeaInputState {
    title: string;
    description: string;
    focused: boolean;
}

const CACHE_TITLE_KEY = 'IdeaInput-Title';
const CACHE_DESCRIPTION_KEY = 'IdeaInput-Description';

export class IdeaInput extends React.Component<IdeaInputProps, IdeaInputState> {
    private title: HTMLInputElement;
    private form: Form;

    constructor(props: IdeaInputProps) {
      super(props);
      this.state = {
        title: cache.get(CACHE_TITLE_KEY) || '',
        description: cache.get(CACHE_DESCRIPTION_KEY) || '',
        focused: false
      };
    }

    public componentDidMount() {
      if (this.props.user) {
        this.title.focus();
      }
    }

    private onTitleFocused() {
      if (!this.props.user) {
        this.title.blur();
        page.showSignIn();
      }
    }

    private onTitleChanged(title: string) {
      cache.set(CACHE_TITLE_KEY, title);
      this.setState({ title });
    }

    private onDescriptionChanged(description: string) {
      cache.set(CACHE_DESCRIPTION_KEY, description);
      this.setState({ description });
    }

    private async submit(event: ButtonClickEvent) {
      if (this.state.title) {
        const result = await actions.addIdea(this.state.title, this.state.description);
        if (result.ok) {
          cache.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY);
          this.form.clearFailure();
          location.href = `/ideas/${result.data.number}/${result.data.slug}`;
          event.preventEnable();
        } else if (result.error) {
          this.form.setFailure(result.error);
        }
      }
    }

    public render() {
      const buttonCss = this.state.title === '' ? 'primary disabled' : 'primary';
      const details = (
        <div>
          <div className="field">
            <Textarea
              onChange={(e) => this.onDescriptionChanged(e.currentTarget.value)}
              defaultValue={this.state.description}
              placeholder="Describe your idea"
            />
          </div>
          <Button className={buttonCss} onClick={(e) => this.submit(e)}>
            Submit
          </Button>
        </div>
      );

      return (
      <Form ref={(f) => { this.form = f!; }} >
        <input
          id="new-idea-input"
          type="text"
          ref={(e) => this.title = e!}
          onFocus={() => this.onTitleFocused()}
          maxLength={100}
          defaultValue={this.state.title}
          onChange={(e) => this.onTitleChanged(e.currentTarget.value)}
          placeholder={this.props.placeholder}
        />
        {this.state.title && details}
      </Form>
      );
    }
}
