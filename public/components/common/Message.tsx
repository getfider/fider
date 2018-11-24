import "./Message.scss";

import React from "react";
import { classSet } from "@fider/services";
import { FaBan, FaRegCheckCircle } from "react-icons/fa";

interface MessageProps {
  type: "success" | "error";
  showIcon?: boolean;
}

export const Message: React.StatelessComponent<MessageProps> = props => {
  const className = classSet({
    "c-message": true,
    [`m-${props.type}`]: true
  });

  const icon = props.type === "error" ? <FaBan /> : <FaRegCheckCircle />;

  return (
    <p className={className}>
      {props.showIcon === true && icon}
      <span>{props.children}</span>
    </p>
  );
};
