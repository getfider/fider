import React from "react";
import { Tag } from "@fider/models";
import { ListItem, ShowTag, Button } from "@fider/components";
import { TagFormState, TagForm } from "./TagForm";
import { actions, Failure, Fider } from "@fider/services";
import { FaTimes, FaEdit } from "react-icons/fa";

interface TagListItemProps {
  tag: Tag;
  onTagEdited: (tag: Tag) => void;
  onTagDeleted: (tag: Tag) => void;
}

interface TagListItemState {
  tag: Tag;
  isDeleting: boolean;
  isEditing: boolean;
}

export class TagListItem extends React.Component<TagListItemProps, TagListItemState> {
  constructor(props: TagListItemProps) {
    super(props);
    this.state = {
      tag: props.tag,
      isDeleting: false,
      isEditing: false
    };
  }

  private startDelete = async () => {
    this.setState({ isDeleting: true, isEditing: false });
  };

  private cancelDelete = async () => {
    this.setState({ isDeleting: false });
  };

  private deleteTag = async () => {
    const result = await actions.deleteTag(this.state.tag.slug);
    if (result.ok) {
      this.setState({
        isDeleting: false
      });
      this.props.onTagDeleted(this.state.tag);
    }
  };

  private startEdit = async () => {
    this.setState({ isDeleting: false, isEditing: true });
  };

  private cancelEdit = async () => {
    this.setState({ isEditing: false });
  };

  private updateTag = async (data: TagFormState): Promise<Failure | undefined> => {
    const result = await actions.updateTag(this.state.tag.slug, data.name, data.color, data.isPublic);
    if (result.ok) {
      const tag = this.state.tag;
      tag.name = result.data.name;
      tag.slug = result.data.slug;
      tag.color = result.data.color;
      tag.isPublic = result.data.isPublic;

      this.setState({
        isEditing: false,
        tag
      });

      this.props.onTagEdited(tag);
    } else {
      return result.error;
    }
  };

  private renderDeleteMode() {
    return (
      <>
        <div className="content">
          <b>Are you sure?</b>{" "}
          <span>
            The tag <ShowTag tag={this.state.tag} /> will be removed from all posts.
          </span>
        </div>
        <Button className="right" onClick={this.cancelDelete} color="cancel">
          Cancel
        </Button>
        <Button color="danger" className="right" onClick={this.deleteTag}>
          Delete tag
        </Button>
      </>
    );
  }

  private renderViewMode() {
    const buttons = Fider.session.user.isAdministrator && [
      <Button size="mini" key={0} onClick={this.startDelete} className="right">
        <FaTimes />
        Delete
      </Button>,
      <Button size="mini" key={1} onClick={this.startEdit} className="right">
        <FaEdit />
        Edit
      </Button>
    ];

    return (
      <>
        <ShowTag tag={this.state.tag} />
        {buttons}
      </>
    );
  }

  private renderEditMode() {
    return (
      <TagForm
        name={this.props.tag.name}
        color={this.props.tag.color}
        isPublic={this.props.tag.isPublic}
        onSave={this.updateTag}
        onCancel={this.cancelEdit}
      />
    );
  }

  public render() {
    const view = this.state.isDeleting
      ? this.renderDeleteMode()
      : this.state.isEditing
      ? this.renderEditMode()
      : this.renderViewMode();

    return <ListItem>{view}</ListItem>;
  }
}
