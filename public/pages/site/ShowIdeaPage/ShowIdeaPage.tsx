import * as React from 'react';

import { inject, injectables } from '@fider/di';
import { CurrentUser, Comment, Idea, Tag } from '@fider/models';
import { Session, IdeaService, Failure } from '@fider/services';

import { setTitle } from '@fider/utils/page';

import { TagsPanel, CommentInput, CommentList, ResponseForm } from './';
import { SupportCounter, ShowIdeaResponse, DisplayError, Button, Textarea, UserName, Gravatar, Moment, Form, MultiLineText, Footer, Header } from '@fider/components';

import './ShowIdeaPage.scss';

interface ShowIdeaPageState {
  editMode: boolean;
  newTitle: string;
  newDescription: string;
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

  constructor(props: {}) {
    super(props);

    this.user = this.session.getCurrentUser();
    this.idea = this.session.get<Idea>('idea');
    this.comments = this.session.getArray<Comment>('comments');
    this.tags = this.session.getArray<Tag>('tags');

    this.state = {
      editMode: false,
      newTitle: this.idea.title,
      newDescription: this.idea.description,
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

  public render() {
    return (
      <div>
        <Header />
        <div className="page ui container">
          <div className="ui stackable vertically padded grid container">
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
            <div className="three wide column action-col">
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

              <TagsPanel idea={this.idea} tags={this.tags} />
            </div>
            <div className="thirteen wide column">
                <div className="ui comments">
                  <span className="subtitle">Discussion</span>
                  <CommentList comments={this.comments} />
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
