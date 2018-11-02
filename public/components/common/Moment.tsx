import React from "react";
import { formatDate, timeSince } from "@fider/services";

interface MomentText {
  date: Date | string;
}

export const Moment = (props: MomentText) => {
  if (!props.date) {
    return <span />;
  }

  const now = new Date();
  const date = props.date instanceof Date ? props.date : new Date(props.date);

  const diff = (now.getTime() - date.getTime()) / (60 * 60 * 24 * 1000);
  const display = diff >= 365 ? formatDate(props.date) : timeSince(now, date);

  return (
    <span className="date" title={formatDate(props.date)}>
      {display}
    </span>
  );
};
