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
  title: "Shoot! Well, this is unexpected…",
  message: "An error has occurred and we're working to fix the problem!",
  showError: () => !Fider.isProduction()
};

export const ErrorPage: React.SFC<ErrorPageProps> = ({
  title = "Shoot! Well, this is unexpected…",
  message = "An error has occurred and we're working to fix the problem!",
  error,
  errorInfo,
  showError
}: ErrorPageProps) => {
  return (
    <div id="p-error" className="container failure-page">
      <TenantLogo size={100} useFiderIfEmpty={true} />
      <h1>{title}</h1>
      <p>{message}</p>
      {Fider.settings && (
        <span>
          Take me back to <a href={Fider.settings.baseURL}>{Fider.settings.baseURL}</a> home page.
        </span>
      )}
      {showError && showError() && (
        <pre className="error">
          {error.toString()}
          {errorInfo.componentStack}
        </pre>
      )}
    </div>
  );
};

ErrorPage.defaultProps = defaultProps;
