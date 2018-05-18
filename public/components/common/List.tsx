import "./List.scss";

import * as React from "react";
import { classSet } from "@fider/services";

interface ListProps {
  className?: string;
}

interface ListItemProps {
  className?: string;
  onClick?: () => void;
}

export const List: React.StatelessComponent<ListProps> = props => {
  return <div className={`c-list ${props.className || ""}`}>{props.children}</div>;
};

export const ListItem: React.StatelessComponent<ListItemProps> = props => {
  const className = classSet({
    "c-list-item": true,
    [props.className || ""]: true,
    "m-selectable": props.onClick
  });

  if (props.onClick) {
    return (
      <div className={className} onClick={props.onClick}>
        {props.children}
      </div>
    );
  }
  return <div className={className}>{props.children}</div>;
};
