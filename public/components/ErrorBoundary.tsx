import React from "react"
import { ErrorPage } from "@fider/pages/Error/Error.page"
import { FiderContext } from "@fider/services"

interface ErrorBoundaryProps {
  onError?: (err: Error) => void
}

interface ErrorBoundaryState {
  error?: Error
  errorInfo?: React.ErrorInfo
}

export class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: any) {
    super(props)

    this.state = {
      error: undefined,
      errorInfo: undefined,
    }
  }

  public componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    const onError = this.props.onError
    if (onError) {
      onError(error)
    }

    this.setState({
      error,
      errorInfo,
    })
  }

  public render() {
    const { error, errorInfo } = this.state

    if (error && errorInfo) {
      return <FiderContext.Consumer>{(fider) => <ErrorPage error={error} errorInfo={errorInfo} showDetails={!fider.isProduction()} />}</FiderContext.Consumer>
    } else {
      return this.props.children
    }
  }
}
