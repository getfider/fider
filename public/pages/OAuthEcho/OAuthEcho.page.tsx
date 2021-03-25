import "./OAuthEcho.page.scss"

import React from "react"
import { navigator } from "@fider/services"
import { Segments, Segment } from "@fider/components"
import { FaCheckCircle, FaTimesCircle, FaExclamationTriangle } from "react-icons/fa"

interface OAuthEchoPageProps {
  err: string | undefined
  body: string
  profile: {
    id: string
    name: string
    email: string
  }
}

const ok = <FaCheckCircle className="check" />
const error = <FaTimesCircle className="error" />
const warn = <FaExclamationTriangle className="warn" />

export default class OAuthEchoPage extends React.Component<OAuthEchoPageProps, any> {
  public componentDidMount() {
    navigator.replaceState("/")
  }

  private renderError() {
    return (
      <>
        <h5>Error</h5>
        <pre>{this.props.err}</pre>
      </>
    )
  }

  private renderParseResult() {
    const idOk = this.props.profile && this.props.profile.id !== ""
    const nameOk = this.props.profile && this.props.profile.name !== "Anonymous"
    const emailOk = this.props.profile && this.props.profile.email !== ""

    let responseBody = ""
    try {
      responseBody = JSON.stringify(JSON.parse(this.props.body), null, "  ")
    } catch {
      responseBody = this.props.body
    }

    return (
      <>
        <h5>Raw Body</h5>
        <pre>{responseBody}</pre>
        <h5>Parsed Profile</h5>
        <Segments>
          <Segment>
            <p>
              {idOk ? ok : error}
              <strong>ID:</strong> {this.props.profile && this.props.profile.id}
              {!idOk && <p className="info">ID is required. If not found, users will see an error during sign in process.</p>}
            </p>
          </Segment>
          <Segment>
            <p>
              {nameOk ? ok : warn}
              <strong>Name:</strong> {this.props.profile && this.props.profile.name}
              {!nameOk && (
                <p className="info">
                  Name is required, if not found we&apos;ll use <strong>Anonymous</strong> as the name of every new user.
                </p>
              )}
            </p>
          </Segment>
          <Segment>
            <p>
              {emailOk ? ok : warn}
              <strong>Email:</strong> {this.props.profile && this.props.profile.email}
              {!emailOk && (
                <p className="info">
                  Email is not required, but highly recommended. If invalid or not found, new users won&apos;t be able to receive notifications.
                </p>
              )}
            </p>
          </Segment>
        </Segments>
      </>
    )
  }

  public render() {
    return (
      <div id="p-oauth-echo" className="page container">
        {this.props.err ? this.renderError() : this.renderParseResult()}
      </div>
    )
  }
}
