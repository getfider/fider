import "./ShowIdea.page.scss";

import * as React from "react";

import { CurrentUser, Comment, Idea, Tag } from "@fider/models";
import { actions, Failure } from "@fider/services";

import { TagsPanel, DiscussionPanel, ResponseForm, NotificationsPanel, ModerationPanel } from "./";
import {
  SupportCounter,
  ShowIdeaResponse,
  DisplayError,
  Button,
  UserName,
  Gravatar,
  Moment,
  MultiLineText,
  List,
  ListItem,
  Input,
  Form,
  TextArea
} from "@fider/components";

interface ShowIdeaPageProps {
  user?: CurrentUser;
  idea: Idea;
  subscribed: boolean;
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
      newDescription: this.props.idea.description
    };
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
      <div id="p-show-idea" className="page container">
        <div className="header-col">
          <List>
            <ListItem>
              <SupportCounter user={this.props.user} idea={this.props.idea} />

              <div className="idea-header">
                {this.state.editMode ? (
                  <Form error={this.state.error}>
                    <Input
                      field="title"
                      maxLength={100}
                      value={this.state.newTitle}
                      onChange={newTitle => this.setState({ newTitle })}
                    />
                  </Form>
                ) : (
                  <h1>{this.props.idea.title}</h1>
                )}

                <span className="info">
                  Shared <Moment date={this.props.idea.createdOn} /> by <Gravatar user={this.props.idea.user} />{" "}
                  <UserName user={this.props.idea.user} />
                </span>
              </div>
            </ListItem>
          </List>

          <span className="subtitle">Description</span>
          {this.state.editMode ? (
            <Form error={this.state.error}>
              <TextArea
                field="description"
                value={this.state.newDescription}
                onChange={newDescription => this.setState({ newDescription })}
              />
            </Form>
          ) : this.props.idea.description ? (
            <MultiLineText className="description" text={this.props.idea.description} style="simple" />
          ) : (
            <p className="description">This idea doesn't have a description.</p>
          )}

          <ShowIdeaResponse status={this.props.idea.status} response={this.props.idea.response} />
        </div>

        <div className="action-col">
          {this.props.user &&
            this.props.user.isCollaborator && [
              <span key={0} className="subtitle">
                Actions
              </span>,
              this.state.editMode ? (
                <List key={1}>
                  <ListItem>
                    <Button color="positive" fluid={true} onClick={async () => this.saveChanges()}>
                      <i className="save icon" /> Save
                    </Button>
                  </ListItem>
                  <ListItem>
                    <Button fluid={true} onClick={async () => this.setState({ error: undefined, editMode: false })}>
                      <i className="cancel icon" /> Cancel
                    </Button>
                  </ListItem>
                </List>
              ) : (
                <List key={1}>
                  <ListItem>
                    <Button fluid={true} onClick={async () => this.setState({ editMode: true })}>
                      <i className="edit icon" /> Edit
                    </Button>
                  </ListItem>
                  <ListItem>
                    <ResponseForm idea={this.props.idea} />
                  </ListItem>
                </List>
              )
            ]}

          <TagsPanel user={this.props.user} idea={this.props.idea} tags={this.props.tags} />
          <NotificationsPanel user={this.props.user} idea={this.props.idea} subscribed={this.props.subscribed} />
          <ModerationPanel user={this.props.user} idea={this.props.idea} />
        </div>

        <DiscussionPanel user={this.props.user} idea={this.props.idea} comments={this.props.comments} />
      </div>
    );
  }
}
