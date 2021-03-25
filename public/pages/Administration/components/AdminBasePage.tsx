import "./AdminBasePage.scss"

import React from "react"
import { Heading } from "@fider/components"
import { SideMenu, SideMenuToggler } from "./SideMenu"
import { IconType } from "react-icons"

export abstract class AdminBasePage<P, S> extends React.Component<P, S> {
  public abstract id: string
  public abstract name: string
  public abstract icon: IconType
  public abstract title: string
  public abstract subtitle: string
  public abstract content(): JSX.Element

  private toggleSideMenu = (active: boolean) => {
    const el = document.querySelector(".hidden-lg .c-side-menu") as HTMLElement
    if (el) {
      el.style.display = active ? "" : "none"
    }
  }

  public render() {
    return (
      <div id={this.id} className="page container">
        <Heading title={this.title} icon={this.icon} subtitle={this.subtitle} className="l-admin-heading" />
        <SideMenuToggler onToggle={this.toggleSideMenu} />

        <div className="row">
          <div className="col-lg-2 hidden-sm hidden-md">
            <SideMenu visible={true} activeItem={this.name} />
          </div>
          <div className="col-lg-10 col-md-12">
            <SideMenu className="hidden-lg hidden-xl" visible={false} activeItem={this.name} />
            {this.content()}
          </div>
        </div>
      </div>
    )
  }
}
