import * as React from 'react';
import { User, IdeaResponse, IdeaStatus } from '@fider/models';
import { Gravatar, MultiLineText, Moment, UserName } from '@fider/components/common';

interface IdeaResponseProps {
  status: number;
  response: IdeaResponse;
}

const DuplicateDetails = (props: IdeaResponseProps): JSX.Element | null => {
  const duplicate = props.response.duplicate;
  if (!duplicate) {
    return null;
  }
  const status = IdeaStatus.Get(duplicate.status);

  return (
    <div className="content">
      <span>&#8618;</span>
      <span title={status.title} className={`ui mini empty circular ${status.color} label`} />
      <a href={`/ideas/${duplicate.number}/${duplicate.slug}`}>{duplicate.title}</a>
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
      <div className="fdr-response item ui segment">
        <span className={`ui mini label ${status.color}`}>{status.title}</span>
        <Gravatar user={props.response.user}/> <UserName user={props.response.user} />
        <span className="info">
            <Moment date={props.response.respondedOn} />
        </span>
        {status === IdeaStatus.Duplicate ? DuplicateDetails(props) : StatusDetails(props)}
      </div>
    );
  }

  return <div/>;
};
