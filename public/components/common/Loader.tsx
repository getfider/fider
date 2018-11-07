import "./Loader.scss";

import React from "react";

interface LoaderState {
  show: boolean;
}

export class Loader extends React.Component<{}, LoaderState> {
  private unmounted: boolean;
  constructor(props: {}) {
    super(props);
    this.unmounted = false;
    this.state = {
      show: false
    };
  }

  public componentDidMount() {
    setTimeout(() => {
      if (!this.unmounted) {
        this.setState({
          show: true
        });
      }
    }, 500);
  }
  public componentWillUnmount() {
    this.unmounted = true;
  }

  public render() {
    return this.state.show && <div className="c-loader" />;
  }
}
