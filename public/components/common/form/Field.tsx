import * as React from "react";

interface FieldProps {
  label?: string;
}

export const Field: React.StatelessComponent<FieldProps> = props => {
  return (
    <div className="c-form-field">
      {!!props.label && <label>{props.label}</label>}
      {props.children}
    </div>
  );
};
