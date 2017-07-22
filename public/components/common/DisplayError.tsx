import * as React from 'react';
import { Session, Failure } from '@fider/services';

const arrayToTag = (items: string[]) => {
  return items.map((m) => <li>{m}</li>);
};

export const DisplayError = (props: {error?: Failure}) => {
  if (!props.error) {
    return <div></div>;
  }

  const items = props.error.messages ? arrayToTag(props.error.messages) : [];

  for (const field in props.error.failures) {
    if (props.error.failures.hasOwnProperty(field)) {
      const tags = arrayToTag(props.error.failures[field]);
      tags.forEach((t) => items.push(t));
    }
  }

  return <div className="display-error ui pointing below red basic label">
            { items && <ul>{ items }</ul> }
         </div>;
};
