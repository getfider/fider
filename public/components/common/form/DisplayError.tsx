import * as React from "react";
import { Failure } from "@fider/services";

const arrayToTag = (items: string[]) => {
  return items.map(m => <li key={m}>{m}</li>);
};

interface DisplayErrorProps {
  error?: Failure;
  fields?: string[];
}

export const hasError = (field: string, error?: Failure): boolean => {
  if (error && error.failures) {
    return error.failures.hasOwnProperty(field);
  }
  return false;
};

export const DisplayError = (props: DisplayErrorProps) => {
  if (!props.error) {
    return null;
  }

  let items: JSX.Element[] = [];

  if (props.error.messages && !props.fields) {
    items = arrayToTag(props.error.messages);
  } else if (props.error.failures && props.fields) {
    for (const field of props.fields || Object.keys(props.error.failures)) {
      if (props.error.failures.hasOwnProperty(field)) {
        const tags = arrayToTag(props.error.failures[field]);
        tags.forEach(t => items.push(t));
      }
    }
  }

  return items.length > 0 ? (
    <div className={`c-form-field-error`}>
      <ul>{items}</ul>
    </div>
  ) : null;
};
