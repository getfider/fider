import * as React from "react";
import { CurrentUser, Tag, Idea } from "@fider/models";
import { actions } from "@fider/services";
import { ShowTag } from "@fider/components";

interface TagsPanelProps {
  user?: CurrentUser;
  idea: Idea;
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
      canEdit: !!this.props.user && this.props.user.isCollaborator && this.props.tags.length > 0,
      isEditing: false,
      assignedTags: this.props.tags.filter(t => this.props.idea.tags.indexOf(t.slug) >= 0)
    };
  }

  private async assignOrUnassignTag(tag: Tag) {
    const idx = this.state.assignedTags.indexOf(tag);
    let assignedTags: Tag[] = [];
    if (idx >= 0) {
      const response = await actions.unassignTag(tag.slug, this.props.idea.number);
      if (response.ok) {
        assignedTags = this.state.assignedTags.splice(idx, 1) && this.state.assignedTags;
      }
    } else {
      const response = await actions.assignTag(tag.slug, this.props.idea.number);
      if (response.ok) {
        assignedTags = this.state.assignedTags.concat(tag);
      }
    }

    this.setState({
      assignedTags
    });
  }

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
        <div className="ui list tag-list">
          {this.state.assignedTags.map(tag => (
            <div key={tag.id} className="item">
              <ShowTag tag={tag} />
            </div>
          ))}
        </div>
      ) : (
        <span className="info">None yet</span>
      );

    const editTagsList = this.props.tags.length > 0 && (
      <div className="ui list tag-list">
        {this.props.tags.map(tag => (
          <div key={tag.id} className="item selectable" onClick={async () => this.assignOrUnassignTag(tag)}>
            <i className={`icon ${this.state.assignedTags.indexOf(tag) >= 0 && "check"}`} />
            <ShowTag tag={tag} circular={true} />
            <span>{tag.name}</span>
          </div>
        ))}
      </div>
    );

    const subtitleClasses = `subtitle ${this.state.canEdit && "active"}`;
    const icon =
      this.state.canEdit &&
      (this.state.isEditing ? <i className="check circle icon" /> : <i className="setting icon" />);

    return (
      <div>
        <span className={subtitleClasses} onClick={this.onSubtitleClick}>
          Tags {icon}
        </span>

        {!this.state.isEditing && tagsList}
        {this.state.isEditing && editTagsList}
      </div>
    );
  }
}
