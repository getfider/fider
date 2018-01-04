import * as React from 'react';
import { Idea, IdeaStatus } from '@fider/models';

export type IdeaFilterFunction = (ideas: Idea[]) => Idea[];

interface IdeaFilterProps {
    ideas: Idea[];
    activeFilter: string;
    filterChanged: (name: string) => void;
}

const filterers: {[key: string]: (idea: Idea) => boolean} = {
    'trending': (idea: Idea) =>  IdeaStatus.Get(idea.status).closed === false,
    'recent': (idea: Idea) => IdeaStatus.Get(idea.status).closed === false,
    'most-wanted': (idea: Idea) => IdeaStatus.Get(idea.status).closed === false,
    'most-discussed': (idea: Idea) => IdeaStatus.Get(idea.status).closed === false,
    'planned': (idea: Idea) => idea.status === IdeaStatus.Planned.value,
    'started': (idea: Idea) => idea.status === IdeaStatus.Started.value,
    'completed': (idea: Idea) => idea.status === IdeaStatus.Completed.value,
    'declined': (idea: Idea) => idea.status === IdeaStatus.Declined.value
};

const names: {[key: string]: string} = {
    'trending': 'trending',
    'recent': 'recent',
    'most-wanted': 'most wanted',
    'most-discussed': 'most discussed',
    'planned': 'planned',
    'started': 'started',
    'completed': 'completed',
    'declined': 'declined'
};

const sorterers: {[key: string]: (left: Idea, right: Idea) => number} = {
    'trending': (left: Idea, right: Idea) => right.ranking !== left.ranking ? right.ranking - left.ranking : new Date(right.createdOn).getTime() - new Date(left.createdOn).getTime(),
    'recent': (left: Idea, right: Idea) => new Date(right.createdOn).getTime() - new Date(left.createdOn).getTime(),
    'most-wanted': (left: Idea, right: Idea) => right.totalSupporters - left.totalSupporters,
    'most-discussed': (left: Idea, right: Idea) => right.totalComments - left.totalComments,
    'planned': (left: Idea, right: Idea) => new Date(right.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime(),
    'started': (left: Idea, right: Idea) => new Date(right.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime(),
    'completed': (left: Idea, right: Idea) => new Date(right.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime(),
    'declined': (left: Idea, right: Idea) => new Date(right.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime()
};

export class IdeaFilter extends React.Component<IdeaFilterProps, {}> {
    private element: HTMLDivElement;

    constructor(props: IdeaFilterProps) {
        super(props);
    }

    public static getFilter(value: string): IdeaFilterFunction {
        if (!filterers[value]) {
            value = 'trending';
        }

        return (ideas: Idea[]) => ideas.filter(filterers[value]).sort(sorterers[value]);
    }

    public componentDidMount() {
        $(this.element).dropdown({
            onChange: (value: string) => {
                this.props.filterChanged(value);
            }
        });
    }

    public render() {
        let activeFilter = this.props.activeFilter;
        if (!(this.props.activeFilter in names)) {
            activeFilter = 'trending';
        }

        const grouped = this.props.ideas.reduce<{ [status: number]: number }>((group, idea) => {
            group[idea.status] = (group[idea.status] || 0) + 1;
            return group;
        }, {});

        const statusFilterItems = IdeaStatus.All.filter((s) => s.filterable && grouped[s.value]).map((s) => (
            <div key={s.value} className={`item ${activeFilter === s.slug && 'active'}`} data-value={s.slug} data-text={s.title.toLowerCase()}>
                {s.title}
                <a className="ui mini circular label">{grouped[s.value]}</a>
            </div>
        ));

        return (
          <h4 className="ui header">
              <div className="content">
                Showing {' '}
                <div className="ideas-filter ui inline dropdown" ref={(e) => this.element = e!}>
                    <div className="text">{names[activeFilter]}</div>
                    <i className="dropdown icon" />
                    <div className="menu">
                        <div className="header">What do you want to see?</div>
                        <div className={`item ${activeFilter === 'trending' && 'active'}`} data-value="trending" data-text="trending">Trending</div>
                        <div className={`item ${activeFilter === 'recent' && 'active'}`} data-value="recent" data-text="recent">Recent</div>
                        <div className={`item ${activeFilter === 'most-wanted' && 'active'}`} data-value="most-wanted" data-text="most wanted">Most Wanted</div>
                        <div className={`item ${activeFilter === 'most-discussed' && 'active'}`} data-value="most-discussed" data-text="most discussed">Most Discussed</div>
                        {statusFilterItems}
                    </div>
                </div>
              </div>
          </h4>
        );
    }
}
