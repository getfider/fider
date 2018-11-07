import React from "react";
import { Tag } from "@fider/models";
import { ListItem, ShowTag } from "@fider/components";

interface TagListItemProps {
  tag: Tag;
  assigned: boolean;
  onClick: (tag: Tag) => void;
}

export class TagListItem extends React.Component<TagListItemProps, {}> {
  private onClick = () => {
    this.props.onClick(this.props.tag);
  };

  public render() {
    return (
      <ListItem onClick={this.onClick}>
        <i className={`icon ${this.props.assigned && "check"}`} />
        <ShowTag tag={this.props.tag} size="mini" circular={true} />
        <span>{this.props.tag.name}</span>
      </ListItem>
    );
  }
}
