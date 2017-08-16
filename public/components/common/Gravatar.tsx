import * as React from 'react';
import md5 = require('md5');
import { getBaseUrl } from '@fider/utils/page';

interface GravatarProps {
  name: string;
  hash?: string;
  email?: string;
}

export const Gravatar = (props: GravatarProps) => {
  let element: any;
  const error = () => {
    element.src = `${getBaseUrl()}/avatars/50/${props.name}`;
  };
  const hash = props.email ? md5(props.email) : props.hash || '';

  return (hash && props.name) ?
    <img ref={(e) => element = e} onError={error} className="ui avatar image" src={ `https://www.gravatar.com/avatar/${hash}?d=404` }/> : 
    <div />;
};
