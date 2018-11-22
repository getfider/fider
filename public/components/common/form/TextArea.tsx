import React from "react";
import { classSet } from "@fider/services";
import { ValidationContext } from "../";
import { DisplayError, hasError } from "./DisplayError";
import Textarea from "react-textarea-autosize";

import { runInThisContext } from 'vm';

export interface TextAreaProps {
  field: string;
  label?: string;
  value?: string;
  disabled?: boolean;
  minRows?: number;
  placeholder?: string;
  onChange?: (value: string) => void;
  onKeyDown?: (e : React.KeyboardEvent<HTMLTextAreaElement>) => void;
  inputRef?: (node: HTMLTextAreaElement) => void;
  onFocus?: React.FocusEventHandler<HTMLTextAreaElement>;
}

export class TextArea extends React.Component<TextAreaProps, {}> {
  constructor(props: TextAreaProps) {
    super(props);
  }

  private onChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    if (this.props.onChange) {
      this.props.onChange(e.currentTarget.value);
    }
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
              {!!this.props.label && <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>}
              <div className="c-form-field-wrapper">
                <Textarea
                  id={`input-${this.props.field}`}
                  disabled={this.props.disabled}
                  onChange={this.onChange}
                  onKeyDown={this.props.onKeyDown}
                  value={this.props.value}
                  minRows={this.props.minRows || 3}
                  placeholder={this.props.placeholder}
                  inputRef={this.props.inputRef}
                  onFocus={this.props.onFocus}
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
