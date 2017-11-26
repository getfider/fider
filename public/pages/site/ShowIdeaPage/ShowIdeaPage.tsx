import * as React from 'react';

import { CurrentUser, Comment, Idea, Tag } from '@fider/models';
import { setTitle } from '@fider/utils/page';
import { SupportCounter } from '@fider/components/SupportCounter';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { CommentInput } from './CommentInput';
import { ResponseForm } from './ResponseForm';
import { ShowTag } from '@fider/components/ShowTag';
import { DisplayError, Button, Textarea, UserName, Gravatar, Moment, Form, MultiLineText, Footer, Header, SocialSignInButton } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session, IdeaService, TagService, Failure } from '@fider/services';

import './ShowIdeaPage.scss';

interface ShowIdeaPageState {
  editMode: boolean;
  editTags: boolean;
  newTitle: string;
  newDescription: string;
  assignedTags: Tag[];
  error?: Failure;
}

export class ShowIdeaPage extends React.Component<{}, ShowIdeaPageState> {
  private user?: CurrentUser;
  private idea: Idea;
  private comments: Comment[];
  private tags: Tag[];

  @inject(injectables.Session)
  public session: Session;

  @inject(injectables.IdeaService)
  public ideaService: IdeaService;

  @inject(injectables.TagService)
  public tagService: TagService;

  constructor(props: {}) {
    super(props);

    this.user = this.session.getCurrentUser();
    this.idea = this.session.get<Idea>('idea');
    this.comments = this.session.getArray<Comment>('comments');
    this.tags = this.session.getArray<Tag>('tags');

    this.state = {
      editMode: false,
      editTags: false,
      newTitle: this.idea.title,
      newDescription: this.idea.description,
      assignedTags: this.tags.filter((t) => this.idea.tags.indexOf(t.id) >= 0),
    };

    setTitle(`${this.idea.title} · ${document.title}`);
  }

  private async saveChanges() {
    const result = await this.ideaService.updateIdea(this.idea.number, this.state.newTitle, this.state.newDescription);
    if (result.ok) {
      this.setState({
        error: undefined,
        editMode: false
      });
      this.idea.title = this.state.newTitle;
      this.idea.description = this.state.newDescription;
      this.forceUpdate();
    } else {
      this.setState({
        error: result.error
      });
    }
  }

  private async assignOrUnassignTag(tag: Tag) {
    const idx = this.state.assignedTags.indexOf(tag);
    let assignedTags: Tag[] = [];
    if (idx >= 0) {
      const response = await this.tagService.unassign(tag.slug, this.idea.number);
      if (response.ok) {
        assignedTags = this.state.assignedTags.splice(idx, 1) && this.state.assignedTags;
      }
    } else {
      const response = await this.tagService.assign(tag.slug, this.idea.number);
      if (response.ok) {
        assignedTags = this.state.assignedTags.concat(tag);
      }
    }

    this.setState({
      assignedTags
    });
  }

  public render() {
    const commentsList = this.comments.map((c) => (
      <div key={c.id} className="comment">
        <Gravatar user={c.user} />
        <div className="content">
          <UserName user={c.user} />
          <div className="metadata">
            · <Moment date={c.createdOn} />
          </div>
          <div className="text">
            <MultiLineText text={c.content} style="simple" />
          </div>
        </div>
      </div>
    ));

    const tagsList = this.state.assignedTags.length
      ?
      this.state.assignedTags.map((tag) => (
        <div key={tag.id} className="item">
          <ShowTag tag={tag} />
        </div>
      ))
      :
      <span className="info">None yet</span>;

    return (
      <div>
        <Header />
        <div className="page ui container">
          <div className="ui stackable grid container">
            <div className="thirteen wide column">
              <div className="ui items unstackable">
                <div className="item">
                  <SupportCounter user={this.user} idea={this.idea} />

                  <div className="idea-header">
                    { this.state.editMode
                      ? [
                        <div key={1} className="ui input huge fluid">
                          <input type="text" onChange={(e) => this.setState({ newTitle: e.currentTarget.value })} defaultValue={this.state.newTitle} />
                        </div>,
                        <DisplayError key={0} fields={['title']} pointing="above" error={this.state.error} />
                        ]
                      : <h1 className="ui header">{this.idea.title}</h1>
                    }

                    <span className="info">
                      Shared <Moment date={this.idea.createdOn} /> by <Gravatar user={this.idea.user} /> <UserName user={this.idea.user} />
                    </span>
                  </div>
                </div>
              </div>

              <span className="subtitle">Description</span>
              {
                this.state.editMode
                ? <div className="ui form">
                    <div className="field">
                      <DisplayError fields={['description']} error={this.state.error} />
                      <Textarea onChange={(e) => this.setState({ newDescription: e.currentTarget.value })} defaultValue={this.state.newDescription} />
                    </div>
                  </div>
                : this.idea.description
                ? <MultiLineText className="description" text={this.idea.description} style="simple" />
                : <p className="description">This idea doesn't have a description.</p>
              }

              <ShowIdeaResponse status={this.idea.status} response={this.idea.response} />

            </div>

            <div className="three wide column">
              {
                this.session.isCollaborator() && [
                  <span key={0} className="subtitle">Actions</span>,
                  this.state.editMode
                    ?
                    <div key={1} className="ui list">
                      <div className="item">
                        <Button className="positive icon fluid text-left" onClick={async () => this.saveChanges()}>
                          <i className="save icon" /> Save
                        </Button>
                      </div>
                      <div className="item">
                        <Button className="icon fluid text-left" onClick={async () => this.setState({ error: undefined, editMode: false })}>
                          <i className="cancel icon" /> Cancel
                        </Button>
                      </div>
                    </div>
                    :
                    <div key={1} className="ui list">
                      <div className="item">
                          <Button className="icon fluid text-left" onClick={async () => this.setState({ editMode: true })}>
                            <i className="edit icon" /> Edit
                          </Button>
                      </div>
                      <div className="item">
                        <ResponseForm idea={this.idea} />
                      </div>
                    </div>
                ]
              }

              <span
                className={`subtitle ${this.session.isCollaborator() && this.tags.length > 0 && 'active'}`}
                onClick={() => this.session.isCollaborator() && this.tags.length > 0 && this.setState({ editTags: !this.state.editTags })}
              >
                Tags
                {this.session.isCollaborator() && this.tags.length > 0 && <i className="setting icon" />}
              </span>

              <div className="ui list tag-list">
                {
                  !this.state.editTags && tagsList
              }
              {
                this.state.editTags &&
                this.tags.map((tag) => (
                  <div key={tag.id} className="item selectable" onClick={async () => this.assignOrUnassignTag(tag)}>
                    <i className={`icon ${this.state.assignedTags.indexOf(tag) >= 0 && 'check'}`} />
                    <ShowTag tag={tag} circular={true}/>
                    <span>{tag.name}</span>
                  </div>
                ))
              }
              </div>
            </div>

            <div className="sixteen wide column">
                <div className="ui comments">
                  <span className="subtitle">Discussion</span>
                  {commentsList}
                  <CommentInput idea={this.idea} />
                </div>
            </div>
          </div>

        </div>
        <Footer />
      </div>
    );
  }
}
