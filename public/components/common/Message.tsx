import "./Message.scss";

import React from "react";
import { classSet } from "@fider/services";
import { FaBan, FaCheckCircle } from "react-icons/fa";

interface MessageProps {
  type: "success" | "error";
  showIcon?: boolean;
}

export const Message: React.StatelessComponent<MessageProps> = props => {
  const className = classSet({
    "c-message": true,
    [`m-${props.type}`]: true
  });

  const icon = props.type === "error" ? <FaBan /> : <FaCheckCircle />;

  return (
    <div className={className}>
      {props.showIcon === true && icon}
      {props.children}
    </div>
  );
};
