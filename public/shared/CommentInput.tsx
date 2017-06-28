import * as React from 'react';

import { Idea } from '../models';
import { DisplayError } from './Common';
import { SocialSignInList } from './SocialSignInList';

import { inject, injectables } from '../di';
import { Session, IdeaService, Failure } from '../services';

interface CommentInputProps {
    idea: Idea;
}

interface CommentInputState {
    content: string;
    clicked: boolean;
    error?: Failure;
}

export class CommentInput extends React.Component<CommentInputProps, CommentInputState> {

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.IdeaService)
    public service: IdeaService;

    constructor() {
        super();
        this.state = {
          content: '',
          clicked: false
        };
    }

    public async submit() {
      this.setState({
        clicked: true,
        error: undefined
      });

      const result = await this.service.addComment(this.props.idea.number, this.state.content);
      if (result.ok) {
        location.reload();
      } else {
        this.setState({
          clicked: false,
          error: result.error,
        });
      }
    }

    public render() {
        const user = this.session.getCurrentUser();
        const buttonClasses = `ui blue labeled submit icon button ${this.state.clicked && 'loading disabled'}`;

        const addComment = user ? <form className="ui reply form">
          <DisplayError error={this.state.error} />
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
            Please log in before leaving a comment.
          </div>
          <p>This only takes a second and you'll be good to go!</p>
          <SocialSignInList orientation="horizontal" size="small" />
        </div>;

        return addComment;
    }
}
