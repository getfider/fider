import * as React from "react";

import { mentions, Result } from "@fider/services";
import { User } from "@fider/models";
import { TextArea } from "@fider/components/common/";
import { classSet } from "@fider/services";
import { ValidationContext } from "../";
import { DisplayError, hasError } from "./DisplayError";
import { MentionsInput, Mention, MentionItem, SuggestionDataItem } from "react-mentions";

export interface MentionableTextAreaProps {
  field: string;
  label?: string;
  value?: string;
  disabled?: boolean;
  minRows?: number;
  placeholder?: string;
  onChange?: (value: string) => void;
  onKeyDown?: (event: React.KeyboardEvent<HTMLTextAreaElement> | React.KeyboardEvent<HTMLInputElement>) => void;
  inputRef?: React.RefObject<HTMLTextAreaElement>;
  onFocus?: React.FocusEventHandler<HTMLTextAreaElement>;
}

export class MentionableTextArea extends React.Component<MentionableTextAreaProps, {}> {
  private inputRef!: React.RefObject<HTMLTextAreaElement>;
  constructor(props: MentionableTextAreaProps) {
    super(props);
    this.fetchUsers = this.fetchUsers.bind(this);

    if (this.props.inputRef == null) {
      this.inputRef = React.createRef();
    } else {
      this.inputRef = this.props.inputRef;
    }
  }

  private onChange = (
    event: { target: { value: string } },
    newValue: string,
    newPlainTextValue: string,
    m: MentionItem[]
  ) => {
    console.log(newValue);
    if (this.props.onChange) {
      this.props.onChange(newValue);
    }
  };

  private fetchUsers = (query: string, callback: (data: SuggestionDataItem[]) => void) => {
    console.log("fetch users " + query);
    mentions.get(query).then((usersPromise: Result<User[]>) => {
      console.log(usersPromise);
      if (usersPromise.ok) {
        const items = usersPromise.data.map(element => {
          return {
            id: element.id,
            display: element.name
          };
        });
        console.log(items);
        callback(items);
      }
    });
  };

  public render() {
    //   return <MentionsInput
    //       value={this.props.value}
    //       onChange={this.onChange}
    //       onKeyDown={this.props.onKeyDown}
    //       placeholder={this.props.placeholder}
    //       inputRef={this.props.inputRef}
    //       >
    //   <Mention trigger="@" data={this.fetchUsers} />
    // </MentionsInput>

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
              {!!this.props.label && (
                <label htmlFor={this.inputRef.current ? this.inputRef.current.id : undefined}>{this.props.label}</label>
              )}
              <div >
                <MentionsInput
                  // id={`input-${this.props.field}`}
                  value={this.props.value}
                  onChange={this.onChange}
                  onKeyDown={this.props.onKeyDown}
                  placeholder={this.props.placeholder}
                  inputRef={this.inputRef}
                  markup="@__display__"
                >
                  <Mention trigger="@" data={this.fetchUsers} />
                </MentionsInput>

                 {/* <TextArea
                  value={this.props.value}
                  disabled={this.props.disabled}
                  
                  
                  
                  minRows={this.props.minRows || 3}
                  field="hi"
                  // inputRef={this.props.inputRef}
                  onFocus={this.props.onFocus}
                />  */}
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
