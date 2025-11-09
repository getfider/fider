import "./ShareFeedback.scss"

import React, { useEffect, useRef, useState } from "react"
import { SignInControl } from "@fider/components/common/SignInControl"
import { Modal, CloseIcon, Form, Button, Input, LegalFooter } from "@fider/components/common"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { actions, Failure, querystring, classSet } from "@fider/services"
import { plainText } from "@fider/services/markdown"
import { i18n } from "@lingui/core"
import { Tag } from "@fider/models"
import { SimilarPosts } from "../components/SimilarPosts"
import { TagsSelect } from "@fider/components/common/TagsSelect"
import CommentEditor from "@fider/components/common/form/CommentEditor"
import {
  CACHE_KEYS,
  clearCache,
  getCachedDescription,
  getCachedTags,
  getCachedTitle,
  setCachedDescription,
  setCachedTags,
  setCachedTitle,
  setPostCreated,
  setPostPending,
} from "./PostCache"
import { useAttachments } from "@fider/hooks/useAttachments"

interface ShareFeedbackProps {
  isOpen: boolean
  placeholder: string
  onClose: () => void
  tags: Tag[]
}

export const ShareFeedback: React.FC<ShareFeedbackProps> = (props) => {
  const fider = useFider()
  const { isOpen, onClose } = props

  const getTagsCachedValue = (): Tag[] => {
    if (!canEditTags) {
      return []
    }

    const cacheValue = getCachedTags()
    const urlValue = querystring.get("tags")
    const combined = [...cacheValue, ...urlValue.split(",")]
    const tagsAsStrings = Array.from(new Set(combined.map((s) => s.trim()).filter((s) => s.length > 0)))

    return props.tags.filter((tag) => tagsAsStrings.includes(tag.slug))
  }

  const getTitleManuallyEditedValue = (): boolean => {
    // If the cached title deviates from the description, it means the user manually edited it
    return getCachedTitle() !== getCachedDescription()
  }

  const canEditTags = fider.settings.postWithTags && props.tags.length > 0
  const [title, setTitle] = useState(getCachedTitle())
  const [description, setDescription] = useState(getCachedDescription())
  const { attachments, handleImageUploaded, getImageSrc, clearAttachments } = useAttachments({
    cacheKey: CACHE_KEYS.ATTACHMENT,
    useLocalStorage: true,
    maxAttachments: 3,
  })
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
      // Find newline in the original markdown content for truncation
      let newlineIndex = Math.min(description.indexOf("\n"), 80)
      if (newlineIndex == -1) {
        newlineIndex = 80
      }

      // Get the truncated markdown content and convert to plain text
      const truncatedMarkdown = description.substring(0, newlineIndex)
      const autoTitle = plainText(truncatedMarkdown)

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
    setCachedTitle(value)
    // If this is a manual edit (not auto-generated from description),
    // mark the title as manually edited so we stop auto-populating
    // If the user clears the title, we still want to allow auto-population
    if (isManualEdit) {
      setTitleManuallyEdited(value !== "")
    }
  }

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault()
    }
  }

  const handleTagsChanged = (newTags: Tag[]) => {
    setCachedTags(newTags.map((tag) => tag.slug))
    setTags(newTags)
  }

  const handleDescriptionChange = (value: string) => {
    setCachedDescription(value)

    // If the description starts with an image attachment, we don't want to set it as the title
    if (value.startsWith("![](fider-image:attachments")) {
      return
    }

    setDescription(value)
  }

  const onSubmitFeedback = () => {
    setPostPending(true)
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
        clearCache()
        clearAttachments()
        setPostCreated()
        location.href = `/posts/${result.data.number}/${result.data.slug}`
      } else if (result.error) {
        setError(result.error)
      }
    }
  }

  const onCodeVerified = (result: { showProfileCompletion?: boolean; code?: string }): void => {
    if (result.showProfileCompletion && result.code) {
      // User needs to complete profile - redirect to profile completion page
      // The cached feedback will be preserved for after profile setup
      location.href = `/signin/complete?code=${encodeURIComponent(result.code)}`
    } else {
      // User is authenticated - finalize the feedback submission
      finaliseFeedback()
    }
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
                  disabled={fider.isReadOnly}
                  maxAttachments={3}
                  maxImageSizeKB={5 * 1024}
                  placeholder={i18n._({
                    id: "newpost.modal.description.placeholder",
                    message: "Tell us about it. Explain it fully, don't hold back, the more information the better.",
                  })}
                  onImageUploaded={handleImageUploaded}
                  onGetImageSrc={getImageSrc}
                />
              </div>
              <SimilarPosts title={title} tags={props.tags} />
              <Input
                field="title"
                inputRef={titleRef}
                maxLength={255}
                label={i18n._({ id: "newpost.modal.title.label", message: "Give your idea a title" })}
                value={title}
                disabled={fider.isReadOnly}
                onChange={handleTitleChange}
                onKeyDown={handleKeyDown}
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
                onCodeVerified={onCodeVerified}
                signInButtonText={i18n._({ id: "signin.message.email", message: "Continue with Email" })}
                useEmail={true}
                redirectTo={fider.settings.baseURL}
              />
            </div>
          </div>
        ) : (
          /* For authenticated users, only show the submit button container when title is long enough */
          title.replace(/\s+/g, " ").trim().length > 9 && (
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
