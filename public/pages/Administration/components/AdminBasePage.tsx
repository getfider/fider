import "./AdminBasePage.scss"

import React from "react"
import { PageTitle } from "@fider/components"
import { SideMenu, SideMenuToggler } from "./SideMenu"
import { HStack } from "@fider/components/layout"

export abstract class AdminBasePage<P, S> extends React.Component<P, S> {
  public abstract id: string
  public abstract name: string
  public abstract title: string
  public abstract subtitle: string
  public abstract content(): JSX.Element

  public render() {
    return (
      <div id={this.id} className="page container">
        <HStack justify="between">
          <PageTitle title={this.title} subtitle={this.subtitle} />
          <SideMenuToggler />
        </HStack>

        <div className="c-admin-basepage">
          <SideMenu activeItem={this.name} />
          <div>{this.content()}</div>
        </div>
      </div>
    )
  }
}
