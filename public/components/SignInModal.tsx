import React from "react";
import { Modal, SignInControl, LegalFooter } from "@fider/components/common";

interface SignInModalProps {
  isOpen: boolean;
}

interface SignInModalState {
  isOpen: boolean;
  email: string;
  emailSent: boolean;
}

export class SignInModal extends React.Component<SignInModalProps, SignInModalState> {
  constructor(props: SignInModalProps) {
    super(props);
    this.state = {
      isOpen: this.props.isOpen,
      email: "",
      emailSent: false
    };

    this.onEmailSent = this.onEmailSent.bind(this);
  }

  public componentWillReceiveProps(nextProps: SignInModalProps) {
    this.setState({
      isOpen: nextProps.isOpen
    });
  }

  private onEmailSent(email: string): void {
    this.setState({ email, emailSent: true }, () => {
      setTimeout(() => {
        this.setState({ email: "", emailSent: false });
      }, 5000);
    });
  }

  private closeModal = () => {
    this.setState({ isOpen: false, emailSent: false });
  };

  public render() {
    const content = this.state.emailSent ? (
      <>
        <p>
          We have just sent a confirmation link to <b>{this.state.email}</b>. <br /> Click the link and you’ll be signed
          in.
        </p>
        <p>
          <a href="#" onClick={this.closeModal}>
            OK
          </a>
        </p>
      </>
    ) : (
      <SignInControl useEmail={true} onEmailSent={this.onEmailSent} />
    );

    return (
      <Modal.Window isOpen={this.state.isOpen}>
        <Modal.Header>Sign in to raise your voice</Modal.Header>
        <Modal.Content>{content}</Modal.Content>
        <LegalFooter />
      </Modal.Window>
    );
  }
}
