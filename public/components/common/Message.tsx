import "./Message.scss";

import * as React from "react";
import { classSet } from "@fider/services";

interface MessageProps {
  type: "success" | "error";
}

export const Message: React.StatelessComponent<MessageProps> = props => {
  const className = classSet({
    "c-message": true,
    [`m-${props.type}`]: true
  });

  const icon = props.type === "error" ? "ban" : "check circle outline";

  return (
    <div className={className}>
      <i className={`${icon} icon`} />
      {props.children}
    </div>
  );
};
