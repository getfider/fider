import "./Error.page.scss";

import React from "react";
import { Fider } from "@fider/services";
import { TenantLogo } from "@fider/components";

interface ErrorPageProps {
  title?: string;
  message?: string;
  error: Error;
  errorInfo: React.ErrorInfo;
  // Must be a function to avoid Fider.isProduction call prior to init
  showError?: () => boolean;
}

const defaultProps: Partial<ErrorPageProps> = {
  title: "Whoops!",
  message: "An unexpected rendering error has occurred.",
  showError: () => !Fider.isProduction()
};

export const ErrorPage: React.SFC<ErrorPageProps> = ({
  title = "Whoops!",
  message = "An unexpected rendering error has occurred.",
  error,
  errorInfo,
  showError
}: ErrorPageProps) => {
  return (
    <div id="p-error" className="container failure-page">
      <TenantLogo size={100} />
      <div className="content">
        <h2>{title}</h2>
        <p>{message}</p>
        {showError &&
          showError() && (
            <pre className="error">
              {error.toString()}
              {errorInfo.componentStack}
            </pre>
          )}
      </div>
    </div>
  );
};

ErrorPage.defaultProps = defaultProps;
