import * as React from "react";
import { Post, PostStatus, CurrentUser } from "@fider/models";
import { Dropdown, DropdownItemProps, DropdownProps } from "@fider/components";
import { Fider } from "@fider/services";

interface PostFilterProps {
  activeFilter: string;
  countPerStatus: { [key: string]: number };
  filterChanged: (name: string) => void;
}

export class PostFilter extends React.Component<PostFilterProps, {}> {
  constructor(props: PostFilterProps) {
    super(props);
  }

  private handleChangeFilter = (item: React.SyntheticEvent<HTMLElement>, data: DropdownProps) => {
    this.props.filterChanged(data.value as string);
  };

  public render() {
    const options: DropdownItemProps[] = [
      { text: "trending", value: "trending", content: "Trending" },
      { text: "recent", value: "recent", content: "Recent" },
      { text: "most wanted", value: "most-wanted", content: "Most Wanted" },
      { text: "most discussed", value: "most-discussed", content: "Most Discussed" }
    ];

    if (Fider.session.isAuthenticated) {
      options.push({ text: "my votes", value: "my-votes", content: "My Votes" });
    }

    PostStatus.All.filter(s => s.filterable && this.props.countPerStatus[s.value]).forEach(s => {
      options.push({
        text: s.title.toLowerCase(),
        value: s.slug,
        content: (
          <span>
            {s.title} <a className="counter">{this.props.countPerStatus[s.value]}</a>
          </span>
        )
      });
    });

    const filterExists = options.filter(x => x.value === this.props.activeFilter).length > 0;
    const activeFilter = filterExists ? this.props.activeFilter : "trending";

    return (
      <>
        Show{" "}
        <Dropdown
          className="l-post-filter"
          header="What do you want to see?"
          inline={true}
          options={options}
          defaultValue={activeFilter}
          onChange={this.handleChangeFilter}
        />
      </>
    );
  }
}
