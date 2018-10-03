import "./ShowError.scss";
import { Fider } from "@fider/services";
import { TenantLogo } from "@fider/components";

import * as React from "react";

interface ShowErrorProps {
  title?: string;
  message?: string;
  error: Error;
  errorInfo: React.ErrorInfo;
  // Must be a function to avoid Fider.isProduction call prior to init
  showError?: () => boolean;
}

const defaultProps: Partial<ShowErrorProps> = {
  title: "Whoops!",
  message: "An unexpected rendering error has occurred.",
  showError: () => !Fider.isProduction()
};

export const ShowError: React.SFC<ShowErrorProps> = ({
  title = "Whoops!",
  message = "An unexpected rendering error has occurred.",
  error,
  errorInfo,
  showError
}: ShowErrorProps) => {
  return (
    <div className="c-show-error container failure-page">
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

ShowError.defaultProps = defaultProps;
