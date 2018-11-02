import React from "react";
import { Tag } from "@fider/models";
import { ShowTag } from "@fider/components/ShowTag";
import { DropDown, DropDownItem } from "@fider/components";

import "./TagsFilter.scss";

interface TagsFilterProps {
  tags: Tag[];
  defaultSelection: string[];
  selectionChanged: (selected: string[]) => void;
}

interface TagsFilterState {
  selected: string[];
}

export class TagsFilter extends React.Component<TagsFilterProps, TagsFilterState> {
  constructor(props: TagsFilterProps) {
    super(props);
    this.state = {
      selected: props.defaultSelection
    };
  }

  private onChange = (item: DropDownItem) => {
    let selected = [];
    const idx = this.state.selected.indexOf(item.value as string);
    if (idx >= 0) {
      selected = this.state.selected.splice(idx, 1) && this.state.selected;
    } else {
      selected = this.state.selected.concat(item.value as string);
    }
    this.setState({ selected });
    this.props.selectionChanged(selected);
  };

  private renderText = () => {
    const text =
      this.state.selected.length === 0
        ? "any tag"
        : this.state.selected.length === 1
          ? "1 tag"
          : `${this.state.selected.length} tags`;
    return <>{text}</>;
  };

  public render() {
    if (this.props.tags.length === 0) {
      return null;
    }

    const items = this.props.tags.map(t => {
      return {
        value: t.slug,
        label: t.name,
        render: (
          <div className={this.state.selected.indexOf(t.slug) >= 0 ? "selected-tag" : ""}>
            <i className="icon check" />
            <ShowTag tag={t} size="mini" circular={true} />
            {t.name}
          </div>
        )
      };
    });

    return (
      <>
        with{" "}
        <DropDown
          className="l-tags-filter"
          inline={true}
          highlightSelected={false}
          items={items}
          renderText={this.renderText}
          onChange={this.onChange}
        />
      </>
    );
  }
}
