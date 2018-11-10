import React from "react";
import { Tag } from "@fider/models";
import { ListItem, ShowTag } from "@fider/components";
import { FaCheck } from "react-icons/fa";

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
        {this.props.assigned ? <FaCheck /> : <svg className="icon" />}
        <ShowTag tag={this.props.tag} size="mini" circular={true} />
        <span>{this.props.tag.name}</span>
      </ListItem>
    );
  }
}
