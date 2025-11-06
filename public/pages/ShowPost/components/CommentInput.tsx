import React, { useCallback, useState, useEffect } from "react"

import { Post } from "@fider/models"
import { Avatar, UserName, Button, Form } from "@fider/components"
import { SignInModal } from "@fider/components"

import { cache, actions, Failure, Fider } from "@fider/services"
import { HStack } from "@fider/components/layout"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

import { useFider } from "@fider/hooks"
import { useAttachments } from "@fider/hooks/useAttachments"
import CommentEditor from "@fider/components/common/form/CommentEditor"

interface CommentInputProps {
  post: Post
}

const CACHE_TITLE_KEY = "CommentInput-Comment-Title-"
const CACHE_ATTACHMENTS_KEY = "CommentInput-Comment-Attachments-"

export const CommentInput = (props: CommentInputProps) => {
  const getCacheKey = (cachePrefix: string) => `${cachePrefix}${props.post.id}`

  const getContentFromCache = () => {
    return cache.session.get(getCacheKey(CACHE_TITLE_KEY))
  }

  const fider = useFider()
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [error, setError] = useState<Failure | undefined>(undefined)
  const [isClient, setIsClient] = useState(false)

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

    const result = await actions.createComment(props.post.number, content || "", attachments)
    if (result.ok) {
      clearAttachments()
      cache.session.remove(getCacheKey(CACHE_TITLE_KEY))
      if (fider.session.tenant.isModerationEnabled && !fider.session.user.isCollaborator) {
        cache.session.set("COMMENT_CREATED_MODERATION", "true")
      }
      location.reload()
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
  }, [])

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <HStack spacing={2} className="c-comment-input" align="start">
        {Fider.session.isAuthenticated && <Avatar user={Fider.session.user} />}
        <div className="flex-grow bg-gray-100 rounded-md p-2">
          <Form error={error}>
            {Fider.session.isAuthenticated && (
              <div className="mb-1">
                <UserName user={Fider.session.user} />
              </div>
            )}

            {isClient ? (
              <>
                <CommentEditor
                  field="content"
                  disabled={!Fider.session.isAuthenticated}
                  onChange={commentChanged}
                  onFocus={editorFocused}
                  initialValue={getContentFromCache()}
                  placeholder={i18n._({ id: "showpost.commentinput.placeholder", message: "Leave a comment" })}
                  maxAttachments={2}
                  maxImageSizeKB={5 * 1024}
                  onGetImageSrc={getImageSrc}
                  onImageUploaded={handleImageUploaded}
                />

                {fider.session.isModerationRequired && Fider.session.isAuthenticated && !Fider.session.user.isCollaborator && (
                  <div className="mt-2 text-muted text-sm p-2 bg-gray-50 rounded border-l-4 border-yellow-500">
                    <Trans id="comment.moderation.notice">Your comment will be reviewed by an administrator before being visible to other users.</Trans>
                  </div>
                )}

                {hasContent && (
                  <>
                    <Button variant="primary" onClick={submit} className="mt-2">
                      <Trans id="action.submit">Submit</Trans>
                    </Button>
                  </>
                )}
              </>
            ) : (
              <div className="comment-input-placeholder p-2">{i18n._({ id: "showpost.commentinput.placeholder", message: "Leave a comment" })}</div>
            )}
          </Form>
        </div>
      </HStack>
    </>
  )
}
