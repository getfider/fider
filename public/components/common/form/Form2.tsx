import "./Form2.scss";

import * as React from "react";
import { Failure } from "@fider/services";
import { Button } from "@fider/components";

interface ValidationContext {
  error?: Failure;
}

interface FormProps {
  error?: Failure;
}

export const ValidationContext = React.createContext<ValidationContext>({});

export class Form2 extends React.Component<FormProps, {}> {
  constructor(props: FormProps) {
    super(props);
    this.state = {
      error: this.props.error
    };
  }

  public render() {
    return (
      <form autoComplete="off" className="c-form">
        <ValidationContext.Provider value={{ error: this.props.error }}>
          {this.props.children}
        </ValidationContext.Provider>
      </form>
    );
  }
}
