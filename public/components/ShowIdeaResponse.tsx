import * as React from "react";
import { User, IdeaResponse, IdeaStatus } from "@fider/models";
import { Gravatar, MultiLineText, Moment, UserName } from "@fider/components/common";

interface IdeaResponseProps {
  status: number;
  response: IdeaResponse;
}

const DuplicateDetails = (props: IdeaResponseProps): JSX.Element | null => {
  const original = props.response.original;
  if (!original) {
    return null;
  }
  const status = IdeaStatus.Get(original.status);

  return (
    <div className="content">
      <span>&#8618;</span>
      <span title={status.title} className={`status-${status.slug}`} />
      <a href={`/ideas/${original.number}/${original.slug}`}>{original.title}</a>
    </div>
  );
};

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
      <div className="c-response item ui segment">
        <span className={`status-label status-${status.slug}`}>{status.title}</span>
        <Gravatar user={props.response.user} /> <UserName user={props.response.user} />
        <span className="info">
          <Moment date={props.response.respondedOn} />
        </span>
        {status === IdeaStatus.Duplicate ? DuplicateDetails(props) : StatusDetails(props)}
      </div>
    );
  }

  return <div />;
};
