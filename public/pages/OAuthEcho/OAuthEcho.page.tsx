import "./OAuthEcho.page.scss";

import * as React from "react";
import { navigator } from "@fider/services";
import { Segments, Segment } from "@fider/components";

interface OAuthEchoPageProps {
  statusCode: number;
  response: string;
  profile: {
    id: string;
    name: string;
    email: string;
  };
}

const ok = <i className="check circle icon" />;
const error = <i className="times circle icon" />;
const warn = <i className="exclamation triangle icon" />;

export class OAuthEchoPage extends React.Component<OAuthEchoPageProps, {}> {
  public componentDidMount() {
    navigator.replaceState("/");
  }

  public render() {
    const idOk = this.props.profile && this.props.profile.id !== "";
    const nameOk = this.props.profile && this.props.profile.name !== "Anonymous";
    const emailOk = this.props.profile && this.props.profile.email !== "";

    return (
      <div id="p-oauth-echo" className="page container">
        <h5>Status Code: {this.props.statusCode}</h5>
        <h5>Raw Response</h5>
        <pre>{JSON.stringify(JSON.parse(this.props.response), null, "  ")}</pre>
        <h5>Parsed Profile</h5>
        <Segments>
          <Segment>
            <p>
              {idOk ? ok : error}
              <strong>ID:</strong> {this.props.profile && this.props.profile.id}
              {!idOk && (
                <p className="info">ID is required. If not found, users will see an error during sign in process.</p>
              )}
            </p>
          </Segment>
          <Segment>
            <p>
              {nameOk ? ok : warn}
              <strong>Name:</strong> {this.props.profile && this.props.profile.name}
            </p>
            {!nameOk && (
              <p className="info">
                Name is required, if not found we'll use <strong>Anonymous</strong> as the name of every new user.
              </p>
            )}
          </Segment>
          <Segment>
            <p>
              {emailOk ? ok : warn}
              <strong>Email:</strong> {this.props.profile && this.props.profile.email}
            </p>
            {!emailOk && (
              <p className="info">
                Email is not required, but highly recommended. If invalid or not found, new users won't be able to
                receive notifications.
              </p>
            )}
          </Segment>
        </Segments>
      </div>
    );
  }
}
