import "./ShowIdeaResponse.scss";

import * as React from "react";
import { User, IdeaResponse, IdeaStatus } from "@fider/models";
import { Gravatar, MultiLineText, Moment, UserName, Segment } from "@fider/components/common";

interface ShowIdeaStatusProps {
  status: IdeaStatus;
}

export const ShowIdeaStatus = (props: ShowIdeaStatusProps) => {
  return <span className={`status-label status-${props.status.slug}`}>{props.status.title}</span>;
};

const DuplicateDetails = (props: IdeaResponseProps): JSX.Element | null => {
  const original = props.response.original;
  if (!original) {
    return null;
  }
  const status = IdeaStatus.Get(original.status);

  return (
    <div className="content">
      <span>&#8618;</span> <a href={`/ideas/${original.number}/${original.slug}`}>{original.title}</a>
    </div>
  );
};

interface IdeaResponseProps {
  status: number;
  response: IdeaResponse;
}

const StatusDetails = (props: IdeaResponseProps): JSX.Element => {
  return (
    <div className="content">
      <MultiLineText text={props.response.text} style="full" />
    </div>
  );
};

export const ShowIdeaResponse = (props: IdeaResponseProps): JSX.Element => {
  const status = IdeaStatus.Get(props.status);

  if (props.response && status.show) {
    return (
      <Segment className="l-response">
        <ShowIdeaStatus status={status} />
        <Gravatar user={props.response.user} size="small" /> <UserName user={props.response.user} />
        {status === IdeaStatus.Duplicate ? DuplicateDetails(props) : StatusDetails(props)}
      </Segment>
    );
  }

  return <div />;
};
