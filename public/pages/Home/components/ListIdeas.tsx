import "./ListIdeas.scss";

import * as React from "react";
import { Post, Tag, IdeaStatus, CurrentUser } from "@fider/models";
import {
  ShowTag,
  ShowIdeaResponse,
  SupportCounter,
  Gravatar,
  MultiLineText,
  Moment,
  ListItem,
  List
} from "@fider/components";

interface ListIdeasProps {
  ideas?: Post[];
  tags: Tag[];
  emptyText: string;
}

const ListIdeaItem = (props: { idea: Post; user?: CurrentUser; tags: Tag[] }) => {
  return (
    <ListItem>
      <SupportCounter idea={props.idea} />
      <div className="c-list-item-content">
        {props.idea.totalComments > 0 && (
          <div className="info right">
            {props.idea.totalComments} <i className="comments outline icon" />
          </div>
        )}
        <a className="c-list-item-title" href={`/ideas/${props.idea.number}/${props.idea.slug}`}>
          {props.idea.title}
        </a>
        <MultiLineText className="c-list-item-description" text={props.idea.description} style="simple" />
        <ShowIdeaResponse status={props.idea.status} response={props.idea.response} />
        {props.tags.map(tag => <ShowTag key={tag.id} size="tiny" tag={tag} />)}
      </div>
    </ListItem>
  );
};

export class ListIdeas extends React.Component<ListIdeasProps, {}> {
  constructor(props: ListIdeasProps) {
    super(props);
  }

  public render() {
    if (!this.props.ideas) {
      return null;
    }

    if (this.props.ideas.length === 0) {
      return <p className="center">{this.props.emptyText}</p>;
    }

    return (
      <List className="c-idea-list" divided={true}>
        {this.props.ideas.map(idea => (
          <ListIdeaItem
            key={idea.id}
            idea={idea}
            tags={this.props.tags.filter(tag => idea.tags.indexOf(tag.slug) >= 0)}
          />
        ))}
      </List>
    );
  }
}
