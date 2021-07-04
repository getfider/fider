import "./AdminBasePage.scss"

import React from "react"
import { PageTitle } from "@fider/components"
import { SideMenu, SideMenuToggler } from "./SideMenu"
import { HStack } from "@fider/components/layout"

interface AdminPageContainerProps {
  id: string
  name: string
  title: string
  subtitle: string
  children: React.ReactNode
}

export const AdminPageContainer = (props: AdminPageContainerProps) => {
  return (
    <div id={props.id} className="page container">
      <HStack justify="between">
        <PageTitle title={props.title} subtitle={props.subtitle} />
        <SideMenuToggler />
      </HStack>

      <div className="c-admin-basepage">
        <SideMenu activeItem={props.name} />
        <div>{props.children}</div>
      </div>
    </div>
  )
}

export abstract class AdminBasePage<P, S> extends React.Component<P, S> {
  public abstract id: string
  public abstract name: string
  public abstract title: string
  public abstract subtitle: string
  public abstract content(): JSX.Element

  public render() {
    return (
      <AdminPageContainer id={this.id} name={this.name} title={this.title} subtitle={this.subtitle}>
        {this.content()}
      </AdminPageContainer>
    )
  }
}
