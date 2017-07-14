import * as React from 'react';

import { Idea } from '@fider/models';
import { Button, DisplayError, SocialSignInList } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, Failure } from '@fider/services';

interface CommentInputProps {
    idea: Idea;
}

interface CommentInputState {
    content: string;
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
          content: ''
        };
    }

    public async submit() {
      this.setState({
        error: undefined
      });

      const result = await this.service.addComment(this.props.idea.number, this.state.content);
      if (result.ok) {
        location.reload();
      } else {
        this.setState({
          error: result.error,
        });
      }
    }

    public render() {
        const user = this.session.getCurrentUser();

        const addComment = user ? <form className="ui reply form">
          <DisplayError error={this.state.error} />
          <div className="field">
            <textarea onKeyUp={(e) => { this.setState({ content: e.currentTarget.value }); }}
                      placeholder="Leave your comment here..."></textarea>
          </div>
          <Button className="primary" onClick={ () => this.submit() }>
            Add Comment
          </Button>
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
