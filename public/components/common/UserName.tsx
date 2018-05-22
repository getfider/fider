import "./UserName.scss";

import * as React from "react";
import { User, UserRole } from "@fider/models";
import { classSet } from "@fider/services";

interface UserNameProps {
  user: User;
}

export const UserName = (props: UserNameProps) => {
  const isCollaborator = props.user.role >= UserRole.Collaborator;

  const className = classSet({
    "c-username": true,
    "m-staff": isCollaborator
  });

  return <span className={className}>{props.user.name}</span>;
};
