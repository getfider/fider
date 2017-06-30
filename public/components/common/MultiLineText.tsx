import * as React from 'react';

export const MultiLineText = (props: {text?: string}) => {
  if (!props.text) {
    return <p></p>;
  }

  return <div>{props.text.split('\n').map((item, i) =>
   <span>{item}<br/></span>
  )}</div>;
};
