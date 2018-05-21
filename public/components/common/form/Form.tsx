import "./Form.scss";

import * as React from "react";
import { Failure, classSet } from "@fider/services";
import { Button, DisplayError } from "@fider/components";

interface ValidationContext {
  error?: Failure;
}

interface FormProps {
  className?: string;
  size?: "mini" | "normal";
  error?: Failure;
}

export const ValidationContext = React.createContext<ValidationContext>({});

export class Form extends React.Component<FormProps, {}> {
  constructor(props: FormProps) {
    super(props);
    this.state = {
      error: this.props.error
    };
  }

  public render() {
    const className = classSet({
      "c-form": true,
      [this.props.className!]: this.props.className,
      [`m-${this.props.size}`]: this.props.size
    });
    return (
      <form autoComplete="off" className={className}>
        <DisplayError error={this.props.error} />
        <ValidationContext.Provider value={{ error: this.props.error }}>
          {this.props.children}
        </ValidationContext.Provider>
      </form>
    );
  }
}
