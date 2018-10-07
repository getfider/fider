import * as React from "react";
import { Tag } from "@fider/models";
import { ShowTag } from "@fider/components/ShowTag";
import { FiderDropDown, FiderDropDownItem } from "@fider/components";

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

  private onChange = (item: FiderDropDownItem) => {
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

  private renderItem = (item: FiderDropDownItem) => {
    const tag = this.props.tags.filter(t => t.slug === item.value)[0];
    return (
      <div className={this.state.selected.indexOf(tag.slug) >= 0 ? "selected-tag" : ""}>
        <i className="icon check" />
        <ShowTag tag={tag} size="mini" circular={true} />
        {tag.name}
      </div>
    );
  };

  private renderSelected = () => {
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
        label: t.name
      };
    });

    return (
      <>
        with{" "}
        <FiderDropDown
          className="l-tags-filter"
          inline={true}
          items={items}
          renderItem={this.renderItem}
          renderSelected={this.renderSelected}
          onChange={this.onChange}
        />
      </>
    );
  }
}
