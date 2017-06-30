import * as React from 'react';
import * as moment from 'moment';
import { User, IdeaResponse, IdeaStatus } from '@fider/models';
import { Gravatar, MultiLineText } from '@fider/components/common';

interface IdeaResponseProps {
  status: number;
  response: IdeaResponse;
}

export const ShowIdeaResponse = (props: IdeaResponseProps): JSX.Element => {
    const status = IdeaStatus.Get(props.status);

    if (props.response && status.show) {
        return <div className="fdr-response item ui raised segment">
                <span className={`ui ribbon label ${status.color}`}>{ status.title }</span>
                <div className="info">
                    <Gravatar email={props.response.user.email}/> <u>{props.response.user.name}</u>
                    <span title={props.response.respondedOn.toString()}>
                    { moment(props.response.respondedOn).fromNow() }
                    </span>
                </div>
                <div className="content">
                    <MultiLineText text={ props.response.text } />
                </div>
                </div>;
    }
    return <div/>;
};
