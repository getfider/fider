import "./ManageTags.page.scss";

import * as React from "react";
import { ShowTag, Button, Gravatar, UserName, Segment, List, ListItem, Heading } from "@fider/components";
import { AdminBasePage, TagForm, TagFormState } from "../components";

import { Tag, CurrentUser, UserRole } from "@fider/models";
import { actions, Failure } from "@fider/services";

interface ManageTagsPageProps {
  user: CurrentUser;
  tags: Tag[];
}

interface ManageTagsPageState {
  isAdding: boolean;
  allTags: Tag[];
  deleting?: number;
  editing?: number;
}

const tagSorter = (t1: Tag, t2: Tag) => {
  if (t1.name < t2.name) {
    return -1;
  } else if (t1.name > t2.name) {
    return 1;
  }
  return 0;
};

export class ManageTagsPage extends AdminBasePage<ManageTagsPageProps, ManageTagsPageState> {
  public id = "p-admin-tags";
  public name = "tags";
  public icon = "tags";
  public title = "Tags";
  public subtitle = "Manage your site tags";

  constructor(props: ManageTagsPageProps) {
    super(props);
    this.state = {
      isAdding: false,
      allTags: this.props.tags
    };
  }

  private async saveNewTag(data: TagFormState): Promise<Failure | undefined> {
    const result = await actions.createTag(data.name, data.color, data.isPublic);
    if (result.ok) {
      this.setState({
        isAdding: false,
        allTags: this.state.allTags.concat(result.data).sort(tagSorter)
      });
    } else {
      return result.error;
    }
  }

  private async updateTag(tag: Tag, data: TagFormState): Promise<Failure | undefined> {
    const result = await actions.updateTag(tag.slug, data.name, data.color, data.isPublic);
    if (result.ok) {
      tag.name = result.data.name;
      tag.slug = result.data.slug;
      tag.color = result.data.color;
      tag.isPublic = result.data.isPublic;
      this.setState({
        editing: undefined,
        allTags: this.state.allTags.sort(tagSorter)
      });
    } else {
      return result.error;
    }
  }

  private async deleteTag(tag: Tag) {
    const result = await actions.deleteTag(tag.slug);
    const idx = this.state.allTags.indexOf(tag);
    if (result.ok) {
      this.setState({
        deleting: undefined,
        allTags: this.state.allTags.splice(idx, 1) && this.state.allTags
      });
    }
  }

  private getTagList(filter: (tag: Tag) => boolean) {
    return this.state.allTags.filter(filter).map(t => {
      if (this.state.editing === t.id) {
        return (
          <ListItem key={t.id}>
            <TagForm
              name={t.name}
              color={t.color}
              isPublic={t.isPublic}
              onSave={async data => this.updateTag(t, data)}
              onCancel={() => this.setState({ editing: undefined })}
            />
          </ListItem>
        );
      }

      if (this.state.deleting === t.id) {
        return (
          <ListItem key={t.id}>
            <div className="content">
              <b>Are you sure?</b>{" "}
              <span>
                The tag <ShowTag tag={t} /> will be removed from all ideas.
              </span>
            </div>
            <Button className="right" onClick={async () => this.setState({ deleting: undefined })}>
              Cancel
            </Button>
            <Button color="danger" className="right" onClick={() => this.deleteTag(t)}>
              Delete tag
            </Button>
          </ListItem>
        );
      }

      return (
        <ListItem key={t.id}>
          <ShowTag tag={t} />
          {this.props.user.isAdministrator && [
            <Button
              key={0}
              onClick={async () =>
                this.setState({
                  isAdding: false,
                  editing: undefined,
                  deleting: t.id
                })
              }
              className="right"
            >
              <i className="remove icon" />Remove
            </Button>,
            <Button
              key={1}
              onClick={async () =>
                this.setState({
                  isAdding: false,
                  editing: t.id,
                  deleting: undefined
                })
              }
              className="right"
            >
              <i className="edit icon" />Edit
            </Button>
          ]}
        </ListItem>
      );
    });
  }

  public content() {
    const publicTaglist = this.getTagList(t => t.isPublic);
    const privateTagList = this.getTagList(t => !t.isPublic);

    const form =
      this.props.user.isAdministrator &&
      (this.state.isAdding ? (
        <Segment>
          <TagForm onSave={async data => this.saveNewTag(data)} onCancel={() => this.setState({ isAdding: false })} />
        </Segment>
      ) : (
        <Button
          color="positive"
          onClick={async e =>
            this.setState({
              isAdding: true,
              deleting: undefined,
              editing: undefined
            })
          }
        >
          Add new
        </Button>
      ));

    return (
      <>
        {form}
        <Segment>
          <List divided={true}>
            <ListItem>
              <Heading size="small" title="Public Tags" subtitle="These tags are visible to all visitors." />
            </ListItem>
            {publicTaglist.length === 0 ? <ListItem>There aren’t any public tags yet.</ListItem> : publicTaglist}
          </List>
        </Segment>

        <Segment>
          <List divided={true}>
            <ListItem>
              <Heading
                size="small"
                title="Private Tags"
                subtitle="These tags are only visible for members of this site."
              />
            </ListItem>
            {privateTagList.length === 0 ? <ListItem>There aren’t any private tags yet.</ListItem> : privateTagList}
          </List>
        </Segment>
      </>
    );
  }
}
