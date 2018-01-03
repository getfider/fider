import * as React from 'react';
import { User, UserRole } from '@fider/models';

interface UserNameProps {
  user: User;
}

export const UserName = (props: UserNameProps) => {
  const isCollaborator = props.user.role >= UserRole.Collaborator;
  return <span className={`name ${isCollaborator ? 'staff' : ''}`}>{props.user.name}</span>;
};
