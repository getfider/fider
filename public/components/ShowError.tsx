import "./ShowError.scss";

import * as React from "react";

interface ShowErrorProps {
  title?: string;
  message?: string;
}

export const ShowError = ({
  title = "Whoops!",
  message = "An unexpected rendering error has occurred."
}: ShowErrorProps) => {
  return (
    <div className="c-show-error">
      <img className="logo" src="https://getfider.com/images/logo-100x100.png" />
      <div className="content">
        <h2>{title}</h2>
        <p>{message}</p>
      </div>
    </div>
  );
};
