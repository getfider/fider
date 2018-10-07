import * as React from "react";
import { PostStatus } from "@fider/models";
import { DropDown, DropDownItem } from "@fider/components";
import { Fider } from "@fider/services";

import "./PostFilter.scss";

interface PostFilterProps {
  activeView: string;
  countPerStatus: { [key: string]: number };
  viewChanged: (name: string) => void;
}

export class PostFilter extends React.Component<PostFilterProps, {}> {
  constructor(props: PostFilterProps) {
    super(props);
  }

  private handleChangeView = (item: DropDownItem) => {
    this.props.viewChanged(item.value as string);
  };

  public renderSelected = (item?: DropDownItem) => {
    return <>{item!.label.toLowerCase()}</>;
  };

  public renderItem = (item: DropDownItem) => {
    return (
      <span>
        {item.label} <a className="counter">{this.props.countPerStatus[item.value]}</a>
      </span>
    );
  };

  public render() {
    const options: DropDownItem[] = [
      { value: "trending", label: "Trending" },
      { value: "recent", label: "Recent" },
      { value: "most-wanted", label: "Most Wanted" },
      { value: "most-discussed", label: "Most Discussed" }
    ];

    if (Fider.session.isAuthenticated) {
      options.push({ value: "my-votes", label: "My Votes" });
    }

    PostStatus.All.filter(s => s.filterable && this.props.countPerStatus[s.value]).forEach(s => {
      options.push({
        label: s.title,
        value: s.value
      });
    });

    const viewExists = options.filter(x => x.value === this.props.activeView).length > 0;
    const activeView = viewExists ? this.props.activeView : "trending";

    return (
      <>
        Show{" "}
        <DropDown
          className="l-post-filter"
          header="What do you want to see?"
          inline={true}
          items={options}
          renderSelected={this.renderSelected}
          renderItem={this.renderItem}
          defaultValue={activeView}
          onChange={this.handleChangeView}
        />{" "}
      </>
    );
  }
}
