import * as React from "react";
import { Idea, IdeaStatus } from "@fider/models";
import { actions } from "@fider/services";
import { Dropdown, DropdownProps, DropdownItemProps, DropdownOnSearchChangeData } from "@fider/components";

interface IdeaSearchProps {
  exclude?: number[];
  onChanged(ideaNumber: number): void;
}

interface IdeaSearchState {
  ideas: Idea[];
}

export class IdeaSearch extends React.Component<IdeaSearchProps, IdeaSearchState> {
  private timer?: number;

  constructor(props: IdeaSearchProps) {
    super(props);
    this.state = {
      ideas: []
    };
    this.search("");
  }

  private onSearchChange = (e: React.SyntheticEvent<HTMLElement>, data: DropdownOnSearchChangeData) => {
    this.search(data.searchQuery);
  };

  private onChange = (e: React.SyntheticEvent<HTMLElement>, data: DropdownProps) => {
    this.props.onChanged(parseInt(data.value as string, 10));
  };

  private search = (searchQuery: string) => {
    window.clearTimeout(this.timer);
    this.timer = window.setTimeout(() => {
      actions.searchIdeas(searchQuery, "", []).then(res => {
        const ideas =
          this.props.exclude && this.props.exclude.length > 0
            ? res.data.filter(i => this.props.exclude!.indexOf(i.number) === -1)
            : res.data;
        this.setState({ ideas });
      });
    }, 200);
  };

  private returnAll = (options: DropdownItemProps[], value: string) => options;

  public render() {
    const options = this.state.ideas.map(i => {
      const status = IdeaStatus.Get(i.status);
      return {
        key: i.number,
        text: i.title,
        value: i.number,
        content: (
          <>
            <span className="support">
              <i className="medium caret up icon" />
              {i.totalSupporters}
            </span>
            <span className={`gm-status-label gm-status-${status.slug}`}>{status.title}</span>
            {i.title}
          </>
        )
      };
    });

    return (
      <Dropdown
        className="c-idea-search"
        fluid={true}
        selectOnBlur={false}
        selection={true}
        search={this.returnAll}
        options={options}
        placeholder="Search original idea"
        onChange={this.onChange}
        onSearchChange={this.onSearchChange}
      />
    );
  }
}
