import * as React from "react";

interface SiteMenuProps {
  active?: string;
}

export const SideMenu = (props: SiteMenuProps) => {
  const active = props.active || "general";
  return (
    <div className="ui vertical pointing menu fluid">
      <a className={`item ${active === "general" ? "active" : ""}`}>General</a>
      <a className={`item ${active === "members" ? "active" : ""}`}>Members</a>
      <a className={`item ${active === "tags" ? "active" : ""}`}>Tags</a>
    </div>
  );
};
