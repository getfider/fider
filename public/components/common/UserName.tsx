import * as React from 'react';
import { User } from '@fider/models';

interface UserNameProps {
  user: User;
}

export const UserName = (props: UserNameProps) => {
  const isCollaborator = props.user.role >= 2;
  return <span className={`name ${isCollaborator ? 'staff' : ''}`}>{ props.user.name }</span>;
};
