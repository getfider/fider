import * as React from "react";
import { Tag } from "@fider/models";
import { ShowTag } from "@fider/components/ShowTag";
import { Dropdown, DropdownProps } from "@fider/components";

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

  private onChange = (e: React.SyntheticEvent<HTMLElement>, data: DropdownProps) => {
    let selected = [];
    const idx = this.state.selected.indexOf(data.value as string);
    if (idx >= 0) {
      selected = this.state.selected.splice(idx, 1) && this.state.selected;
    } else {
      selected = this.state.selected.concat(data.value as string);
    }
    this.setState({ selected });
    this.props.selectionChanged(selected);
  };

  public render() {
    if (this.props.tags.length === 0) {
      return null;
    }

    const options = this.props.tags.map(t => {
      return {
        value: t.slug,
        text: t.name,
        content: (
          <div className={this.state.selected.indexOf(t.slug) >= 0 ? "selected-tag" : ""}>
            <i className="icon check" />
            <ShowTag tag={t} circular={true} />
            {t.name}
          </div>
        )
      };
    });

    const text =
      this.state.selected.length === 0
        ? "any tag"
        : this.state.selected.length === 1
          ? "1 tag"
          : `${this.state.selected.length} tags`;

    return (
      <>
        with{" "}
        <Dropdown
          className="tags-filter"
          selectOnBlur={false}
          text={text}
          defaultValue="0"
          inline={true}
          options={options}
          onChange={this.onChange}
        />
      </>
    );
  }
}
