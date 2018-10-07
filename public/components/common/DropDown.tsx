import * as React from "react";
import { findDOMNode } from "react-dom";
import { classSet } from "@fider/services";

import "./DropDown.scss";

export interface DropDownItem {
  value: any;
  label: string;
}

export interface DropDownProps {
  className?: string;
  defaultValue?: any;
  items: DropDownItem[];
  placeholder?: string;
  inline?: boolean;
  header?: string;
  onChange?: (item: DropDownItem) => void;
  renderSelected?: (item?: DropDownItem) => JSX.Element;
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

    return (
      <div className="c-dropdown-menu">
        {this.props.header && <div className="c-dropdown-menu-header">{this.props.header}</div>}
        <div className="c-dropdown-menu-items">
          {items.length ? items : <div className={`c-dropdown-noresults`}>No results found</div>}
        </div>
      </div>
    );
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
    const text = this.state.selected ? (
      this.state.selected.label
    ) : (
      <span className="c-dropdown-placeholder">{this.props.placeholder}</span>
    );

    const dropdownClass = classSet({
      "c-dropdown": true,
      [`${this.props.className}`]: this.props.className,
      "is-open": this.state.isOpen,
      "is-inline": this.props.inline
    });

    return (
      <div className={dropdownClass}>
        <div onMouseDown={this.handleMouseDown} onTouchEnd={this.handleMouseDown}>
          {this.props.renderSelected ? (
            <div className="c-dropdown-text">{this.props.renderSelected(this.state.selected)}</div>
          ) : (
            <div className="c-dropdown-control">
              <div>{text}</div>
              <span className="c-dropdown-arrow" />
            </div>
          )}
        </div>
        {this.state.isOpen && this.buildItemList()}
      </div>
    );
  }
}
