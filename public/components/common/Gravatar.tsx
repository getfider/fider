import * as React from 'react';
import { getBaseUrl } from '@fider/utils/page';
import { User } from '@fider/models';

interface GravatarProps {
  user?: User;
  name?: string;
  email?: string;
}

export const Gravatar = (props: GravatarProps) => {
  const name = props.name ? props.name : props.user ? props.user.name : '_';
  const id = props.user ? props.user.id : 0;
  const queryString = props.email ? `?e=${encodeURIComponent(props.email)}` : '';

  const url = `${getBaseUrl()}/avatars/50/${id}/${encodeURIComponent(name)}${queryString}`;
  const isCollaborator = props.user ? props.user.role >= 2 : false;

  let element: any;
  return (
    <img
      ref={(e) => element = e}
      className={`fdr-avatar image ${isCollaborator && 'staff'}`}
      title={name}
      src={url}
    />
  );
};
