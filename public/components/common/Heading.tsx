import "./Heading.scss";

import * as React from "react";
import { classSet } from "@fider/services";

interface HeadingLogo {
  title: string;
  dividing?: boolean;
  level?: 1 | 2 | 3 | 4 | 5;
  icon?: string;
  subtitle?: string;
}

export const Heading = (props: HeadingLogo) => {
  const level = props.level || 2;
  const Tag = `h${level}`;
  const className = classSet({
    "c-heading": true,
    "m-dividing": props.dividing || false
  });

  const iconClassName = classSet({
    "c-heading-icon": true,
    circular: level <= 2,
    [props.icon!]: props.icon,
    icon: true
  });

  return (
    <Tag className={className}>
      <i className={iconClassName} />
      <div className="c-heading-content">
        {props.title}
        <div className="c-heading-subtitle">{props.subtitle}</div>
      </div>
    </Tag>
  );
};
