import "./Heading.scss";

import * as React from "react";
import { classSet } from "@fider/services";

interface HeadingLogo {
  title: string;
  dividing?: boolean;
  level?: 1 | 2 | 3 | 4 | 5;
  icon?: string;
  subtitle?: string;
  iconClassName?: string;
}

export const Heading = (props: HeadingLogo) => {
  const Tag = `h${props.level || 2}`;
  const iconClassName = props.iconClassName || "circular";
  const className = classSet({
    "c-heading": true,
    "m-dividing": props.dividing || false
  });

  return (
    <Tag className={className}>
      <i className={`c-heading-icon ${iconClassName} ${props.icon} icon`} />
      <div className="c-heading-content">
        {props.title}
        <div className="c-heading-subtitle">{props.subtitle}</div>
      </div>
    </Tag>
  );
};
