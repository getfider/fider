import * as React from "react";
import { ShowError } from "./ShowError";

interface ErrorBoundaryState {
  error?: Error;
  errorinfo?: React.ErrorInfo;
}

export class ErrorBoundary extends React.Component<{}, ErrorBoundaryState> {
  constructor(props: any) {
    super(props);

    this.state = {
      error: undefined,
      errorinfo: undefined
    };
  }

  public componentDidCatch(error: Error, errorinfo: React.ErrorInfo) {
    this.setState({
      error,
      errorinfo
    });
  }

  public render() {
    const { error, errorinfo } = this.state;

    if (error) {
      return <ShowError />;
    }

    return this.props.children;
  }
}
