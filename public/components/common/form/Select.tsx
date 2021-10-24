import React from "react";
import { classSet } from "@fider/services";
import { ValidationContext } from "./Form";
import { DisplayError, hasError } from "./DisplayError";

export interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps {
  field: string;
  label?: string;
  maxLength?: number;
  defaultValue?: string;
  options: SelectOption[];
  onChange?: (option?: SelectOption) => void;
}

export const Select: React.FC<SelectProps> = (props) => {
  const onChange = (e: React.FormEvent<HTMLSelectElement>) => {
    const selectedOption = e.currentTarget.value
      ? props.options.find(({ value }) => value === e.currentTarget.value)
      : undefined;

    props.onChange?.(selectedOption);
  };

  const options = props.options.map((option) => (
    <option key={option.value} value={option.value}>
      {option.label}
    </option>
  ));

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <div
          className={classSet({
            "c-form-field": true,
            "m-error": hasError(props.field, ctx.error),
          })}
        >
          {!!props.label && <label htmlFor={`input-${props.field}`}>{props.label}</label>}
          <div className="c-form-field-wrapper">
            <select id={`input-${props.field}`} defaultValue={props.defaultValue} onChange={onChange}>
              {options}
            </select>
          </div>
          <DisplayError fields={[props.field]} error={ctx.error} />
          {props.children}
        </div>
      )}
    </ValidationContext.Consumer>
  );
};
