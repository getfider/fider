import "./Home.page.scss";

import * as React from "react";
import { Post, Tag, PostStatus } from "@fider/models";
import { MultiLineText } from "@fider/components";
import { IdeaInput, ListIdeas, IdeasContainer } from "./";
import { actions, Fider } from "@fider/services";
import { SimilarIdeas } from "./components/SimilarIdeas";

export interface HomePageProps {
  ideas: Post[];
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
        <i className="icon lightbulb outline" aria-hidden="true" />
      </p>
      <p>It's lonely out here. Start by sharing a suggestion!</p>
    </div>
  );
};

const defaultWelcomeMessage = `## Welcome to our feedback site!

We'd love to hear what you're thinking about. What can we do better? This is the place for you to vote, discuss and share ideas.`;

export class HomePage extends React.Component<HomePageProps, HomePageState> {
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
            <IdeaInput
              placeholder={Fider.session.tenant.invitation || "Enter your idea here..."}
              onTitleChanged={this.setTitle}
            />
          </div>
          <div className="col-md-8">
            {this.isLonely() ? (
              <Lonely />
            ) : this.state.title ? (
              <SimilarIdeas title={this.state.title} tags={this.props.tags} />
            ) : (
              <IdeasContainer
                ideas={this.props.ideas}
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
