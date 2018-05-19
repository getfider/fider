import * as React from "react";
import { Failure, classSet } from "@fider/services";
import { Textarea, ValidationContext } from "../";
import { DisplayError, hasError } from "./DisplayError";

interface TextAreaProps {
  label: string;
  field: string;
  defaultValue?: string;
  disabled?: boolean;
  minRows?: number;
  placeholder?: string;
  onChange?: (value: string) => void;
}

interface TextAreaState {
  value: string;
  error?: Failure;
}

export class TextArea extends React.Component<TextAreaProps, TextAreaState> {
  constructor(props: TextAreaProps) {
    super(props);
    this.state = {
      value: props.defaultValue || ""
    };
  }

  private onChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    this.setState({ value: e.currentTarget.value }, () => {
      if (this.props.onChange) {
        this.props.onChange(this.state.value);
      }
    });
  };

  public render() {
    return (
      <ValidationContext.Consumer>
        {ctx => (
          <>
            <div
              className={classSet({
                "c-form-field": true,
                "m-error": hasError(this.props.field, ctx.error)
              })}
            >
              <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>
              <div className="c-form-field-wrapper">
                <Textarea
                  id={`input-${this.props.field}`}
                  disabled={this.props.disabled}
                  onChange={this.onChange}
                  value={this.state.value}
                  minRows={this.props.minRows}
                  placeholder={this.props.placeholder}
                />
              </div>
              <DisplayError fields={[this.props.field]} error={ctx.error} />
              {this.props.children}
            </div>
          </>
        )}
      </ValidationContext.Consumer>
    );
  }
}
