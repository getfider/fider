import * as React from "react";
import { number } from "prop-types";

interface SuggesionBoxProps {
  top: number;
  left: number;
  shown: boolean;
  data: string[];
  itemSelected?: number;
  onItemClick?: (item: string, itemNumber: number) => void;
}

export class SuggestionBox extends React.Component<SuggesionBoxProps, {}> {
  constructor(props: SuggesionBoxProps) {
    super(props);
  }

  private onItemClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (this.props.onItemClick) {
      const index = Number(e.currentTarget.getAttribute("key"));
      const data = this.props.data[index];
      this.props.onItemClick(data, index);
    }
  };

  public render() {
    return (
      <div
        style={{
          position: "absolute",
          width: "200px",
          borderRadius: "6px",
          background: "white",
          boxShadow: "rgba(0, 0, 0, 0.4) 0px 1px 4px",

          display: this.props.shown ? "block" : "none",
          top: this.props.top,
          left: this.props.left
        }}
      >
        {this.props.data.map((item, index) => (
          <div
            style={{
              padding: "5px 5px",
              background: index === this.props.itemSelected ? "#eee" : ""
            }}
            key={index}
            onClick={this.onItemClick}
          >
            {item}
          </div>
        ))}
      </div>
    );
  }
}
