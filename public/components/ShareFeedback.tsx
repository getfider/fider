import "./ShareFeedback.scss"

import React from "react"
import { PostInputAnonymous } from "@fider/pages/Home/components/PostInputAnonymous"
import { SignInControl } from "@fider/components/common/SignInControl"
import { Modal, CloseIcon } from "@fider/components/common"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"

interface ShareFeedbackProps {
  isOpen: boolean
  placeholder: string
  onClose: () => void
}

export const ShareFeedback: React.FC<ShareFeedbackProps> = (props) => {
  const fider = useFider()
  const { isOpen, onClose } = props

  return (
    <Modal.Window className="c-share-feedback" isOpen={isOpen} onClose={onClose} size="fullscreen" center={false}>
      <Modal.Header>
        <div className="flex flex-items-center justify-end">
          <CloseIcon closeModal={onClose} />
        </div>
      </Modal.Header>
      <Modal.Content>
        <div className="c-share-feedback__content mb-4">
          <h2 className="text-title pb-6">
            <Trans id="newpost.modal.title">Share your feedback...</Trans>
          </h2>
          <div className="c-share-feedback-form">
            <PostInputAnonymous placeholder={props.placeholder} />
          </div>
        </div>
        <div className="c-share-feedback__content">
          <div className="c-share-feedback-signin">
            <h2 className="text-title text-center mb-4">Submit your feedback</h2>
            <SignInControl
              signInButtonText={i18n._({ id: "signin.message.email", message: "Continue with Email" })}
              useEmail={true}
              redirectTo={fider.settings.baseURL}
            />
          </div>
        </div>
      </Modal.Content>
    </Modal.Window>
  )
}
