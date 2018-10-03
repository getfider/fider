import * as React from "react";
import { ShowError } from "./ShowError";

interface ErrorBoundaryProps {
  showError?: boolean;
  onError?: (err: Error) => void;
}

interface ErrorBoundaryState {
  error?: Error;
}

export class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: any) {
    super(props);

    this.state = {
      error: undefined
    };
  }

  public componentDidCatch(error: Error, errorinfo: React.ErrorInfo) {
    const onError = this.props.onError;
    if (onError) {
      onError(error);
    }

    this.setState({
      error
    });
  }

  public render() {
    const showError = this.props.showError;
    const error = this.state.error;

    if (error) {
      const message = showError ? error.message : undefined;
      return <ShowError message={message} />;
    } else {
      return this.props.children;
    }
  }
}
