import "./Heading.scss";

import React from "react";
import { classSet } from "@fider/services";

interface HeadingLogo {
  title: string;
  dividing?: boolean;
  size?: "normal" | "small";
  icon?: string;
  subtitle?: string;
  className?: string;
}

export const Heading = (props: HeadingLogo) => {
  const size = props.size || "normal";
  const level = size === "normal" ? 2 : 3;
  const Tag = `h${level}`;
  const className = classSet({
    "c-heading": true,
    "m-dividing": props.dividing || false,
    [`m-${size}`]: true,
    [`${props.className}`]: props.className
  });

  const iconClassName = classSet({
    "c-heading-icon": true,
    circular: level <= 2,
    [props.icon!]: props.icon,
    icon: true
  });

  return (
    <Tag className={className}>
      {props.icon && <i className={iconClassName} />}
      <div className="c-heading-content">
        {props.title}
        <div className="c-heading-subtitle">{props.subtitle}</div>
      </div>
    </Tag>
  );
};
