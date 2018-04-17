import "./Home.page.scss";

import * as React from "react";
import { Idea, Tag, IdeaStatus, CurrentUser, Tenant } from "@fider/models";
import { MultiLineText } from "@fider/components";
import { IdeaInput, ListIdeas, IdeasContainer } from "./";
import { page, actions } from "@fider/services";

export interface HomePageProps {
  user?: CurrentUser;
  tenant: Tenant;
  ideas: Idea[];
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
        <i className="icon lightbulb" aria-hidden="true" />
      </p>
      <p>It's lonely out here. Start by sharing an idea!</p>
    </div>
  );
};

const defaultWelcomeMessage = `## Welcome to our feedback forum!

We'd love to hear what you're thinking about. What can we do better? This is the place for you to vote, discuss and share ideas.`;

export class HomePage extends React.Component<HomePageProps, HomePageState> {
  constructor(props: HomePageProps) {
    super(props);
    this.state = {
      title: ""
    };
  }

  public render() {
    return (
      <div id="p-home" className="page ui container">
        <div className="ui grid stackable">
          <div className="six wide column">
            <MultiLineText
              className="welcome-message"
              text={this.props.tenant.welcomeMessage || defaultWelcomeMessage}
              style="full"
            />
            <IdeaInput
              user={this.props.user}
              placeholder={this.props.tenant.invitation || "I suggest you..."}
              onTitleChanged={title => this.setState({ title })}
            />
          </div>
          <div className="ten wide column">
            {Object.keys(this.props.countPerStatus).length === 0 ? (
              <Lonely />
            ) : (
              <IdeasContainer
                user={this.props.user}
                ideas={this.props.ideas}
                tags={this.props.tags}
                countPerStatus={this.props.countPerStatus}
                newIdeaTitle={this.state.title}
              />
            )}
          </div>
        </div>
      </div>
    );
  }
}
