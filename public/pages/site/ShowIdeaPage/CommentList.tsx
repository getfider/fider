import * as React from 'react';
import { Comment } from '@fider/models';
import { UserName, Gravatar, Moment, MultiLineText } from '@fider/components/common';

interface CommentListProps {
  comments: Comment[];
}

export class CommentList extends React.Component<CommentListProps, {}> {
  constructor(props: CommentListProps) {
    super(props);
  }

  public render() {
    return this.props.comments.map((c) => (
      <div key={c.id} className="comment">
        <Gravatar user={c.user} />
        <div className="content">
          <UserName user={c.user} />
          <div className="metadata">
            · <Moment date={c.createdOn} />
          </div>
          <div className="text">
            <MultiLineText text={c.content} style="simple" />
          </div>
        </div>
      </div>
    ));
  }
}
