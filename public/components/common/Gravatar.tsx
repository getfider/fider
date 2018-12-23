import "./Gravatar.scss";

import React from "react";
import { classSet, Fider } from "@fider/services";
import { isCollaborator, UserRole } from "@fider/models";

interface GravatarProps {
  user?: {
    name: string;
    role?: UserRole;
    avatarURL?: string;
  };
  size?: "small" | "normal" | "large";
}

export const Gravatar = (props: GravatarProps) => {
  const size = props.size || "normal";
  const name = props.user ? props.user.name : "";
  const url = `${Fider.settings.tenantAssetsURL}/avatars/50/0/${encodeURIComponent(name || "?")}`;
  const avatarURL = props.user ? props.user.avatarURL || url : url;

  const className = classSet({
    "c-avatar": true,
    [`m-${size}`]: true,
    "m-staff": props.user && props.user.role && isCollaborator(props.user.role)
  });

  return <img className={className} title={name} src={avatarURL} />;
};
