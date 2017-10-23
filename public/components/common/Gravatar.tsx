import * as React from 'react';
import md5 = require('md5');
import { getBaseUrl } from '@fider/utils/page';
import { User } from '@fider/models';

interface GravatarProps {
  user?: User;
  name?: string;
  email?: string;
}

export const Gravatar = (props: GravatarProps) => {
  const name = props.name ? props.name : props.user ? props.user.name : '';
  const hash = props.email ? md5(props.email) : props.user ? props.user.gravatar : '';
  const fallback = `${getBaseUrl()}/avatars/50/${name}`;
  const isCollaborator = props.user ? props.user.role >= 2 : false;

  let element: any;
  return <img ref={(e) => element = e }
              onError={() => { element.src = fallback; }}
              className={`fdr-avatar image ${isCollaborator && 'staff'}`}
              title={ name }
              src={ `https://www.gravatar.com/avatar/${hash}?d=${encodeURIComponent(fallback)}` }/>;
};
