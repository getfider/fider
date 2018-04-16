import * as React from "react";
import { SideMenu } from "./";

export abstract class AdminBasePage<P, S> extends React.Component<P, S> {
  public abstract name: string;
  public abstract icon: string;
  public abstract title: string;
  public abstract subtitle: string;
  public abstract content(): JSX.Element;

  public render() {
    return (
      <div className="page ui container">
        <h2 className="ui header">
          <i className={`circular ${this.icon} icon`} />
          <div className="content">
            {this.title}
            <div className="sub header">{this.subtitle}</div>
          </div>
        </h2>

        <div className="ui grid">
          <div className="three wide computer sixteen wide mobile column">
            <SideMenu activeItem={this.name} />
          </div>
          <div className="thirteen wide computer sixteen wide mobile column">{this.content()}</div>
        </div>
      </div>
    );
  }
}
