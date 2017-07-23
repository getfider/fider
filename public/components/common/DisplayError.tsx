import * as React from 'react';
import { Session, Failure } from '@fider/services';

const arrayToTag = (items: string[]) => {
  return items.map((m) => <li>{m}</li>);
};

interface DisplayErrorProps {
  error?: Failure;
  fields?: string[];
}

export const DisplayError = (props: DisplayErrorProps) => {
  if (!props.error) {
    return <div></div>;
  }

  let items: JSX.Element[] = [];

  if (props.fields && props.error.failures) {
    for (const field of props.fields) {
      if (props.error.failures.hasOwnProperty(field)) {
        const tags = arrayToTag(props.error.failures[field]);
        tags.forEach((t) => items.push(t));
      }
    }
  } else if (props.error.messages) {
    items = arrayToTag(props.error.messages);
  }

  return items.length > 0 ? <div className="display-error ui pointing below red basic label">
            <ul>{ items }</ul>
         </div> : null;
};
