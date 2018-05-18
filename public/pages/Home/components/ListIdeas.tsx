import "./ListIdeas.scss";

import * as React from "react";
import { Idea, Tag, IdeaStatus, CurrentUser } from "@fider/models";
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

const defaultShowCount = 20;

interface ListIdeasProps {
  user?: CurrentUser;
  ideas: Idea[];
  tags: Tag[];
  emptyText: string;
}

interface ListIdeasState {
  showCount: number;
}

const ListIdeaItem = (props: { idea: Idea; user?: CurrentUser; tags: Tag[] }) => {
  return (
    <ListItem>
      <SupportCounter user={props.user} idea={props.idea} />
      <div className="c-list-item-content">
        {props.idea.totalComments > 0 && (
          <div className="info right">
            {props.idea.totalComments} <i className="comments outline icon" />
          </div>
        )}
        <a className="c-list-item-title text primary-hover" href={`/ideas/${props.idea.number}/${props.idea.slug}`}>
          {props.idea.title}
        </a>
        <MultiLineText className="c-list-item-description" text={props.idea.description} style="simple" />
        <ShowIdeaResponse status={props.idea.status} response={props.idea.response} />
        {props.tags.map(tag => <ShowTag key={tag.id} size="mini" tag={tag} />)}
      </div>
    </ListItem>
  );
};

export class ListIdeas extends React.Component<ListIdeasProps, ListIdeasState> {
  constructor(props: ListIdeasProps) {
    super(props);
    this.state = {
      showCount: defaultShowCount
    };
  }

  private showMore(event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void {
    event.preventDefault();
    this.setState({
      showCount: this.state.showCount + defaultShowCount
    });
  }

  public render() {
    if (this.props.ideas.length === 0) {
      return <p>{this.props.emptyText}</p>;
    }

    const ideasToList = this.props.ideas.slice(0, this.state.showCount);
    return (
      <List className="c-idea-list" divided={true}>
        {ideasToList.map(idea => (
          <ListIdeaItem
            key={idea.id}
            user={this.props.user}
            idea={idea}
            tags={this.props.tags.filter(tag => idea.tags.indexOf(tag.slug) >= 0)}
          />
        ))}
        {this.props.ideas.length > this.state.showCount && (
          <h5
            className="c-idea-list-show-more primary"
            onTouchEnd={e => this.showMore(e)}
            onClick={e => this.showMore(e)}
          >
            View {this.props.ideas.length - this.state.showCount} more ideas
          </h5>
        )}
      </List>
    );
  }
}
