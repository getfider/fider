import React, { useCallback, useState, useEffect } from "react"

import { Post } from "@fider/models"
import { Avatar, Button, Form } from "@fider/components"
import { SignInModal } from "@fider/components"

import { cache, actions, Failure, Fider, notify } from "@fider/services"
import { HStack } from "@fider/components/layout"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

import { useFider } from "@fider/hooks"
import { useAttachments } from "@fider/hooks/useAttachments"
import CommentEditor from "@fider/components/common/form/CommentEditor"

interface CommentInputProps {
  post: Post
  onCommentAdded: () => void
  // Set when replying within a thread
  parentId?: number
  // Used when there's no draft, e.g. an @mention when replying to a reply
  initialValue?: string
  onCancel?: () => void
}

const CACHE_TITLE_KEY = "CommentInput-Comment-Title-"
const CACHE_ATTACHMENTS_KEY = "CommentInput-Comment-Attachments-"

export const CommentInput = (props: CommentInputProps) => {
  const getCacheKey = (cachePrefix: string) => `${cachePrefix}${props.post.id}${props.parentId ? `-${props.parentId}` : ""}`

  const getContentFromCache = () => {
    return cache.session.get(getCacheKey(CACHE_TITLE_KEY))
  }

  const COMMENT_MAX_LENGTH = 4000

  const fider = useFider()
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [error, setError] = useState<Failure | undefined>(undefined)
  const [isClient, setIsClient] = useState(false)
  // Bumped after posting to remount the editor empty
  const [editorKey, setEditorKey] = useState(0)
  const [initialContent, setInitialContent] = useState(() => {
    const draft = getContentFromCache()
    if (!draft && props.initialValue) {
      cache.session.set(getCacheKey(CACHE_TITLE_KEY), props.initialValue)
      return props.initialValue
    }
    return draft
  })
  const [contentLength, setContentLength] = useState((initialContent ?? "").length)

  // Use the attachments hook
  const { attachments, handleImageUploaded, getImageSrc, clearAttachments } = useAttachments({
    cacheKey: getCacheKey(CACHE_ATTACHMENTS_KEY),
    maxAttachments: 2,
  })

  // Check if we're running on the client after component mounts
  useEffect(() => {
    setIsClient(true)
  }, [])

  const hideModal = () => setIsSignInModalOpen(false)
  const clearError = () => setError(undefined)

  const submit = async () => {
    clearError()

    const content = getContentFromCache()

    const result = await actions.createComment(props.post.number, content || "", attachments, props.parentId)
    if (result.ok) {
      clearAttachments()
      cache.session.remove(getCacheKey(CACHE_TITLE_KEY))
      setContentLength(0)
      setInitialContent("")
      setEditorKey((k) => k + 1)
      if (fider.session.isModerationRequiredForNewPost) {
        notify.success(i18n._({ id: "showpost.moderation.commentsuccess", message: "Your comment is awaiting moderation 📝" }))
      }
      props.onCommentAdded()
    } else {
      setError(result.error)
    }
  }

  const editorFocused = () => {
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true)
    }
  }

  const hasContent = true

  const commentChanged = useCallback((value: string): void => {
    cache.session.set(getCacheKey(CACHE_TITLE_KEY), value)
    setContentLength(value.length)
  }, [])

  const isOverLimit = contentLength > COMMENT_MAX_LENGTH

  const placeholder = props.parentId
    ? i18n._({ id: "showpost.commentinput.replyplaceholder", message: "Write a reply" })
    : i18n._({ id: "showpost.commentinput.placeholder", message: "Leave a comment" })

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <HStack spacing={4} className="c-comment-input" align="start">
        {Fider.session.isAuthenticated && <Avatar user={Fider.session.user} size="large" />}
        <div className="c-comment-input-card">
          <Form error={error}>
            {isClient ? (
              <>
                <CommentEditor
                  key={editorKey}
                  field="content"
                  disabled={!Fider.session.isAuthenticated}
                  onChange={commentChanged}
                  onFocus={editorFocused}
                  initialValue={initialContent}
                  placeholder={placeholder}
                  maxAttachments={2}
                  maxImageSizeKB={5 * 1024}
                  maxLength={COMMENT_MAX_LENGTH}
                  onGetImageSrc={getImageSrc}
                  onImageUploaded={handleImageUploaded}
                />

                {hasContent && (
                  <>
                    <Button disabled={!fider.session.isAuthenticated || isOverLimit} variant="primary" onClick={submit} className="mt-4">
                      <Trans id="action.postcomment">Post</Trans>
                    </Button>
                    {props.onCancel && (
                      <Button variant="tertiary" onClick={props.onCancel} className="mt-4">
                        <Trans id="action.cancel">Cancel</Trans>
                      </Button>
                    )}
                  </>
                )}
              </>
            ) : (
              <div className="comment-input-placeholder p-2">{placeholder}</div>
            )}
          </Form>
        </div>
      </HStack>
    </>
  )
}
