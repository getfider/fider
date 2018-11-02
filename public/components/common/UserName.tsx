import "./UserName.scss";

import React from "react";
import { User, isCollaborator } from "@fider/models";
import { classSet } from "@fider/services";

interface UserNameProps {
  user: User;
}

export const UserName = (props: UserNameProps) => {
  const className = classSet({
    "c-username": true,
    "m-staff": isCollaborator(props.user.role)
  });

  return <span className={className}>{props.user.name || "Anonymous"}</span>;
};
