import * as React from "react";
import { findDOMNode } from "react-dom";
import { classSet } from "@fider/services";

/* TODO
- Refactor CSS Classes
- CSS Similar to Semantic UI
- Default Option Render
- Custom Option Render
- Default Header Render
- Custom Header Render
- Option List Header
*/

import "./DropDown.scss";

export interface DropDownItem {
  value: any;
  label: string;
}

export interface DropDownProps {
  defaultValue?: any;
  options: DropDownItem[];
  placeholder: string;
  onChange?: (item: DropDownItem) => void;
}

export interface DropDownState {
  isOpen: boolean;
  selected?: DropDownItem;
}

export class DropDown extends React.Component<DropDownProps, DropDownState> {
  private mounted = false;

  constructor(props: DropDownProps) {
    super(props);
    this.state = {
      selected: this.parseValue(props.defaultValue, props.options),
      isOpen: false
    };
    this.mounted = true;
  }

  public componentDidMount() {
    document.addEventListener("click", this.handleDocumentClick, false);
    document.addEventListener("touchend", this.handleDocumentClick, false);
  }

  public componentWillUnmount() {
    this.mounted = false;
    document.removeEventListener("click", this.handleDocumentClick, false);
    document.removeEventListener("touchend", this.handleDocumentClick, false);
  }

  public handleMouseDown = (event: any) => {
    if (event.type === "mousedown" && event.button !== 0) {
      return;
    }

    event.stopPropagation();
    event.preventDefault();

    this.setState({
      isOpen: !this.state.isOpen
    });
  };

  public parseValue(value: any, options: DropDownItem[]): DropDownItem | undefined {
    for (const opt of options) {
      if (opt.value === value) {
        return opt;
      }
    }
    return undefined;
  }

  public setValue(value: any, label: string) {
    const newState = {
      selected: {
        value,
        label
      },
      isOpen: false
    };
    this.fireChangeEvent(newState);
    this.setState(newState);
  }

  public fireChangeEvent(newState: DropDownState) {
    if (newState.selected && newState.selected !== this.state.selected && this.props.onChange) {
      this.props.onChange(newState.selected);
    }
  }

  public renderOption(option: DropDownItem) {
    let value = option.value;
    if (typeof value === "undefined") {
      value = option.label || option;
    }
    const label = option.label || option.value || option;

    const className = classSet({
      "c-dropdown-option": true,
      "is-selected": this.state.selected && (value === this.state.selected.value || value === this.state.selected)
    });

    return (
      <div
        key={value}
        className={className}
        onMouseDown={this.setValue.bind(this, value, label)}
        onClick={this.setValue.bind(this, value, label)}
      >
        {label}
      </div>
    );
  }

  public buildMenu() {
    const ops = this.props.options.map(option => {
      return this.renderOption(option);
    });

    return ops.length ? ops : <div className={`c-dropdown-noresults`}>No options found</div>;
  }

  public handleDocumentClick = (event: any) => {
    if (this.mounted) {
      const node = findDOMNode(this);
      if (node && !node.contains(event.target)) {
        if (this.state.isOpen) {
          this.setState({ isOpen: false });
        }
      }
    }
  };

  public render() {
    const displayLabel = this.state.selected ? this.state.selected.label : this.props.placeholder;

    const dropdownClass = classSet({
      "c-dropdown": true,
      "is-open": this.state.isOpen
    });

    return (
      <div className={dropdownClass}>
        <div className="c-dropdown-control" onMouseDown={this.handleMouseDown} onTouchEnd={this.handleMouseDown}>
          <div>{displayLabel}</div>
          <span className="c-dropdown-arrow" />
        </div>
        {this.state.isOpen && <div className="c-dropdown-menu">{this.buildMenu()}</div>}
      </div>
    );
  }
}
