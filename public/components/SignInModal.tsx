import React from "react"
import { Modal, SignInControl, LegalFooter, TenantLogo } from "@fider/components"
import { CloseIcon } from "./common"
import { Trans } from "@lingui/react/macro"
import { HStack, VStack } from "./layout"

interface SignInModalProps {
  isOpen: boolean
  onClose: () => void
}

export const SignInModal: React.FC<SignInModalProps> = (props) => {
  const onCodeVerified = (result: { showProfileCompletion?: boolean; code?: string }): void => {
    if (result.showProfileCompletion && result.code) {
      // User needs to complete profile - redirect to profile completion page
      location.href = `/signin/complete?code=${encodeURIComponent(result.code)}`
    } else {
      // User is authenticated - close modal and reload to refresh the page
      props.onClose()
      location.reload()
    }
  }

  return (
    <Modal.Window isOpen={props.isOpen} onClose={props.onClose}>
      <Modal.Header>
        <VStack spacing={8}>
          <HStack justify="between">
            <TenantLogo size={24} useFiderIfEmpty={true} />
            <CloseIcon closeModal={props.onClose} />
          </HStack>
          <p>
            <Trans id="modal.signin.header">Join the conversation</Trans>
          </p>
        </VStack>
      </Modal.Header>
      <Modal.Content>
        <SignInControl useEmail={true} onCodeVerified={onCodeVerified} />
      </Modal.Content>
      <LegalFooter />
    </Modal.Window>
  )
}
