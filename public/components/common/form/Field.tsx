import React from "react";
import { classSet } from "@fider/services";

interface FieldProps {
  className?: string;
  label?: string;
}

export const Field: React.StatelessComponent<FieldProps> = props => {
  return (
    <div
      className={classSet({
        "c-form-field": true,
        [props.className!]: props.className
      })}
    >
      {!!props.label && <label>{props.label}</label>}
      {props.children}
    </div>
  );
};
