import * as React from "react";
import { Idea, IdeaStatus } from "@fider/models";

export type IdeaFilterFunction = (ideas: Idea[]) => Idea[];

interface IdeaFilterProps {
  activeFilter: string;
  countPerStatus: { [key: string]: number };
  filterChanged: (name: string) => void;
}

const names: { [key: string]: string } = {
  trending: "trending",
  recent: "recent",
  "most-wanted": "most wanted",
  "most-discussed": "most discussed",
  planned: "planned",
  started: "started",
  completed: "completed",
  declined: "declined"
};

export class IdeaFilter extends React.Component<IdeaFilterProps, {}> {
  private element?: HTMLDivElement;

  constructor(props: IdeaFilterProps) {
    super(props);
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
      activeFilter = "trending";
    }

    const statusFilterItems = IdeaStatus.All.filter(s => s.filterable && this.props.countPerStatus[s.value]).map(s => (
      <div
        key={s.value}
        className={`item ${activeFilter === s.slug && "active"}`}
        data-value={s.slug}
        data-text={s.title.toLowerCase()}
      >
        {s.title}
        <a className="ui mini circular label">{this.props.countPerStatus[s.value]}</a>
      </div>
    ));

    return (
      <>
        <div className="content">
          Showing{" "}
          <div className="ideas-filter ui inline dropdown" ref={e => (this.element = e!)}>
            <div className="text">{names[activeFilter]}</div>
            <i className="dropdown icon" />
            <div className="menu">
              <div className="header">What do you want to see?</div>
              <div
                className={`item ${activeFilter === "trending" && "active"}`}
                data-value="trending"
                data-text="trending"
              >
                Trending
              </div>
              <div className={`item ${activeFilter === "recent" && "active"}`} data-value="recent" data-text="recent">
                Recent
              </div>
              <div
                className={`item ${activeFilter === "most-wanted" && "active"}`}
                data-value="most-wanted"
                data-text="most wanted"
              >
                Most Wanted
              </div>
              <div
                className={`item ${activeFilter === "most-discussed" && "active"}`}
                data-value="most-discussed"
                data-text="most discussed"
              >
                Most Discussed
              </div>
              {statusFilterItems}
            </div>
          </div>
        </div>
      </>
    );
  }
}
