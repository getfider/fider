import * as React from 'react';

import { CurrentUser, Comment, Idea, Tag } from '@fider/models';
import { page, actions, Failure } from '@fider/services';

import { TagsPanel, CommentInput, CommentList, ResponseForm } from './';
import { SupportCounter, ShowIdeaResponse, DisplayError, Button, Textarea, UserName, Gravatar, Moment, MultiLineText } from '@fider/components';

import './ShowIdeaPage.scss';

interface ShowIdeaPageProps {
  user?: CurrentUser;
  idea: Idea;
  comments: Comment[];
  tags: Tag[];
}

interface ShowIdeaPageState {
  editMode: boolean;
  newTitle: string;
  newDescription: string;
  error?: Failure;
}

export class ShowIdeaPage extends React.Component<ShowIdeaPageProps, ShowIdeaPageState> {

  constructor(props: ShowIdeaPageProps) {
    super(props);

    this.state = {
      editMode: false,
      newTitle: this.props.idea.title,
      newDescription: this.props.idea.description,
    };

    page.setTitle(`${this.props.idea.title} · ${document.title}`);
  }

  private async saveChanges() {
    const result = await actions.updateIdea(this.props.idea.number, this.state.newTitle, this.state.newDescription);
    if (result.ok) {
      this.setState({
        error: undefined,
        editMode: false
      });
      this.props.idea.title = this.state.newTitle;
      this.props.idea.description = this.state.newDescription;
      this.forceUpdate();
    } else {
      this.setState({
        error: result.error
      });
    }
  }

  public render() {
    return (
      <div className="page ui container">
        <div className="ui stackable vertically padded grid container">
          <div className="thirteen wide column">
            <div className="ui items unstackable">
              <div className="item">
                <SupportCounter user={this.props.user} idea={this.props.idea} />

                <div className="idea-header">
                  { this.state.editMode
                    ? [
                      <div key={1} className="ui input huge fluid">
                        <input type="text" maxLength={100} onChange={(e) => this.setState({ newTitle: e.currentTarget.value })} defaultValue={this.state.newTitle} />
                      </div>,
                      <DisplayError key={0} fields={['title']} pointing="above" error={this.state.error} />
                    ]
                    : <h1 className="ui header">{this.props.idea.title}</h1>
                  }

                  <span className="info">
                    Shared <Moment date={this.props.idea.createdOn} /> by <Gravatar user={this.props.idea.user} /> <UserName user={this.props.idea.user} />
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
              : this.props.idea.description
              ? <MultiLineText className="description" text={this.props.idea.description} style="simple" />
              : <p className="description">This idea doesn't have a description.</p>
            }

            <ShowIdeaResponse status={this.props.idea.status} response={this.props.idea.response} />

          </div>
          <div className="three wide column action-col">
            {
              this.props.user && this.props.user.isCollaborator && [
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
                      <ResponseForm idea={this.props.idea} />
                    </div>
                  </div>
              ]
            }

            <TagsPanel user={this.props.user} idea={this.props.idea} tags={this.props.tags} />
          </div>
          <div className="thirteen wide column">
              <div className="ui comments">
                <span className="subtitle">Discussion</span>
                <CommentList comments={this.props.comments} />
                <CommentInput user={this.props.user} idea={this.props.idea} />
              </div>
          </div>

        </div>

      </div>
    );
  }
}
