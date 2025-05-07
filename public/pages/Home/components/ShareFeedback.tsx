import "./ShareFeedback.scss"

import React, { useState } from "react"
import { ShareFeedbackForm } from "@fider/pages/Home/components/ShareFeedbackForm"
import { SignInControl, SignInSubmitResponse } from "@fider/components/common/SignInControl"
import { Modal, CloseIcon, Form, Button } from "@fider/components/common"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { actions, Failure, cache } from "@fider/services"
import { i18n } from "@lingui/core"
import { ImageUpload } from "@fider/models"

interface ShareFeedbackProps {
  isOpen: boolean
  placeholder: string
  onClose: () => void
}

export const ShareFeedback: React.FC<ShareFeedbackProps> = (props) => {
  const fider = useFider()
  const { isOpen, onClose } = props

  // State for the post form
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)

  // Handlers for post input changes
  const handleTitleChange = (value: string) => {
    setTitle(value)
  }

  const handleDescriptionChange = (value: string) => {
    setDescription(value)
  }

  const handleAttachmentsChange = (value: ImageUpload[]) => {
    setAttachments(value)
  }

  const onSubmitFeedback = async (): Promise<SignInSubmitResponse> => {
    // Always try to save the post first
    const postResult = await actions.createAnonymousPost(title, description, attachments)

    if (postResult.ok) {
      // Post saved successfully, now proceed with sign in
      return { ok: true, code: postResult.data.code }
    } else if (postResult.error) {
      setError(postResult.error)
    }
    return { ok: false }
  }

  const clearError = () => setError(undefined)
  const CACHE_TITLE_KEY = "PostInput-Title"
  const CACHE_DESCRIPTION_KEY = "PostInput-Description"

  const finaliseFeedback = async () => {
    if (title) {
      const result = await actions.createPost(title, description, attachments)
      if (result.ok) {
        clearError()
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY)
        location.href = `/posts/${result.data.number}/${result.data.slug}`
      } else if (result.error) {
        setError(result.error)
      }
    }
  }

  const onEmailSent = (email: string) => {
    window.location.href = "/loginemailsent?email=" + encodeURIComponent(email)
  }

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
            <Form error={error}>
              <ShareFeedbackForm
                placeholder={props.placeholder}
                onTitleChange={handleTitleChange}
                onDescriptionChange={handleDescriptionChange}
                onAttachmentsChange={handleAttachmentsChange}
                error={error}
              />
            </Form>
          </div>
        </div>
        <div className="c-share-feedback__content">
          <div className="c-share-feedback-signin">
            <h2 className="text-title text-center mb-4">Submit your feedback</h2>
            {/*
              Note: The email sign-in flow will save the post before signing in.
              For OAuth sign-in buttons, additional server-side changes would be needed
              to fully implement saving the post before OAuth redirect.
              Currently, only the email sign-in flow will work as expected.
            */}
            {!fider.session.isAuthenticated ? (
              <SignInControl
                onSubmit={onSubmitFeedback}
                onEmailSent={onEmailSent}
                signInButtonText={i18n._({ id: "signin.message.email", message: "Continue with Email" })}
                useEmail={true}
                redirectTo={fider.settings.baseURL}
              />
            ) : (
              <Button variant="primary" onClick={finaliseFeedback}>
                <Trans id="action.submit">Submit</Trans>
              </Button>
            )}
          </div>
        </div>
      </Modal.Content>
    </Modal.Window>
  )
}
