import React, { useCallback, useState } from "react"

import { Post, ImageUpload } from "@fider/models"
import { Avatar, UserName, Button, Form, MultiImageUploader } from "@fider/components"
import { SignInModal } from "@fider/components"

import { cache, actions, Failure, Fider } from "@fider/services"
import { HStack } from "@fider/components/layout"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

// import { CommentEditor } from "@fider/components"
import { useFider } from "@fider/hooks"
// import Tiptap from "@fider/components/common/form/CommentEditor2"
import CommentEditor from "@fider/components/common/form/CommentEditor"

interface CommentInputProps {
  post: Post
}

const CACHE_TITLE_KEY = "CommentInput-Comment-"

export const CommentInput = (props: CommentInputProps) => {
  const getCacheKey = () => `${CACHE_TITLE_KEY}${props.post.id}`

  const getContentFromCache = () => {
    return cache.session.get(getCacheKey())
  }

  const fider = useFider()
  // const inputRef = useRef<HTMLTextAreaElement>()
  // const [content, setContent] = useState<string>((fider.session.isAuthenticated && getContentFromCache()) || "")
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)

  const hideModal = () => setIsSignInModalOpen(false)
  const clearError = () => setError(undefined)

  const submit = async () => {
    clearError()

    // Since the comment is being cached, we can save the content that's in the cache
    const content = getContentFromCache()

    const result = await actions.createComment(props.post.number, content || "", attachments)
    if (result.ok) {
      cache.session.remove(getCacheKey())
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
    cache.session.set(getCacheKey(), value)
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
            <CommentEditor
              disabled={!Fider.session.isAuthenticated}
              onChange={commentChanged}
              onFocus={editorFocused}
              initialValue={getContentFromCache()}
              placeholder={i18n._("showpost.commentinput.placeholder", { message: "Leave a comment" })}
            />

            {hasContent && (
              <>
                <MultiImageUploader field="attachments" maxUploads={2} onChange={setAttachments} />
                <Button variant="primary" onClick={submit}>
                  <Trans id="action.submit">Submit</Trans>
                </Button>
              </>
            )}
          </Form>
        </div>
      </HStack>
    </>
  )
}
