import "./ShowPostResponse.scss";

import * as React from "react";
import { PostResponse, PostStatus } from "@fider/models";
import { Gravatar, MultiLineText, Moment, UserName, Segment } from "@fider/components/common";

interface ShowPostStatusProps {
  status: PostStatus;
}

export const ShowPostStatus = (props: ShowPostStatusProps) => {
  return <span className={`status-label status-${props.status.slug}`}>{props.status.title}</span>;
};

const DuplicateDetails = (props: PostResponseProps): JSX.Element | null => {
  if (!props.response) {
    return null;
  }

  const original = props.response.original;
  if (!original) {
    return null;
  }
  const status = PostStatus.Get(original.status);

  return (
    <div className="content">
      <span>&#8618;</span> <a href={`/posts/${original.number}/${original.slug}`}>{original.title}</a>
    </div>
  );
};

interface PostResponseProps {
  status: number;
  response: PostResponse | null;
}

const StatusDetails = (props: PostResponseProps): JSX.Element | null => {
  if (!props.response || !props.response.text) {
    return null;
  }

  return (
    <div className="content">
      <MultiLineText text={props.response.text} style="full" />
    </div>
  );
};

export const ShowPostResponse = (props: PostResponseProps): JSX.Element => {
  const status = PostStatus.Get(props.status);

  if (props.response && status.show) {
    return (
      <Segment className="l-response">
        <ShowPostStatus status={status} />
        <Gravatar user={props.response.user} size="small" /> <UserName user={props.response.user} />
        {status === PostStatus.Duplicate ? DuplicateDetails(props) : StatusDetails(props)}
      </Segment>
    );
  }

  return <div />;
};
