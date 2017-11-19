import * as React from 'react';
import { Header, Footer, Button, Gravatar, UserName } from '@fider/components/common';
import { ShowTag } from '@fider/components/ShowTag';
import { inject, injectables } from '@fider/di';
import { Session, TagService } from '@fider/services';
import { Tag, CurrentUser, UserRole } from '@fider/models';

import { NewTagForm, NewTagFormState } from './';
import { Failure } from 'services/http';

interface ManageTagsPageState {
  isAdding: boolean;
  allTags: Tag[];
  deleting?: number;
  editing?: number;
}

import './ManageTagsPage.scss';

export class ManageTagsPage extends React.Component<{}, ManageTagsPageState> {
    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.TagService)
    public tagService: TagService;

    constructor(props: {}) {
        super(props);
        this.state = {
          isAdding: false,
          allTags: this.session.get<Tag[]>('tags') || []
        };
    }

    private async saveNewTag(data: NewTagFormState): Promise<Failure | undefined> {
      const result = await this.tagService.add(data.name, data.color, data.isPublic);
      if (result.ok) {
        this.setState({
          isAdding: false,
          allTags: this.state.allTags.concat(result.data)
        });
      } else {
        return result.error;
      }
    }

    private async updateTag(tag: Tag, data: NewTagFormState): Promise<Failure | undefined> {
      const result = await this.tagService.update(tag.slug, data.name, data.color, data.isPublic);
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
      const result = await this.tagService.delete(tag.slug);
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
                <NewTagForm
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
                  <b>Are you sure?</b> <span>The tag <ShowTag name={t.name} color={t.color} isPublic={t.isPublic} /> will be removed from all ideas.</span>
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
              <ShowTag name={t.name} color={t.color} isPublic={t.isPublic} />
              {
                this.session.isAdmin() && [
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
        <div>
          <Header />
            <div className="page ui container">
              <h2 className="ui header">
                <i className="circular tags icon" />
                <div className="content">
                  Tags
                  <div className="sub header">Manage your account tags.</div>
                </div>
              </h2>

              {
                this.session.isAdmin() && (
                  this.state.isAdding
                  ? <div className="ui segment">
                    <NewTagForm
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
          <Footer />
        </div>
      );
    }
}
