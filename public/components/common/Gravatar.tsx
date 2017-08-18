import * as React from 'react';
import md5 = require('md5');
import { getBaseUrl } from '@fider/utils/page';

interface GravatarProps {
  name: string;
  hash?: string;
  email?: string;
}

export const Gravatar = (props: GravatarProps) => {
  const fallback = `${getBaseUrl()}/avatars/50/${props.name}`;
  const hash = props.email ? md5(props.email) : props.hash || '';

  let element: any;
  return (hash && props.name) ?
    <img ref={(e) => element = e }
         onError={() => { element.src = fallback; }}
         className="ui avatar image"
         title={props.name}
         src={ `https://www.gravatar.com/avatar/${hash}?d=${encodeURIComponent(fallback)}` }/> :
    <div />;
};
