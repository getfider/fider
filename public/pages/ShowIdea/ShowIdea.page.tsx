import "./ShowIdea.page.scss";

import * as React from "react";

import { Comment, Post, Tag } from "@fider/models";
import { actions, Failure, Fider } from "@fider/services";

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
  idea: Post;
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

  private saveChanges = async () => {
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
  };

  private setNewTitle = (newTitle: string) => {
    this.setState({ newTitle });
  };

  private setNewDescription = (newDescription: string) => {
    this.setState({ newDescription });
  };

  private cancelEdit = async () => {
    this.setState({ error: undefined, editMode: false });
  };

  private startEdit = async () => {
    this.setState({ editMode: true });
  };

  public render() {
    return (
      <div id="p-show-idea" className="page container">
        <div className="header-col">
          <List>
            <ListItem>
              <SupportCounter post={this.props.idea} />

              <div className="idea-header">
                {this.state.editMode ? (
                  <Form error={this.state.error}>
                    <Input field="title" maxLength={100} value={this.state.newTitle} onChange={this.setNewTitle} />
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
              <TextArea field="description" value={this.state.newDescription} onChange={this.setNewDescription} />
            </Form>
          ) : (
            <MultiLineText
              className="description"
              text={this.props.idea.description || "No description provided."}
              style="simple"
            />
          )}

          <ShowIdeaResponse status={this.props.idea.status} response={this.props.idea.response} />
        </div>

        <div className="action-col">
          {Fider.session.isAuthenticated &&
            Fider.session.user.isCollaborator && [
              <span key={0} className="subtitle">
                Actions
              </span>,
              this.state.editMode ? (
                <List key={1}>
                  <ListItem>
                    <Button className="save" color="positive" fluid={true} onClick={this.saveChanges}>
                      <i className="save icon" /> Save
                    </Button>
                  </ListItem>
                  <ListItem>
                    <Button className="cancel" fluid={true} onClick={this.cancelEdit}>
                      <i className="cancel icon" /> Cancel
                    </Button>
                  </ListItem>
                </List>
              ) : (
                <List key={1}>
                  <ListItem>
                    <Button className="edit" fluid={true} onClick={this.startEdit}>
                      <i className="edit icon" /> Edit
                    </Button>
                  </ListItem>
                  <ListItem>
                    <ResponseForm idea={this.props.idea} />
                  </ListItem>
                </List>
              )
            ]}

          <TagsPanel idea={this.props.idea} tags={this.props.tags} />
          <NotificationsPanel idea={this.props.idea} subscribed={this.props.subscribed} />
          <ModerationPanel idea={this.props.idea} />
        </div>

        <DiscussionPanel idea={this.props.idea} comments={this.props.comments} />
      </div>
    );
  }
}
