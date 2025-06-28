import "./ShareFeedback.scss"

import React, { useEffect, useRef, useState } from "react"
import { SignInControl, SignInSubmitResponse } from "@fider/components/common/SignInControl"
import { Modal, CloseIcon, Form, Button, Input, LegalFooter } from "@fider/components/common"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { actions, Failure, cache, querystring, classSet } from "@fider/services"
import { i18n } from "@lingui/core"
import { ImageUpload, Tag } from "@fider/models"
import { SimilarPosts } from "../components/SimilarPosts"
import { TagsSelect } from "@fider/components/common/TagsSelect"
import CommentEditor from "@fider/components/common/form/CommentEditor"

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

  const getTitleManuallyEditedValue = (): boolean => {
    // If the cached title deviates from the description, it means the user manually edited it
    return getCachedValue(CACHE_TITLE_KEY) !== getCachedValue(CACHE_DESCRIPTION_KEY)
  }

  const canEditTags = fider.settings.postWithTags && props.tags.length > 0
  const [title, setTitle] = useState(getCachedValue(CACHE_TITLE_KEY))
  const [description, setDescription] = useState(getCachedValue(CACHE_DESCRIPTION_KEY))
  const [attachments, setAttachments] = useState<ImageUpload[]>(getDraftAttachments())
  const [tags, setTags] = useState(getTagsCachedValue())
  const [error, setError] = useState<Failure | undefined>(undefined)
  const titleRef = useRef<HTMLInputElement>()
  const editorRef = useRef<HTMLDivElement>(null)
  const [titleManuallyEdited, setTitleManuallyEdited] = useState(getTitleManuallyEditedValue())
  const [isInitialMount, setIsInitialMount] = useState(true)

  useEffect(() => {
    setIsInitialMount(false)
  }, [])

  // Handle browser back button
  useEffect(() => {
    if (isOpen) {
      // Push a new state when modal opens
      window.history.pushState({ modalOpen: true }, "", window.location.href)

      const handlePopState = () => {
        // If we're going back and the modal is open, close it
        if (isOpen) {
          onClose()
        }
      }

      window.addEventListener("popstate", handlePopState)

      return () => {
        window.removeEventListener("popstate", handlePopState)
      }
    }
  }, [isOpen, onClose])

  // Handle modal close - go back in history if we pushed a state
  const handleClose = () => {
    // Check if we can go back (and if the previous state was pushed by us)
    if (window.history.state?.modalOpen) {
      window.history.back()
    } else {
      onClose()
    }
  }

  useEffect(() => {
    if (!titleManuallyEdited && !isInitialMount) {
      let newlineIndex = Math.min(description.indexOf("\n"), 80)
      if (newlineIndex == -1) {
        newlineIndex = 80
      }
      const autoTitle = description.substring(0, newlineIndex)
      handleTitleChange(autoTitle, false)
    }
  }, [description, titleManuallyEdited])

  useEffect(() => {
    if (isOpen && editorRef.current) {
      // Small delay to ensure modal is fully rendered
      setTimeout(() => {
        // Focus the editor
        const editorContent = editorRef.current?.querySelector(".ProseMirror")
        if (editorContent) {
          ;(editorContent as HTMLElement).focus()
        }
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
    cache.session.set(CACHE_DESCRIPTION_KEY, value)
    setDescription(value)
  }

  const handleImageUploaded = (image: ImageUpload) => {
    setAttachments((prev) => {
      // If this is a removal request, find and mark the attachment for removal
      if (image.remove && image.bkey) {
        return prev.map((att) => (att.bkey === image.bkey ? { ...att, remove: true } : att))
      }
      // Otherwise add the new upload
      const newAttachments = [...prev, image]

      // Update the cache
      cache.session.set(CACHE_ATTACHMENT_KEY, JSON.stringify(newAttachments))

      return newAttachments
    })
  }

  const onSubmitFeedback = async (): Promise<SignInSubmitResponse> => {
    // Always try to save the post first
    const postResult = await actions.createAnonymousPost(title, description, attachments, canEditTags ? tags.map((tag) => tag.slug) : undefined)

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

  const handleEditorFocus = () => {
    // This function is called when the editor is focused
    // We don't need to do anything special here
  }

  return (
    <Modal.Window className="c-share-feedback" isOpen={isOpen} onClose={handleClose} size="fullscreen" center={false}>
      <Modal.Header>
        <div className="flex flex-items-center justify-end">
          <CloseIcon closeModal={handleClose} />
        </div>
      </Modal.Header>
      <Modal.Content>
        <div className="c-share-feedback__content mb-4">
          <h1 className="text-large pb-6">
            <Trans id="newpost.modal.title">Share your idea...</Trans>
          </h1>
          <div className="c-share-feedback-form">
            <Form error={error}>
              <div ref={editorRef} className="mb-4">
                <CommentEditor
                  field="description"
                  onChange={handleDescriptionChange}
                  onFocus={handleEditorFocus}
                  initialValue={description}
                  disabled={false}
                  placeholder={i18n._({
                    id: "newpost.modal.description.placeholder",
                    message: "Tell us about it. Explain it fully, don't hold back, the more information the better.",
                  })}
                  onImageUploaded={handleImageUploaded}
                />
              </div>
              <SimilarPosts title={title} tags={props.tags} />
              <Input
                field="title"
                inputRef={titleRef}
                maxLength={255}
                label={i18n._({ id: "newpost.modal.title.label", message: "Give your idea a title" })}
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
        {/* For unauthenticated users, always show the sign-in control */}
        {!fider.session.isAuthenticated ? (
          <div className="c-share-feedback__content">
            <div className="c-share-feedback-signin">
              <h2 className="text-title text-center mb-4">
                <Trans id="newpost.modal.submit">Submit your idea</Trans>
              </h2>
              <SignInControl
                onSubmit={onSubmitFeedback}
                onEmailSent={onEmailSent}
                signInButtonText={i18n._({ id: "signin.message.email", message: "Continue with Email" })}
                useEmail={true}
                redirectTo={fider.settings.baseURL}
              />
            </div>
          </div>
        ) : (
          /* For authenticated users, only show the submit button container when title is long enough */
          title.length > 9 && (
            <div className="c-share-feedback__content animate-fade-in">
              <div className="c-share-feedback-signin">
                <div className="flex justify-center">
                  <Button variant="primary" onClick={finaliseFeedback}>
                    <Trans id="newpost.modal.submit">Submit your idea</Trans>
                  </Button>
                </div>
              </div>
            </div>
          )
        )}
        {!fider.session.isAuthenticated ? <LegalFooter /> : null}
      </Modal.Content>
    </Modal.Window>
  )
}
