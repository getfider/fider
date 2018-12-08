import "./Hint.scss";

import React from "react";
import { FaTimes } from "react-icons/fa";

import { cache } from "@fider/services";

interface HintProps {
  permanentCloseKey?: string;
  condition?: boolean;
}

interface HintState {
  isClosed: boolean;
}

export class Hint extends React.Component<HintProps, HintState> {
  private cacheKey: string | undefined;

  constructor(props: HintProps) {
    super(props);
    this.cacheKey = this.props.permanentCloseKey ? `Hint-Closed-${this.props.permanentCloseKey}` : undefined;
    this.state = {
      isClosed: this.cacheKey ? cache.local.has(this.cacheKey) : false
    };
  }

  private close = () => {
    if (this.cacheKey) {
      cache.local.set(this.cacheKey, "true");
    }

    this.setState({
      isClosed: true
    });
  };

  public render() {
    if (this.props.condition === false || this.state.isClosed) {
      return null;
    }

    return (
      <p className="c-hint">
        <strong>HINT:</strong> {this.props.children}
        {this.cacheKey && <FaTimes onClick={this.close} className="close" />}
      </p>
    );
  }
}
