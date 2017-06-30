import * as React from 'react';
import { Session, Failure } from '@fider/services';

export const DisplayError = (props: {error?: Failure}) => {
  if (!props.error) {
    return <div></div>;
  }

  return <div className="ui negative message">
            <p>{ props.error.message }</p>
         </div>;
};
