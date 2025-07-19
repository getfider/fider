import React, { useState, useEffect } from "react"
import { Modal, SignInControl, LegalFooter, TenantLogo } from "@fider/components"
import { Button, CloseIcon } from "./common"
import { Trans } from "@lingui/react/macro"
import { HStack, VStack } from "./layout"

interface SignInModalProps {
  isOpen: boolean
  onClose: () => void
}

export const SignInModal: React.FC<SignInModalProps> = (props) => {
  const [email, setEmail] = useState("")

  useEffect(() => {
    if (email) {
      setTimeout(() => setEmail(""), 5000)
    }
  }, [email])

  const onEmailSent = (value: string): void => {
    setEmail(value)
  }

  const closeModal = () => {
    setEmail("")
    props.onClose()
  }

  const content = email ? (
    <>
      <p>
        <Trans id="signin.message.emailsent">
          We have just sent a confirmation link to <b>{email}</b>. Click the link and you’ll be signed in.
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
        <VStack spacing={8}>
          <HStack justify="between">
            <TenantLogo size={24} useFiderIfEmpty={true} />
            <CloseIcon closeModal={closeModal} />
          </HStack>
          <p>
            <Trans id="modal.signin.header">Join the conversation</Trans>
          </p>
        </VStack>
      </Modal.Header>
      <Modal.Content>{content}</Modal.Content>
      <LegalFooter />
    </Modal.Window>
  )
}
