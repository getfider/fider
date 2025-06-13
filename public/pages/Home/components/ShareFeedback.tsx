import "./ShareFeedback.scss"

import React, { useEffect, useRef, useState } from "react"
import { SignInControl, SignInSubmitResponse } from "@fider/components/common/SignInControl"
import { Modal, CloseIcon, Form, Button, TextArea, MultiImageUploader, Input, Icon, LegalFooter } from "@fider/components/common"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { actions, Failure, cache, querystring, classSet } from "@fider/services"
import { i18n } from "@lingui/core"
import { ImageUpload, Tag } from "@fider/models"
import IconAttach from "@fider/assets/images/heroicons-paperclip.svg"
import { SimilarPosts } from "../components/SimilarPosts"
import { TagsSelect } from "@fider/components/common/TagsSelect"

interface ShareFeedbackProps {
  isOpen: boolean
  placeholder: string
  onClose: () => void
  tags: Tag[]
}

export const CACHE_TITLE_KEY = "PostInput-Title"
export const CACHE_DESCRIPTION_KEY = "PostInput-Description"
export const CACHE_ATTACHMENT_KEY = "PostInput-Attachment"
export const CACHE_TAGS_KEY = "PostInput-Tags"

export const ShareFeedback: React.FC<ShareFeedbackProps> = (props) => {
  const fider = useFider()
  const { isOpen, onClose } = props

  const getCachedValue = (key: string): string => {
    return cache.session.get(key) || ""
  }

  const getDraftAttachments = (): ImageUpload[] => {
    const json = getCachedValue(CACHE_ATTACHMENT_KEY)
    return json.length ? JSON.parse(json) : []
  }

  const getTagsCachedValue = (): Tag[] => {
    if (!canEditTags) {
      return []
    }

    const cacheValue = getCachedValue(CACHE_TAGS_KEY)
    const urlValue = querystring.get("tags")
    const combined = [...cacheValue.split(","), ...urlValue.split(",")]
    const tagsAsStrings = Array.from(new Set(combined.map((s) => s.trim()).filter((s) => s.length > 0)))

    return props.tags.filter((tag) => tagsAsStrings.includes(tag.slug))
  }

  const canEditTags = fider.settings.postWithTags && props.tags.length > 0
  const [title, setTitle] = useState(getCachedValue(CACHE_TITLE_KEY))
  const [description, setDescription] = useState(getCachedValue(CACHE_DESCRIPTION_KEY))
  const [attachments, setAttachments] = useState<ImageUpload[]>(getDraftAttachments())
  const [tags, setTags] = useState(getTagsCachedValue())
  const [error, setError] = useState<Failure | undefined>(undefined)
  const titleRef = useRef<HTMLInputElement>()
  const descriptionRef = useRef<HTMLTextAreaElement>()
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

  useEffect(() => {
    if (isOpen && descriptionRef.current) {
      // Small delay to ensure modal is fully rendered
      setTimeout(() => {
        descriptionRef.current?.focus()
      }, 100)
    }
  }, [isOpen])

  // Handlers for post input changes
  const handleTitleChange = (value: string, isManualEdit = true) => {
    setTitle(value)
    cache.session.set(CACHE_TITLE_KEY, value)
    // If this is a manual edit (not auto-generated from description),
    // mark the title as manually edited so we stop auto-populating
    if (isManualEdit) {
      setTitleManuallyEdited(true)
    }
  }

  const handleTagsChanged = (newTags: Tag[]) => {
    cache.session.set(CACHE_TAGS_KEY, newTags.map((tag) => tag.slug).join(","))
    setTags(newTags)
  }

  const handleDescriptionChange = (value: string) => {
    setDescription(value)
    cache.session.set(CACHE_DESCRIPTION_KEY, value)
  }

  const handleAttachmentsChange = (images: ImageUpload[]) => {
    setAttachments(images)
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

  const finaliseFeedback = async () => {
    if (title) {
      const minDelay = new Promise((resolve) => setTimeout(resolve, 1000))

      const [result] = await Promise.all([
        actions.createPost(
          title,
          description,
          attachments,
          tags.map((tag) => tag.slug)
        ),
        minDelay,
      ])

      if (result.ok) {
        clearError()
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY, CACHE_ATTACHMENT_KEY, CACHE_TAGS_KEY)
        cache.session.set("POST_CREATED_SUCCESS", "true")
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
          <h1 className="text-large pb-6">
            <Trans id="newpost.modal.title">Share your feedback...</Trans>
          </h1>
          <div className="c-share-feedback-form">
            <Form error={error}>
              <TextArea
                label={props.placeholder}
                field="description"
                onChange={handleDescriptionChange}
                value={description}
                minRows={5}
                inputRef={descriptionRef}
                placeholder={i18n._({
                  id: "newpost.modal.description.placeholder",
                  message: "Tell us about your idea. Explain it fully, don't hold back, the more information the better.",
                })}
              />
              <SimilarPosts title={title} tags={props.tags} />
              <MultiImageUploader
                field="attachments"
                maxUploads={3}
                bkeys={attachments.filter((a) => a.bkey).map((a) => a.bkey ?? "")}
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
              {canEditTags && (
                <div className="c-form-field">
                  <label>
                    <Trans id="label.tags">Tags</Trans>
                  </label>
                  <div className={classSet({ "c-form-field": true })}>
                    <TagsSelect tags={props.tags} selectionChanged={handleTagsChanged} selected={tags} alwaysEditing={true} canEdit={true} />
                  </div>
                </div>
              )}
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
                <Button variant="primary" onClick={finaliseFeedback} disabled={title.length < 10}>
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
