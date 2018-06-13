import * as React from "react";
import { CurrentUser } from "@fider/models";
import { Heading } from "@fider/components";
import { SideMenu } from "./";

export abstract class AdminBasePage<P, S> extends React.Component<P, S> {
  public abstract id: string;
  public abstract name: string;
  public abstract icon: string;
  public abstract title: string;
  public abstract subtitle: string;
  public abstract content(): JSX.Element;

  public render() {
    return (
      <div id={this.id} className="page container">
        <Heading title={this.title} icon={this.icon} subtitle={this.subtitle} />

        <div className="row">
          <div className="col-lg-2">
            <SideMenu activeItem={this.name} />
          </div>
          <div className="col-lg-10">{this.content()}</div>
        </div>
      </div>
    );
  }
}
