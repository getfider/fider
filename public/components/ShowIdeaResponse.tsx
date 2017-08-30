import * as React from 'react';
import { User, IdeaResponse, IdeaStatus } from '@fider/models';
import { Gravatar, MultiLineText, Moment, UserName } from '@fider/components/common';

interface IdeaResponseProps {
  status: number;
  response: IdeaResponse;
}

export const ShowIdeaResponse = (props: IdeaResponseProps): JSX.Element => {
    const status = IdeaStatus.Get(props.status);

    if (props.response && status.show) {
        return <div className="fdr-response item ui segment">
                    <span className={`ui mini label ${status.color}`}>{ status.title }</span>
                    <Gravatar name={props.response.user.name} hash={props.response.user.gravatar}/> <UserName user={props.response.user} />
                    <span className="info">
                        <Moment date={props.response.respondedOn} />
                    </span>
                    <div className="content">
                        <MultiLineText text={ props.response.text } style="full" />
                    </div>
                </div>;
    }
    return <div/>;
};
