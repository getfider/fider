import "./Modal.scss";

import React from "react";
import ReactDOM from "react-dom";
import { classSet } from "@fider/services";

interface ModalWindowProps {
  className?: string;
  isOpen: boolean;
  size?: "small" | "large";
  canClose?: boolean;
  center?: boolean;
  onClose?: () => void;
}

interface ModalWindowState {
  isOpen: boolean;
}

interface ModalFooterProps {
  align?: "left" | "center" | "right";
  children?: React.ReactNode;
}

class ModalWindow extends React.Component<ModalWindowProps, ModalWindowState> {
  private root?: HTMLElement;

  constructor(props: ModalWindowProps) {
    super(props);
    this.state = {
      isOpen: this.props.isOpen
    };
  }

  public static defaultProps: Partial<ModalWindowProps> = {
    size: "small",
    canClose: true,
    center: true
  };

  public componentWillUpdate(nextProps: ModalWindowProps, nextState: ModalWindowState) {
    if (nextState.isOpen) {
      document.addEventListener("keydown", this.keyDown, false);
    } else {
      document.removeEventListener("keydown", this.keyDown, false);
    }
  }

  public componentWillReceiveProps(nextProps: ModalWindowProps) {
    this.setState({
      isOpen: nextProps.isOpen
    });
  }

  private keyDown = (event: KeyboardEvent) => {
    if (event.keyCode === 27) {
      // ESC
      this.close();
    } else if (event.keyCode === 9) {
      // TAB
      event.preventDefault();
    }
  };

  private close = () => {
    if (this.props.canClose) {
      this.setState({ isOpen: false });
      if (this.props.onClose) {
        this.props.onClose();
      }
    }
  };

  private swallow = (evt: React.MouseEvent<HTMLDivElement>) => {
    evt.stopPropagation();
  };

  private getContainer = (): HTMLElement => {
    if (!this.root) {
      this.root = document.getElementById("root-modal")!;
    }
    return this.root;
  };

  public render() {
    if (!this.state.isOpen) {
      document.body.style.overflow = "";
      return null;
    }

    document.body.style.overflow = "hidden";

    const className = classSet({
      "c-modal-window": true,
      [`${this.props.className}`]: !!this.props.className,
      "m-center": this.props.center,
      [`m-${this.props.size}`]: true
    });

    return ReactDOM.createPortal(
      <div aria-disabled={true} className="c-modal-dimmer" onClick={this.close}>
        <div className={className} onClick={this.swallow}>
          {this.props.children}
        </div>
      </div>,
      this.getContainer()
    );
  }
}

export const Modal = {
  Window: ModalWindow,
  Header: (props: { children: React.ReactNode }) => {
    return <div className="c-modal-header">{props.children}</div>;
  },
  Content: (props: { children: React.ReactNode }) => {
    return <div className="c-modal-content">{props.children}</div>;
  },
  Footer: (props: ModalFooterProps) => {
    const align = props.align || "right";
    const className = classSet({
      "c-modal-footer": true,
      [`m-${align}`]: true
    });
    return <div className={className}>{props.children}</div>;
  }
};
