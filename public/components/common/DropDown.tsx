import * as React from "react";
import { findDOMNode } from "react-dom";
import { classSet } from "@fider/services";

/* TODO
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
  items: DropDownItem[];
  placeholder: string;
  onChange?: (item: DropDownItem) => void;
  renderItem?: (item: DropDownItem) => JSX.Element;
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
      selected: this.findItem(props.defaultValue, props.items),
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

  public findItem(value: any, items: DropDownItem[]): DropDownItem | undefined {
    for (const item of items) {
      if (item.value === value) {
        return item;
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

  public renderItem(item: DropDownItem) {
    const { label, value } = item;
    const isSelected = this.state.selected && (value === this.state.selected.value || value === this.state.selected);
    const className = classSet({
      "c-dropdown-item": true,
      "is-selected": isSelected
    });

    return (
      <div
        key={value}
        className={className}
        onMouseDown={this.setValue.bind(this, value, label)}
        onClick={this.setValue.bind(this, value, label)}
      >
        {this.props.renderItem ? this.props.renderItem(item) : label}
      </div>
    );
  }

  public buildItemList() {
    const items = this.props.items.map(item => {
      return this.renderItem(item);
    });

    return items.length ? items : <div className={`c-dropdown-noresults`}>No results found</div>;
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
        {this.state.isOpen && <div className="c-dropdown-menu">{this.buildItemList()}</div>}
      </div>
    );
  }
}
