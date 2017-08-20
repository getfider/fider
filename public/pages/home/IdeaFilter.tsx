import * as React from 'react';
import { Idea, IdeaStatus } from '@fider/models';

export type IdeaFilterFunction = (ideas: Idea[]) => Idea[];

interface IdeaFilterProps {
    activeFilter: string;
    filterChanged: (name: string, filter: IdeaFilterFunction) => void;
}

const filterers: {[key: string]: (idea: Idea) => boolean} = {
    'recent': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
    'most-wanted': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
    'most-discussed': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
    'completed': (idea: Idea) => idea.status === IdeaStatus.Completed.value,
    'declined': (idea: Idea) => idea.status === IdeaStatus.Declined.value
};

const names: {[key: string]: string} = {
    'recent': 'recent',
    'most-wanted': 'most wanted ideas',
    'most-discussed': 'most discussed ideas',
    'completed': 'completed ideas',
    'declined': 'declined ideas'
};

const sorterers: {[key: string]: (left: Idea, right: Idea) => number} = {
    'recent': (left: Idea, right: Idea) => new Date(right.createdOn).getTime() - new Date(left.createdOn).getTime(),
    'most-wanted': (left: Idea, right: Idea) => right.totalSupporters - left.totalSupporters,
    'most-discussed': (left: Idea, right: Idea) => right.totalComments - left.totalComments,
    'completed': (left: Idea, right: Idea) => new Date(left.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime(),
    'declined': (left: Idea, right: Idea) => new Date(left.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime()
};

export class IdeaFilter extends React.Component<IdeaFilterProps, {}> {
    private element: HTMLDivElement;

    constructor(props: IdeaFilterProps) {
        super(props);
    }

    public static getFilter(value: string): IdeaFilterFunction {
        if (!filterers[value]) {
            value = 'recent';
        }

        return (ideas: Idea[]) => ideas.filter(filterers[value]).sort(sorterers[value]);
    }

    public componentDidMount() {
        $(this.element).dropdown({
            onChange: (value: string) => {
                this.props.filterChanged(value, IdeaFilter.getFilter(value));
            }
        });
    }

    public render() {
        let activeFilter = this.props.activeFilter;
        if (!(this.props.activeFilter in names)) {
            activeFilter = 'recent';
        }

        return <h4 className="ui header">
                    <div className="content">
                    Showing {' '}
                    <div className="ui inline dropdown" ref={(e) => this.element = e!}>
                        <div className="text">{ names[activeFilter] }</div>
                        <i className="dropdown icon"></i>
                        <div className="menu">
                        <div className="header">What do you want to see?</div>
                        <div className={`item ${activeFilter === 'recent' && 'active'}`} data-value="recent" data-text="recent ideas">Recent</div>
                        <div className={`item ${activeFilter === 'most-wanted' && 'active'}`} data-value="most-wanted" data-text="most wanted ideas">Most Wanted</div>
                        <div className={`item ${activeFilter === 'most-discussed' && 'active'}`} data-value="most-discussed" data-text="most discussed ideas">Most Discussed</div>
                        <div className={`item ${activeFilter === 'completed' && 'active'}`} data-value="completed" data-text="completed ideas">Completed</div>
                        <div className={`item ${activeFilter === 'declined' && 'active'}`} data-value="declined" data-text="declined ideas">Declined</div>
                        </div>
                    </div>
                    </div>
                </h4>;
    }
}
