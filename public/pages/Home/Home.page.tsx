import "./Home.page.scss";

import React from "react";
import { Post, Tag, PostStatus } from "@fider/models";
import { MultiLineText } from "@fider/components";
import { Fider } from "@fider/services";
import { SimilarPosts } from "./components/SimilarPosts";
import { FaRegLightbulb } from "react-icons/fa";
import { PostInput } from "./components/PostInput";
import { PostsContainer } from "./components/PostsContainer";

export interface HomePageProps {
  posts: Post[];
  tags: Tag[];
  countPerStatus: { [key: string]: number };
}

export interface HomePageState {
  title: string;
}

const Lonely = () => {
  return (
    <div className="center">
      <p>
        <FaRegLightbulb />
      </p>
      <p>It's lonely out here. Start by sharing a suggestion!</p>
    </div>
  );
};

const defaultWelcomeMessage = `We'd love to hear what you're thinking about. 

What can we do better? This is the place for you to vote, discuss and share ideas.`;

export default class HomePage extends React.Component<HomePageProps, HomePageState> {
  constructor(props: HomePageProps) {
    super(props);
    this.state = {
      title: ""
    };
  }

  private isLonely(): boolean {
    const len = Object.keys(this.props.countPerStatus).length;
    if (len === 0) {
      return true;
    }

    if (len === 1 && PostStatus.Deleted.value in this.props.countPerStatus) {
      return true;
    }

    return false;
  }

  private setTitle = async (title: string) => {
    this.setState({ title });
  };

  public render() {
    return (
      <div id="p-home" className="page container">
        <div className="row">
          <div className="col-md-4">
            <MultiLineText
              className="welcome-message"
              text={Fider.session.tenant.welcomeMessage || defaultWelcomeMessage}
              style="full"
            />
            <PostInput
              placeholder={Fider.session.tenant.invitation || "Enter your suggestion here..."}
              onTitleChanged={this.setTitle}
            />
          </div>
          <div className="col-md-8">
            {this.isLonely() ? (
              <Lonely />
            ) : this.state.title ? (
              <SimilarPosts title={this.state.title} tags={this.props.tags} />
            ) : (
              <PostsContainer
                posts={this.props.posts}
                tags={this.props.tags}
                countPerStatus={this.props.countPerStatus}
              />
            )}
          </div>
        </div>
      </div>
    );
  }
}
