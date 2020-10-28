import "./TagsFilter.scss";

import React from "react";
import { Tag } from "@fider/models";
import { ShowTag } from "@fider/components/ShowTag";
import { DropDown, DropDownItem } from "@fider/components";
import { FaCheck } from "react-icons/fa";
import { WithTranslation, withTranslation } from "react-i18next";

interface TagsFilterProps extends WithTranslation {
  tags: Tag[];
  defaultSelection: string[];
  selectionChanged: (selected: string[]) => void;
}

interface TagsFilterState {
  selected: string[];
}

class InternalTagsFilter extends React.Component<TagsFilterProps, TagsFilterState> {
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
    const { t } = this.props;
    const text =
      this.state.selected.length === 0
        ? t("home.tagsFilter.anyTag")
        : this.state.selected.length === 1
        ? t("home.tagsFilter.oneTag")
        : t("home.tagsFilter.nTags", { n: this.state.selected.length });
    return <>{text}</>;
  };

  public render() {
    if (this.props.tags.length === 0) {
      return null;
    }

    const items = this.props.tags.map(tag => {
      return {
        value: tag.slug,
        label: tag.name,
        render: (
          <div className={this.state.selected.indexOf(tag.slug) >= 0 ? "selected-tag" : ""}>
            <FaCheck />
            <ShowTag tag={tag} size="mini" circular={true} />
            {t.name}
          </div>
        )
      };
    });
    const { t } = this.props;
    return (
      <div>
        <span className="subtitle">{t("home.tagsFilter.subtitle")}</span>
        <DropDown
          className="l-tags-filter"
          inline={true}
          style="simple"
          highlightSelected={false}
          items={items}
          onChange={this.onChange}
          renderText={this.renderText}
        />
      </div>
    );
  }
}

export const TagsFilter = withTranslation()(InternalTagsFilter);
