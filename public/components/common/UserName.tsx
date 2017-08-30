import * as React from 'react';
import { User } from '@fider/models';

interface UserNameProps {
  user: User;
}

export const UserName = (props: UserNameProps) => {
  const isStaff = props.user && props.user.role >= 2;
  return <span className={`name ${isStaff ? 'staff' : ''}`}>{ props.user.name }</span>;
};
