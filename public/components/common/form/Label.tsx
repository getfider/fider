import * as React from "react";

export const Label: React.StatelessComponent<{}> = props => {
  return (
    <div className="c-form-field">
      <label>{props.children}</label>
    </div>
  );
};
