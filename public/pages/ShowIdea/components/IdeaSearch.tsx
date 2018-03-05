import * as React from 'react';
import { Idea, IdeaStatus } from '@fider/models';
import { actions } from '@fider/services';

interface IdeaSearchProps {
  exclude?: number[];
  onChanged(ideaNumber: number): void;
}

interface IdeaSearchState {
  ideas?: Idea[];
}

export class IdeaSearch extends React.Component<IdeaSearchProps, IdeaSearchState> {
  private element?: HTMLDivElement;

  constructor(props: IdeaSearchProps) {
    super(props);
    this.state = { };
    actions.getAllIdeas().then((res) => {
      const ideas = this.props.exclude && this.props.exclude.length > 0
        ? res.data.filter((i) => this.props.exclude!.indexOf(i.number) === -1)
        : res.data;
      this.setState({ ideas });
    });
  }

  public componentDidUpdate() {
    $(this.element).dropdown({
      fullTextSearch: true,
      onChange: (ideaNumber: string) => {
        if (ideaNumber) {
          this.props.onChanged(parseInt(ideaNumber, 10));
        }
      }
    });
  }

  public render() {
    const items = this.state.ideas && (
      <div className="menu">
        {this.state.ideas.map((i) => {
          const status = IdeaStatus.Get(i.status);
          return (
            <div key={i.id} className="item" data-value={i.number} data-text={i.title}>
              <span className="support"><i className="medium caret up icon" />{i.totalSupporters}</span>
              <span className={`ui mini label ${status.color}`}>
                {status.title}
              </span>
              {i.title}
            </div>
          );
        })}
      </div>
    );

    const className = `ui selection ${items ? 'search' : 'loading'} fluid dropdown fdr-idea-search`;

    return (
      <div className={className} ref={(e) => this.element = e!}>
        <i className="dropdown icon" />
        <div className="default text">Search original idea</div>
        {items}
      </div>
    );
  }
}
