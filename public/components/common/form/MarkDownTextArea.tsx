import * as React from "react";
import { Failure, classSet } from "@fider/services";
import { ValidationContext, MultiLineText } from "../";
import { DisplayError, hasError } from "./DisplayError";
import Textarea from "react-textarea-autosize";
import { Tab } from "semantic-ui-react";
import "./MarkDownTextArea.scss";

interface MarkDownTextAreaProps {
  label?: string;
  field: string;
  value?: string;
  disabled?: boolean;
  minRows?: number;
  placeholder?: string;
  onChange?: (value: string) => void;
  inputRef?: (node: HTMLTextAreaElement) => void;
  onFocus?: React.FocusEventHandler<HTMLTextAreaElement>;
}

export class MarkDownTextArea extends React.Component<MarkDownTextAreaProps, {}> {
  constructor(props: MarkDownTextAreaProps) {
    super(props);
    this.state = {
      value: props.value || ""
    };
  }

  private onChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    if (this.props.onChange) {
      this.setState({ value: e.currentTarget.value });
      this.props.onChange(e.currentTarget.value);
    }
  };

  public render() {
    const panes = [
      {
        menuItem: "Write",
        render: () => (
          <Tab.Pane>
            <Textarea
              id={`input-${this.props.field}`}
              disabled={this.props.disabled}
              onChange={this.onChange}
              value={this.props.value}
              minRows={this.props.minRows || 3}
              placeholder={this.props.placeholder}
              inputRef={this.props.inputRef}
              onFocus={this.props.onFocus}
            />
          </Tab.Pane>
        )
      },
      {
        menuItem: "Preview",
        render: () => (
          <Tab.Pane>
            <MultiLineText
              className={`markdown-preview-${this.props.minRows && this.props.minRows > 3 ? "full" : "simple"}`}
              text={this.props.value}
              style={this.props.minRows && this.props.minRows > 3 ? "full" : "simple"}
            />
          </Tab.Pane>
        )
      }
    ];
    return (
      <ValidationContext.Consumer>
        {ctx => (
          <>
            <div
              className={classSet({
                "c-form-field": true,
                "markdown-text-area": true,
                "m-error": hasError(this.props.field, ctx.error)
              })}
            >
              {!!this.props.label && <label htmlFor={`input-${this.props.field}`}>{this.props.label}</label>}
              <div className="c-form-field-wrapper">
                <Tab panes={panes} />
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
