import "./ShareFeedback.scss"

import React, { useEffect, useRef, useState } from "react"
import { SignInControl, SignInSubmitResponse } from "@fider/components/common/SignInControl"
import { Modal, CloseIcon, Form, Button, TextArea, MultiImageUploader, Input, Icon, LegalFooter } from "@fider/components/common"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { actions, Failure, cache } from "@fider/services"
import { i18n } from "@lingui/core"
import { ImageUpload } from "@fider/models"
import IconAttach from "@fider/assets/images/heroicons-paperclip.svg"

interface ShareFeedbackProps {
  isOpen: boolean
  placeholder: string
  onClose: () => void
  draftPost?: {
    id: number
    code: string
    title: string
    description: string
  }
  draftAttachments?: string[]
}

export const ShareFeedback: React.FC<ShareFeedbackProps> = (props) => {
  const fider = useFider()
  const { isOpen, onClose, draftPost, draftAttachments } = props

  // State for the post form
  const [title, setTitle] = useState(draftPost?.title || "")
  const [description, setDescription] = useState(draftPost?.description || "")
  const [attachments, setAttachments] = useState<ImageUpload[]>(draftAttachments ? draftAttachments.map((bkey) => ({ bkey, url: "", remove: false })) : [])
  const [error, setError] = useState<Failure | undefined>(undefined)
  const titleRef = useRef<HTMLInputElement>()
  const [titleManuallyEdited, setTitleManuallyEdited] = useState(false)

  useEffect(() => {
    if (!titleManuallyEdited) {
      let newlineIndex = Math.min(description.indexOf("\n"), 80)
      if (newlineIndex == -1) {
        newlineIndex = 80
      }
      const autoTitle = description.substring(0, newlineIndex)
      handleTitleChange(autoTitle, false)
    }
  }, [description, titleManuallyEdited])

  // Handlers for post input changes
  const handleTitleChange = (value: string, isManualEdit = true) => {
    setTitle(value)
    // If this is a manual edit (not auto-generated from description),
    // mark the title as manually edited so we stop auto-populating
    if (isManualEdit) {
      setTitleManuallyEdited(true)
    }
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
              <TextArea
                label={props.placeholder}
                field="description"
                onChange={handleDescriptionChange}
                value={description}
                minRows={5}
                placeholder={i18n._({
                  id: "newpost.modal.description.placeholder",
                  message: "Tell us about your idea. Explain it fully, don't hold back, the more information the better.",
                })}
              />
              <MultiImageUploader
                field="attachments"
                maxUploads={3}
                bkeys={props.draftAttachments}
                onChange={handleAttachmentsChange}
                addImageButton={
                  <a className="flex items-center clickable">
                    <Icon sprite={IconAttach} height="18" width="18" />
                    <span className="ml-1">
                      <Trans id="newpost.modal.addimage">Add Images</Trans>
                    </span>
                  </a>
                }
              />
              <Input
                field="title"
                inputRef={titleRef}
                maxLength={255}
                label={i18n._({ id: "newpost.modal.title.label", message: "Give your suggestion a title" })}
                value={title}
                onChange={handleTitleChange}
                placeholder={i18n._({ id: "newpost.modal.title.placeholder", message: "Something short and snappy, sum it up in a few words" })}
              />
            </Form>
          </div>
        </div>
        <div className="c-share-feedback__content">
          <div className="c-share-feedback-signin">
            {/*
              Note: The email sign-in flow will save the post before signing in.
              For OAuth sign-in buttons, additional server-side changes would be needed
              to fully implement saving the post before OAuth redirect.
              Currently, only the email sign-in flow will work as expected.
            */}
            {!fider.session.isAuthenticated ? (
              <>
                <h2 className="text-title text-center mb-4">
                  <Trans id="modal.signin.header">Submit your feedback</Trans>
                </h2>
                <SignInControl
                  onSubmit={onSubmitFeedback}
                  onEmailSent={onEmailSent}
                  signInButtonText={i18n._({ id: "signin.message.email", message: "Continue with Email" })}
                  useEmail={true}
                  redirectTo={fider.settings.baseURL}
                />
              </>
            ) : (
              <div className="flex justify-center">
                <Button variant="primary" onClick={finaliseFeedback}>
                  <Trans id="modal.signin.header">Submit your feedback</Trans>
                </Button>
              </div>
            )}
          </div>
        </div>
        <LegalFooter />
      </Modal.Content>
    </Modal.Window>
  )
}
