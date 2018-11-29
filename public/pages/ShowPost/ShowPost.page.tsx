import "./ShowPost.page.scss";

import React from "react";

import { Comment, Post, Tag, Vote } from "@fider/models";
import { actions, Failure, Fider } from "@fider/services";

import {
  VoteCounter,
  ShowPostResponse,
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
import { FaSave, FaTimes, FaEdit } from "react-icons/fa";
import { ResponseForm } from "./components/ResponseForm";
import { TagsPanel } from "./components/TagsPanel";
import { NotificationsPanel } from "./components/NotificationsPanel";
import { ModerationPanel } from "./components/ModerationPanel";
import { DiscussionPanel } from "./components/DiscussionPanel";
import { VotesPanel } from "./components/VotesPanel";

interface ShowPostPageProps {
  post: Post;
  subscribed: boolean;
  comments: Comment[];
  tags: Tag[];
  votes: Vote[];
}

interface ShowPostPageState {
  editMode: boolean;
  newTitle: string;
  newDescription: string;
  error?: Failure;
}

export default class ShowPostPage extends React.Component<ShowPostPageProps, ShowPostPageState> {
  constructor(props: ShowPostPageProps) {
    super(props);

    this.state = {
      editMode: false,
      newTitle: this.props.post.title,
      newDescription: this.props.post.description
    };
  }

  private saveChanges = async () => {
    const result = await actions.updatePost(this.props.post.number, this.state.newTitle, this.state.newDescription);
    if (result.ok) {
      this.setState({
        error: undefined,
        editMode: false
      });
      this.props.post.title = this.state.newTitle;
      this.props.post.description = this.state.newDescription;
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
      <div id="p-show-post" className="page container">
        <div className="header-col">
          <List>
            <ListItem>
              <VoteCounter post={this.props.post} />

              <div className="post-header">
                {this.state.editMode ? (
                  <Form error={this.state.error}>
                    <Input field="title" maxLength={100} value={this.state.newTitle} onChange={this.setNewTitle} />
                  </Form>
                ) : (
                  <h1>{this.props.post.title}</h1>
                )}

                <span className="info">
                  <Moment date={this.props.post.createdAt} /> &middot; <Gravatar user={this.props.post.user} />{" "}
                  <UserName user={this.props.post.user} />
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
            <MultiLineText className="description" text={this.props.post.description} style="simple" />
          )}

          <ShowPostResponse showUser={true} status={this.props.post.status} response={this.props.post.response} />
        </div>

        <div className="action-col">
          <VotesPanel post={this.props.post} votes={this.props.votes} />

          {Fider.session.isAuthenticated &&
            Fider.session.user.isCollaborator && [
              <span key={0} className="subtitle">
                Actions
              </span>,
              this.state.editMode ? (
                <List key={1}>
                  <ListItem>
                    <Button className="save" color="positive" fluid={true} onClick={this.saveChanges}>
                      <FaSave /> Save
                    </Button>
                  </ListItem>
                  <ListItem>
                    <Button className="cancel" fluid={true} onClick={this.cancelEdit}>
                      <FaTimes /> Cancel
                    </Button>
                  </ListItem>
                </List>
              ) : (
                <List key={1}>
                  <ListItem>
                    <Button className="edit" fluid={true} onClick={this.startEdit}>
                      <FaEdit /> Edit
                    </Button>
                  </ListItem>
                  <ListItem>
                    <ResponseForm post={this.props.post} />
                  </ListItem>
                </List>
              )
            ]}

          <TagsPanel post={this.props.post} tags={this.props.tags} />
          <NotificationsPanel post={this.props.post} subscribed={this.props.subscribed} />
          <ModerationPanel post={this.props.post} />
        </div>

        <DiscussionPanel post={this.props.post} comments={this.props.comments} />
      </div>
    );
  }
}
