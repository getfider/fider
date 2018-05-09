import "./Gravatar.scss";

import * as React from "react";
import { page, classSet } from "@fider/services";
import { User, UserRole } from "@fider/models";

interface GravatarProps {
  user?: User;
}

export const Gravatar = (props: GravatarProps) => {
  const name = props.user ? props.user.name : "_";
  const id = props.user ? props.user.id : 0;

  const url = `${page.getBaseUrl()}/avatars/50/${id}/${encodeURIComponent(name)}`;
  const isCollaborator = props.user ? props.user.role >= UserRole.Collaborator : false;

  const className = classSet({
    "c-avatar": true,
    "gc-staff": isCollaborator
  });

  return <img className={className} title={name} src={url} />;
};
