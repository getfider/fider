import axios from "axios";
import * as React from "react";
import { Idea } from "../models";
import { getCurrentUser } from "../storage";
import { SocialSignInButton } from "./SocialSignInButton";

interface CommentInputProps {
    idea: Idea;
}

interface CommentInputState {
    content: string;
    clicked: boolean;
}

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {
    constructor() {
        super();
        this.state = {
          content: "",
          clicked: false
        };
    }

    public async submit() {
      this.setState({ clicked: true });

      const response = await axios.post(`/api/ideas/${this.props.idea.id}/comments`, {
        content: this.state.content
      });

      location.reload();
    }

    public render() {
        const user = getCurrentUser();
        const buttonClasses = `ui blue labeled submit icon button ${this.state.clicked && "loading disabled"}`;

        const addComment = user ? <form className="ui reply form">
          <div className="field">
            <textarea onKeyUp={(e) => { this.setState({ content: e.currentTarget.value }); }}
                      placeholder="Leave your comment here..."></textarea>
          </div>
          <div className={ buttonClasses } onClick={async () => await this.submit()}>
            <i className="icon edit"></i> Add Comment
          </div>
        </form> :
        <div className="ui message">
          <div className="header">
            Please log in before leaving a comment
          </div>
          <p>This only takes a second and you'll be good to go!</p>
          <SocialSignInButton provider="facebook" small={true} />
          <SocialSignInButton provider="google" small={true} />
        </div>;

        return addComment;
    }
}
