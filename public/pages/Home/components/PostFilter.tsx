import "./PostFilter.scss";

import React from "react";
import { PostStatus } from "@fider/models";
import { DropDown, DropDownItem } from "@fider/components";
import { useTranslation } from "react-i18next";
import { useFider } from "@fider/hooks";

interface PostFilterProps {
  activeView: string;
  countPerStatus: { [key: string]: number };
  viewChanged: (name: string) => void;
}

export const PostFilter = (props: PostFilterProps) => {
  const fider = useFider();
  const { t } = useTranslation();

  const handleChangeView = (item: DropDownItem) => {
    props.viewChanged(item.value as string);
  };
  const renderText = (item: DropDownItem | undefined) => {
    if (!item) {
      return <div />;
    }
    return <div>{t(item.label)}</div>;
  };

  const options: DropDownItem[] = [
    { value: "trending", label: t("home.postFilter.options.trending") },
    { value: "recent", label: t("home.postFilter.options.recent") },
    { value: "most-wanted", label: t("home.postFilter.options.mostWanted") },
    { value: "most-discussed", label: t("home.postFilter.options.mostDiscussed") }
  ];

  if (fider.session.isAuthenticated) {
    options.push({ value: "my-votes", label: t("home.postFilter.options.myVotes") });
  }

  PostStatus.All.filter(s => s.filterable && props.countPerStatus[s.value]).forEach(s => {
    options.push({
      label: s.title,
      value: s.value,
      render: (
        <span>
          {s.title} <a className="counter">{props.countPerStatus[s.value]}</a>
        </span>
      )
    });
  });

  // const viewExists = options.filter(x => x.value === props.activeView).length > 0;
  const activeView = "trending";

  return (
    <div>
      <span className="subtitle">{t("home.postFilter.subtitle")}</span>
      <DropDown
        header={t("home.postFilter.dropDownHeader")}
        className="l-post-filter"
        inline={true}
        style="simple"
        items={options}
        defaultValue={activeView}
        renderText={renderText}
        onChange={handleChangeView}
      />
    </div>
  );
};
