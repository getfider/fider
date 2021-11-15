import React from "react"
import { navigator } from "@fider/services"
import { Icon } from "@fider/components"

import IconXCircle from "@fider/assets/images/heroicons-x-circle.svg"
import IconCheckCircle from "@fider/assets/images/heroicons-check-circle.svg"
import IconExclamation from "@fider/assets/images/heroicons-exclamation.svg"
import { HStack, VStack } from "@fider/components/layout"

interface OAuthEchoPageProps {
  err: string | undefined
  body: string
  profile: {
    id: string
    name: string
    email: string
  }
}

const ok = <Icon sprite={IconCheckCircle} className="h-4 text-green-500" />
const error = <Icon sprite={IconXCircle} className="h-4 text-red-500" />
const warn = <Icon sprite={IconExclamation} className="h-4 text-yellow-500" />

export default class OAuthEchoPage extends React.Component<OAuthEchoPageProps, any> {
  public componentDidMount() {
    navigator.replaceState("/")
  }

  private renderError() {
    return (
      <>
        <h5 className="text-display">Error</h5>
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
        <h5 className="text-display mb-2">Raw Body</h5>
        <pre>{responseBody}</pre>
        <h5 className="text-display mb-2 mt-8">Parsed Profile</h5>
        <VStack divide={true} spacing={2}>
          <VStack>
            <HStack>
              {idOk ? ok : error}
              <strong>ID:</strong> <span>{this.props.profile && this.props.profile.id}</span>
            </HStack>
            {!idOk && <span className="text-muted">ID is required. If not found, users will see an error during sign in process.</span>}
          </VStack>
          <VStack>
            <HStack>
              {nameOk ? ok : warn}
              <strong>Name:</strong> <span>{this.props.profile && this.props.profile.name}</span>
            </HStack>
            {!nameOk && (
              <span className="text-muted">
                Name is required, if not found we&apos;ll use <strong>Anonymous</strong> as the name of every new user.
              </span>
            )}
          </VStack>
          <VStack>
            <HStack>
              {emailOk ? ok : warn}
              <strong>Email:</strong> {this.props.profile && this.props.profile.email}
            </HStack>
            {!emailOk && (
              <span className="text-muted">
                Email is not required, but highly recommended. If invalid or not found, new users won&apos;t receive notifications.
              </span>
            )}
          </VStack>
        </VStack>
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
