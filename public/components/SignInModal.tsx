import React, { useState, useEffect } from "react";
import { Modal, SignInControl, LegalFooter } from "@fider/components/common";

interface SignInModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const SignInModal: React.StatelessComponent<SignInModalProps> = props => {
  const [confirmationAddress, setConfirmationAddress] = useState("");

  useEffect(() => {
    if (confirmationAddress) {
      setTimeout(() => setConfirmationAddress(""), 5000);
    }
  }, [confirmationAddress]);

  const onEmailSent = (email: string): void => {
    setConfirmationAddress(email);
  };

  const closeModal = () => {
    setConfirmationAddress("");
    props.onClose();
  };

  const content = confirmationAddress ? (
    <>
      <p>
        We have just sent a confirmation link to <b>{confirmationAddress}</b>. <br /> Click the link and you’ll be
        signed in.
      </p>
      <p>
        <a href="#" onClick={closeModal}>
          OK
        </a>
      </p>
    </>
  ) : (
    <SignInControl useEmail={true} onEmailSent={onEmailSent} />
  );

  return (
    <Modal.Window isOpen={props.isOpen} onClose={closeModal}>
      <Modal.Header>Sign in to raise your voice</Modal.Header>
      <Modal.Content>{content}</Modal.Content>
      <LegalFooter />
    </Modal.Window>
  );
};
