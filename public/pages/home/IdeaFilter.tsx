import * as React from 'react';
import { Idea, IdeaStatus } from '@fider/models';

type IdeaFilterFunction = (ideas: Idea[]) => Idea[];

interface IdeaFilterProps {
    filterChanged: (filter: IdeaFilterFunction) => void;
}

const filterers: {[key: string]: (idea: Idea) => boolean} = {
    'recent': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
    'most-wanted': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
    'most-discussed': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
    'completed': (idea: Idea) => idea.status === IdeaStatus.Completed.value,
    'declined': (idea: Idea) => idea.status === IdeaStatus.Declined.value
};

const sorterers: {[key: string]: (left: Idea, right: Idea) => number} = {
    'recent': (left: Idea, right: Idea) => new Date(right.createdOn).getTime() - new Date(left.createdOn).getTime(),
    'most-wanted': (left: Idea, right: Idea) => right.totalSupporters - left.totalSupporters,
    'most-discussed': (left: Idea, right: Idea) => right.totalComments - left.totalComments,
    'completed': (left: Idea, right: Idea) => new Date(left.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime(),
    'declined': (left: Idea, right: Idea) => new Date(left.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime()
};

const getFilter = (value: string): IdeaFilterFunction => {
    return (ideas: Idea[]) => ideas.filter(filterers[value]).sort(sorterers[value]);
};

export class IdeaFilter extends React.Component<IdeaFilterProps, {}> {
    private element: HTMLDivElement;

    constructor(props: IdeaFilterProps) {
        super(props);
    }

    public static defaultFilter = getFilter('recent');

    public componentDidMount() {
        $(this.element).dropdown({
          onChange: (value: string) => {
            this.props.filterChanged(getFilter(value));
          }
        });
    }

    public render() {
        return <h4 className="ui header">
                    <div className="content">
                    Showing {' '}
                    <div className="ui inline dropdown" ref={(e) => this.element = e!}>
                        <div className="text">recent ideas</div>
                        <i className="dropdown icon"></i>
                        <div className="menu">
                        <div className="header">What do you want to see?</div>
                        <div className="item" data-value="recent" data-text="recent ideas">Recent</div>
                        <div className="item" data-value="most-wanted" data-text="most wanted ideas">Most Wanted</div>
                        <div className="item" data-value="most-discussed" data-text="most discussed ideas">Most Discussed</div>
                        <div className="item" data-value="completed" data-text="completed ideas">Completed</div>
                        <div className="item" data-value="declined" data-text="declined ideas">Declined</div>
                        </div>
                    </div>
                    </div>
                </h4>;
    }
}
