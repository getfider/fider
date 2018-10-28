import * as React from "react";
import * as ReactDOM from "react-dom";

import { Post, CurrentUser, User } from "@fider/models";
import { Gravatar, UserName, Button, DisplayError, SignInControl, TextAreaTriggerStart, TextArea, Form } from "@fider/components/common";
import { SignInModal, SuggestionBox } from "@fider/components";

import { cache, actions, Failure, Fider, mentions } from "@fider/services";

interface CommentInputProps {
  post: Post;
}

interface CommentInputState {
  content: string;
  error?: Failure;
  showSignIn: boolean;
  showMentionSuggestion: boolean;
  mentionSuggestionTop: number;
  mentionSuggestionLeft: number;
  mentionText: User[];
  mentionSelected: number;
}

const CACHE_TITLE_KEY = "CommentInput-Comment-";

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
  private input!: HTMLTextAreaElement;

  constructor(props: CommentInputProps) {
    super(props);

    this.state = {
      content: (Fider.session.isAuthenticated && cache.get(this.getCacheKey())) || "",
      showSignIn: false,
      showMentionSuggestion: false,
      mentionSuggestionTop: 0,
      mentionSuggestionLeft: 0,
      mentionText: [],
      mentionSelected: 0
    };
  }

  private getCacheKey(): string {
    return `${CACHE_TITLE_KEY}${this.props.post.id}`;
  }

  private commentChanged = (content: string) => {
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
    if (!Fider.session.isAuthenticated) {
      this.input.blur();
      this.setState({ showSignIn: true });
    }
  };

  private setInputRef = (e: HTMLTextAreaElement) => {
    this.input = e;
  };

  private mentionStart = (e : TextAreaTriggerStart) => {
    this.setState({
      mentionSuggestionLeft : e.left,
      mentionSuggestionTop  : e.top,
      showMentionSuggestion: true
    })
  }
  private mentionChange = (text : string) => {

    mentions.get(text).then((response) => {
      this.setState({
        mentionText : response.data
      });
    })

  }

  private mentionEnd = () => {
    this.setState({
      showMentionSuggestion: false,
    })
  }

  // This is too rudimentary as it will override whatever is in front of the mention
  private mentionSelected = (mentionStartIndex : number, cursorPosition: number) =>{
    const {content} = this.state;
    const contentBeforeMention = content.substring(0, mentionStartIndex);
    const contentAfterCursor = content.substring(cursorPosition);
    const mention = this.state.mentionText[this.state.mentionSelected].name;

    const toReplace = contentBeforeMention + mention + contentAfterCursor;
    this.setState({showMentionSuggestion: false, content: toReplace})
  }

  private mentionArrow = (selectEvent : string) => {
    switch (selectEvent){
      case ("up"): {
        const selected = (this.state.mentionSelected - 1) % this.state.mentionText.length;
        this.setState({mentionSelected : selected})
        break;
      }
      case ("down"): {
        const selected = (this.state.mentionSelected + 1) % this.state.mentionText.length;
        this.setState({mentionSelected : selected})
        break;
      }
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
            <TextArea
              placeholder="Write a comment..."
              field="content"
              value={this.state.content}
              minRows={1}
              onChange={this.commentChanged}
              onFocus={this.handleOnFocus}
              onTriggerStart={this.mentionStart}
              onTriggerChange={this.mentionChange}
              onTriggerEnd={this.mentionEnd}
              onTriggerSelected={this.mentionSelected}
              onTriggerArrow={this.mentionArrow}
              inputRef={this.setInputRef}
            >
            <SuggestionBox
              top={this.state.mentionSuggestionTop}
              left={this.state.mentionSuggestionLeft}
              shown={this.state.showMentionSuggestion}
              data={this.state.mentionText.map(u => u.name)}
              itemSelected={this.state.mentionSelected}
            />
            </TextArea>
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
