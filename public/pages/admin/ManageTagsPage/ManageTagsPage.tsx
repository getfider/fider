import * as React from 'react';
import { Button, Gravatar, UserName } from '@fider/components/common';
import { ShowTag } from '@fider/components/ShowTag';
import { Tag, CurrentUser, UserRole } from '@fider/models';

import { TagForm, TagFormState } from './';
import { actions, Failure } from '@fider/services';

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

import './ManageTagsPage.scss';

export class ManageTagsPage extends React.Component<ManageTagsPageProps, ManageTagsPageState> {

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
          allTags: this.state.allTags.concat(result.data)
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
          editing: undefined
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

    public render() {
      const items = this.state.allTags.map(
        (t) => {
          if (this.state.editing === t.id) {
            return (
              <div key={t.id} className="item">
                <TagForm
                  name={t.name}
                  color={t.color}
                  isPublic={t.isPublic}
                  onSave={async (data) => this.updateTag(t, data)}
                  onCancel={() => this.setState({ editing: undefined })}
                />
              </div>
            );
          }

          if (this.state.deleting === t.id) {
            return (
              <div key={t.id} className="item">
                <div className="content">
                  <b>Are you sure?</b> <span>The tag <ShowTag tag={t} /> will be removed from all ideas.</span>
                </div>
                <Button
                  className="right floated"
                  onClick={async () => this.setState({ deleting: undefined })}
                >
                  Cancel
                </Button>
                <Button
                  className="negative right floated"
                  onClick={() => this.deleteTag(t)}
                >
                  Delete tag
                </Button>
              </div>
            );
          }

          return (
            <div key={t.id} className="item">
              <ShowTag tag={t} />
              {
                this.props.user.isAdministrator && [
                  <Button
                    key={0}
                    simple={true}
                    onClick={async () => this.setState({ isAdding: false, editing: undefined, deleting: t.id })}
                    className="icon negative right floated"
                  >
                    <i className="remove icon" />Remove
                  </Button>,
                  <Button
                    key={1}
                    simple={true}
                    onClick={async () => this.setState({ isAdding: false, editing: t.id, deleting: undefined })}
                    className="icon right floated"
                  >
                    <i className="edit icon" />Edit
                  </Button>
                ]
              }
            </div>
          );
        }
      );

      return (
        <div className="page ui container">
          <h2 className="ui header">
            <i className="circular tags icon" />
            <div className="content">
              Tags
              <div className="sub header">Manage your account tags.</div>
            </div>
          </h2>

          {
            this.props.user.isAdministrator && (
              this.state.isAdding
              ? <div className="ui segment">
                <TagForm
                  onSave={async (data) => this.saveNewTag(data)}
                  onCancel={() => this.setState({ isAdding: false })}
                />
              </div>
              : <Button
                className="positive"
                onClick={async (e) => this.setState({ isAdding: true, deleting: undefined, editing: undefined })}
              >
                Add new
              </Button>
            )
          }

          <div className="ui segment">
            <div className="ui middle aligned very relaxed divided list">
            {items.length ? items : <div className="content">There aren’t any tags yet.</div>}
            </div>
          </div>

        </div>
      );
    }
}
