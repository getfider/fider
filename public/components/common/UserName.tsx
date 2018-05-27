import "./UserName.scss";

import * as React from "react";
import { User, UserRole, CurrentUser } from "@fider/models";
import { classSet } from "@fider/services";

interface UserNameProps {
  user: User;
}

export const UserName = (props: UserNameProps) => {
  const className = classSet({
    "c-username": true,
    "m-staff": props.user.role >= UserRole.Collaborator
  });

  return <span className={className}>{props.user.name || "Anonymous"}</span>;
};
