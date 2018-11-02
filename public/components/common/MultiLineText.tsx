import React from "react";
import { markdown } from "@fider/services";

interface MultiLineText {
  className?: string;
  text?: string;
  style: "full" | "simple";
}

export const MultiLineText = (props: MultiLineText) => {
  if (!props.text) {
    return <p />;
  }

  const func = props.style === "full" ? markdown.full : markdown.simple;
  return (
    <div className={`markdown-body ${props.className || ""}`} dangerouslySetInnerHTML={{ __html: func(props.text) }} />
  );
};
