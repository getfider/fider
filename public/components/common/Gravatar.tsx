import "./Gravatar.scss";

import * as React from "react";
import { page, classSet } from "@fider/services";
import { User, UserRole, UserStatus, CurrentUser } from "@fider/models";

interface GravatarProps {
  user?: User;
}

export const Gravatar = (props: GravatarProps) => {
  const id = props.user ? props.user.id : 0;
  const name = props.user ? props.user.name : "";
  const url = `${page.getBaseUrl()}/avatars/50/${id}/${encodeURIComponent(name || "?")}`;

  const className = classSet({
    "c-avatar": true,
    "m-staff": props.user && props.user.role >= UserRole.Collaborator
  });

  return <img className={className} title={name} src={url} />;
};
