import './Modal.scss';

import * as React from 'react';
import { SystemSettings, CurrentUser, Tenant } from '@fider/models';
import { SignInControl, EnvironmentInfo, Gravatar } from '@fider/components/common';
import { page, actions } from '@fider/services';

interface ModalProps {
  isOpen: boolean;
}

interface ModalState {
  isOpen: boolean;
}

export class Modal extends React.Component<ModalProps, ModalState> {
  constructor(props: ModalProps) {
    super(props);
    this.state = {
      isOpen: props.isOpen,
    };

    this.keyDown = this.keyDown.bind(this);
    this.close = this.close.bind(this);
  }

  public componentDidMount() {
    document.addEventListener('keydown', this.keyDown, false);
  }

  public componentWillUnmount() {
    document.removeEventListener('keydown', this.keyDown, false);
  }

  public componentWillReceiveProps(props: ModalProps) {
    this.setState({ isOpen: props.isOpen });
  }

  private keyDown(event: KeyboardEvent) {
    if (event.keyCode === 27) { // ESC
      this.close();
    }
  }

  private close() {
    this.setState({ isOpen: false });
  }

  public render() {
    if (!this.state.isOpen) {
      return null;
    }

    return (
      <div
        className="c-modal__dimmer"
        onClick={this.close}
      >
        <div className="c-modal__window">
          The Modal Content
        </div>
      </div>
    );
  }
}
