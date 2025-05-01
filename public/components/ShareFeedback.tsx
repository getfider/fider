import "./ShareFeedback.scss"

import React from "react"
import { PostInput } from "@fider/pages/Home/components/PostInput"
import { SignInControl } from "@fider/components/common/SignInControl"
import { Modal, CloseIcon } from "@fider/components/common"
import { useFider } from "@fider/hooks"
import { i18n } from "@lingui/core"

interface ShareFeedbackProps {
  isOpen: boolean
  onClose: () => void
}

export const ShareFeedback: React.FC<ShareFeedbackProps> = (props) => {
  const fider = useFider()
  const { isOpen, onClose } = props

  const handleTitleChanged = () => {
    // Handle title changes if needed
  }

  return (
    <Modal.Window className="c-share-feedback" isOpen={isOpen} onClose={onClose} size="fullscreen" center={false}>
      <Modal.Header>
        <div className="flex flex-items-center justify-end">
          {/* <h2 className="text-title">Share your feedback...</h2> */}
          <CloseIcon closeModal={onClose} />
        </div>
      </Modal.Header>
      <Modal.Content>
        <div className="c-share-feedback__content">
          <h2 className="text-title pb-6">Share your feedback...</h2>
          <div className="c-share-feedback-form">
            <p className="text-sm text-muted mb-4">Tell us what you&apos;d like to see in Fider</p>
            <PostInput
              placeholder={i18n._("home.postinput.placeholder", { message: "Something short and snappy, sum it up a few words" })}
              onTitleChanged={handleTitleChanged}
            />
          </div>
          <div className="c-share-feedback-signin">
            <h3 className="text-subtitle text-center mb-4">Submit your feedback</h3>
            <SignInControl useEmail={true} redirectTo={fider.settings.baseURL} />
          </div>
        </div>
      </Modal.Content>
    </Modal.Window>
  )
}
