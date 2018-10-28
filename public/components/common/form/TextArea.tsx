import * as React from "react";
import { Failure, classSet } from "@fider/services";
import { ValidationContext } from "../";
import { DisplayError, hasError } from "./DisplayError";
import Textarea from "react-textarea-autosize";
import getCaretCoordinates = require('textarea-caret');
import { runInThisContext } from 'vm';


export interface TextAreaTriggerStart {
  top: number;
  left: number;
  key: string;
}

interface TextAreaProps {
  label?: string;
  field: string;
  value?: string;
  disabled?: boolean;
  minRows?: number;
  placeholder?: string;
  onChange?: (value: string) => void;
  onTriggerStart?: (e : TextAreaTriggerStart) => void;
  onTriggerEnd?: () => void;
  onTriggerChange?: (text: string) => void;
  onTriggerArrow?: (arrow: string) => void;
  onTriggerSelected?: (triggerPosition: number, cursorPosition : number) => void;
  inputRef?: (node: HTMLTextAreaElement) => void;
  onFocus?: React.FocusEventHandler<HTMLTextAreaElement>;
}

interface TextAreaState {
  suggestionTriggered: boolean;
  suggestionTriggerPosition: number;
}


export class TextArea extends React.Component<TextAreaProps, TextAreaState> {
  constructor(props: TextAreaProps) {
    super(props);
  }

  componentDidMount = () => {
    this.setState({suggestionTriggered : false})
  }

  private onChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    if (this.props.onChange) {
      this.props.onChange(e.currentTarget.value);
    }
  };

  private onKeyDown = (e : React.KeyboardEvent<HTMLTextAreaElement>) => {
    const {onTriggerEnd, onTriggerChange, onTriggerStart, onTriggerArrow, onTriggerSelected} = this.props;

    const {suggestionTriggered, suggestionTriggerPosition} = this.state;

    const {key, which, currentTarget} = e
    if (onTriggerStart ){
      const { selectionStart } = e.currentTarget;

      if( suggestionTriggered ){
        // Backspace handling
        if (which === 8 && selectionStart <= suggestionTriggerPosition) {
          this.setState({
            suggestionTriggered : false,
            suggestionTriggerPosition: 0
          });

          onTriggerEnd && setTimeout(() =>{ onTriggerEnd() }, 0);
        }
        // Down arrow handling
        else if (which === 40 && onTriggerArrow) {
          e.preventDefault();
          setTimeout(() => {
            onTriggerArrow("down");
          }, 0);
        }
        // Up arrow handling
        else if (which === 38 && onTriggerArrow) {
          e.preventDefault();
          setTimeout(() => {
            onTriggerArrow("up");
          }, 0);
        }
        // Enter handling
        else if (which === 13 && onTriggerSelected) {
          console.log("selected");
          const capturedText = currentTarget.value.substr(suggestionTriggerPosition, selectionStart);
          e.preventDefault();
          setTimeout(() => {
            onTriggerSelected(suggestionTriggerPosition, selectionStart);
          }, 0);
          this.setState({
            suggestionTriggered : false,
            suggestionTriggerPosition: 0
          });
        }
        else {
          onTriggerChange && setTimeout(()=> {
            const capturedText = currentTarget.value.substr(suggestionTriggerPosition, selectionStart);
            onTriggerChange(capturedText)
            });
        }
      }
      else {
        if (key === '@'){
          const caretCoordinates =  getCaretCoordinates(e.currentTarget, selectionStart)

          this.setState({ suggestionTriggered : true,
                          suggestionTriggerPosition: selectionStart + 1
                        });

          setTimeout(()=> {
            onTriggerStart( {
              key : key,
              left : caretCoordinates.left,
              top : caretCoordinates.top + caretCoordinates.height + 5
            });
          });
        }

     }
    }
  }

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
                  onKeyDown={this.onKeyDown}
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
