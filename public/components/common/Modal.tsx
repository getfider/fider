import './Modal.scss';

import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { SystemSettings, CurrentUser, Tenant } from '@fider/models';
import { SignInControl, EnvironmentInfo, Gravatar } from '@fider/components/common';
import { page, actions, classSet } from '@fider/services';

interface ModalWindowProps {
  isOpen: boolean;
  size?: 'small' | 'large';
  canClose?: boolean;
  center?: boolean;
}

interface ModalWindowState {
  isOpen: boolean;
}

class ModalWindow extends React.Component<ModalWindowProps, ModalWindowState> {
  private root?: HTMLElement;

  constructor(props: ModalWindowProps) {
    super(props);
    this.state = {
      isOpen: this.props.isOpen,
    };

    this.keyDown = this.keyDown.bind(this);
    this.close = this.close.bind(this);
    this.swallow = this.swallow.bind(this);
  }

  public static defaultProps: Partial<ModalWindowProps> = {
    size: 'small',
    canClose: true,
    center: true,
  };

  public componentDidMount() {
    document.addEventListener('keydown', this.keyDown, false);
  }

  public componentWillUnmount() {
    document.removeEventListener('keydown', this.keyDown, false);
  }

  public componentWillReceiveProps(nextProps: ModalWindowProps) {
    this.setState({
      isOpen: nextProps.isOpen
    });
  }

  private keyDown(event: KeyboardEvent) {
    if (event.keyCode === 27) { // ESC
      this.close();
    } else if (event.keyCode === 9) { // TAB
      event.preventDefault();
    }
  }

  private close() {
    if (this.props.canClose) {
      this.setState({ isOpen: false });
    }
  }

  private swallow(evt: React.MouseEvent<HTMLDivElement>) {
    evt.stopPropagation();
  }

  private getContainer(): HTMLElement {
    if (!this.root) {
      this.root = document.getElementById('root-modal')!;
    }
    return this.root;
  }

  public render() {
    if (!this.state.isOpen) {
      return null;
    }

    const className = classSet({
      'c-modal__window': true,
      'center': this.props.center,
      [this.props.size!]: true,
    });

    return ReactDOM.createPortal((
      <div
        aria-disabled={true}
        className="c-modal__dimmer"
        onClick={this.close}
      >
        <div
          className={className}
          onClick={this.swallow}
        >
          {this.props.children}
        </div>
      </div>
    ), this.getContainer());
  }
}

export const Modal = {
  Window: ModalWindow,
  Header: (props: { children: React.ReactNode }) => {
    return (
      <div className="c-modal__header">
        {props.children}
      </div>
    );
  },
  Content: (props: { children: React.ReactNode }) => {
    return (
      <div className="c-modal__content">
        {props.children}
      </div>
    );
  },
  Footer: (props: { children: React.ReactNode }) => {
    return (
      <div className="c-modal__footer">
        {props.children}
      </div>
    );
  }
};
