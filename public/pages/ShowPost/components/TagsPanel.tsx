import React from "react";
import { Tag, Post } from "@fider/models";
import { actions, Fider } from "@fider/services";
import { ShowTag, List, ListItem } from "@fider/components";
import { TagListItem } from "./TagListItem";
import { FaCheckCircle, FaCog } from "react-icons/fa";

interface TagsPanelProps {
  post: Post;
  tags: Tag[];
}

interface TagsPanelState {
  canEdit: boolean;
  assignedTags: Tag[];
  isEditing: boolean;
}

export class TagsPanel extends React.Component<TagsPanelProps, TagsPanelState> {
  constructor(props: TagsPanelProps) {
    super(props);
    this.state = {
      canEdit: Fider.session.isAuthenticated && Fider.session.user.isCollaborator && this.props.tags.length > 0,
      isEditing: false,
      assignedTags: this.props.tags.filter(t => this.props.post.tags.indexOf(t.slug) >= 0)
    };
  }

  private assignOrUnassignTag = async (tag: Tag) => {
    const idx = this.state.assignedTags.indexOf(tag);
    let assignedTags: Tag[] = [];
    if (idx >= 0) {
      const response = await actions.unassignTag(tag.slug, this.props.post.number);
      if (response.ok) {
        assignedTags = this.state.assignedTags.splice(idx, 1) && this.state.assignedTags;
      }
    } else {
      const response = await actions.assignTag(tag.slug, this.props.post.number);
      if (response.ok) {
        assignedTags = this.state.assignedTags.concat(tag);
      }
    }

    this.setState({
      assignedTags
    });
  };

  private onSubtitleClick = () => {
    if (this.state.canEdit) {
      this.setState({ isEditing: !this.state.isEditing });
    }
  };

  public render() {
    if (!this.state.canEdit && this.state.assignedTags.length === 0) {
      return null;
    }

    const tagsList =
      this.state.assignedTags.length > 0 ? (
        <List className="c-tag-list">
          {this.state.assignedTags.map(tag => (
            <ListItem key={tag.id}>
              <ShowTag tag={tag} />
            </ListItem>
          ))}
        </List>
      ) : (
        <span className="info">None yet</span>
      );

    const editTagsList = this.props.tags.length > 0 && (
      <List className="c-tag-list">
        {this.props.tags.map(tag => (
          <TagListItem
            key={tag.id}
            tag={tag}
            assigned={this.state.assignedTags.indexOf(tag) >= 0}
            onClick={this.assignOrUnassignTag}
          />
        ))}
      </List>
    );

    const subtitleClasses = `subtitle ${this.state.canEdit && "active"}`;
    const icon = this.state.canEdit && (this.state.isEditing ? <FaCheckCircle /> : <FaCog />);

    return (
      <>
        <span className={subtitleClasses} onClick={this.onSubtitleClick}>
          Tags {icon}
        </span>

        {!this.state.isEditing && tagsList}
        {this.state.isEditing && editTagsList}
      </>
    );
  }
}
