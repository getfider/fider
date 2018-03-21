import * as React from "react";
import { Idea, IdeaStatus } from "@fider/models";
import { Dropdown, DropdownItemProps } from "@fider/components";

interface IdeaFilterProps {
  activeFilter: string;
  countPerStatus: { [key: string]: number };
  filterChanged: (name: string) => void;
}

export class IdeaFilter extends React.Component<IdeaFilterProps, {}> {
  constructor(props: IdeaFilterProps) {
    super(props);
  }

  public render() {
    const options: DropdownItemProps[] = [
      { text: "trending", value: "trending", content: "Trending" },
      { text: "recent", value: "recent", content: "Recent" },
      { text: "most wanted", value: "most-wanted", content: "Most Wanted" },
      { text: "most discussed", value: "most discussed", content: "Most Discussed" }
    ];

    IdeaStatus.All.filter(s => s.filterable && this.props.countPerStatus[s.value]).forEach(s => {
      options.push({
        text: s.title.toLowerCase(),
        value: s.slug,
        content: (
          <span>
            {s.title} <a className="ui mini circular label">{this.props.countPerStatus[s.value]}</a>
          </span>
        )
      });
    });

    const activeFilter = options.find(x => x.value === this.props.activeFilter) ? this.props.activeFilter : "trending";

    return (
      <>
        Show{" "}
        <Dropdown
          className="ideas-filter"
          header="What do you want to see?"
          inline={true}
          options={options}
          defaultValue={activeFilter}
          onChange={(item, data) => this.props.filterChanged(data.value as string)}
        />
      </>
    );
  }
}
