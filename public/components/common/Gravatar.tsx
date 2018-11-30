import "./Gravatar.scss";

import React from "react";
import { classSet, Fider } from "@fider/services";
import { isCollaborator, UserRole } from "@fider/models";

interface GravatarProps {
  user?: {
    id: number;
    name: string;
    role?: UserRole;
  };
  size?: "small" | "normal" | "large";
}

export const Gravatar = (props: GravatarProps) => {
  const size = props.size || "normal";
  const id = props.user ? props.user.id : 0;
  const name = props.user ? props.user.name : "";
  const url = `${Fider.settings.tenantAssetsURL}/avatars/50/${id}/${encodeURIComponent(name || "?")}`;

  const className = classSet({
    "c-avatar": true,
    [`m-${size}`]: true,
    "m-staff": props.user && props.user.role && isCollaborator(props.user.role)
  });

  return <img className={className} title={name} src={url} />;
};
