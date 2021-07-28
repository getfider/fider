import React, { useState, useEffect } from "react"
import { Modal, SignInControl, LegalFooter } from "@fider/components"
import { Button } from "./common"
import { Trans } from "@lingui/macro"

interface SignInModalProps {
  isOpen: boolean
  onClose: () => void
}

export const SignInModal: React.StatelessComponent<SignInModalProps> = (props) => {
  const [confirmationAddress, setConfirmationAddress] = useState("")

  useEffect(() => {
    if (confirmationAddress) {
      setTimeout(() => setConfirmationAddress(""), 5000)
    }
  }, [confirmationAddress])

  const onEmailSent = (email: string): void => {
    setConfirmationAddress(email)
  }

  const closeModal = () => {
    setConfirmationAddress("")
    props.onClose()
  }

  const content = confirmationAddress ? (
    <>
      <p>
        <Trans id="signin.message.emailsent">
          We have just sent a confirmation link to <b>{confirmationAddress}</b>. <br /> Click the link and you’ll be signed in.
        </Trans>
      </p>
      <p>
        <Button variant="tertiary" onClick={closeModal}>
          <Trans id="action.ok">OK</Trans>
        </Button>
      </p>
    </>
  ) : (
    <SignInControl useEmail={true} onEmailSent={onEmailSent} />
  )

  return (
    <Modal.Window isOpen={props.isOpen} onClose={closeModal}>
      <Modal.Header>
        <Trans id="modal.signin.header">Sign in to post and vote</Trans>
      </Modal.Header>
      <Modal.Content>{content}</Modal.Content>
      <LegalFooter />
    </Modal.Window>
  )
}
