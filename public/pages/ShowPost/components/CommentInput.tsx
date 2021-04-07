import React, { useState, useRef } from "react"

import { Post, ImageUpload } from "@fider/models"
import { Avatar, UserName, Button, TextArea, Form, MultiImageUploader } from "@fider/components/common"
import { SignInModal } from "@fider/components"

import { cache, actions, Failure, Fider } from "@fider/services"
import { useFider } from "@fider/hooks"

interface CommentInputProps {
  post: Post
}

const CACHE_TITLE_KEY = "CommentInput-Comment-"

export const CommentInput = (props: CommentInputProps) => {
  const getCacheKey = () => `${CACHE_TITLE_KEY}${props.post.id}`

  const fider = useFider()
  const inputRef = useRef<HTMLTextAreaElement>()
  const [content, setContent] = useState((fider.session.isAuthenticated && cache.session.get(getCacheKey())) || "")
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)

  const commentChanged = (newContent: string) => {
    cache.session.set(getCacheKey(), newContent)
    setContent(newContent)
  }

  const hideModal = () => setIsSignInModalOpen(false)
  const clearError = () => setError(undefined)

  const submit = async () => {
    clearError()

    const result = await actions.createComment(props.post.number, content, attachments)
    if (result.ok) {
      cache.session.remove(getCacheKey())
      location.reload()
    } else {
      setError(result.error)
    }
  }

  const handleOnFocus = () => {
    if (!fider.session.isAuthenticated && inputRef.current) {
      inputRef.current.blur()
      setIsSignInModalOpen(true)
    }
  }

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <div className={`c-comment-input ${Fider.session.isAuthenticated && "m-authenticated"}`}>
        {Fider.session.isAuthenticated && <Avatar size="large" user={Fider.session.user} />}
        <Form error={error}>
          {Fider.session.isAuthenticated && <UserName user={Fider.session.user} />}
          <TextArea
            placeholder="Leave a comment"
            field="content"
            value={content}
            minRows={1}
            onChange={commentChanged}
            onFocus={handleOnFocus}
            inputRef={inputRef}
          />
          {content && (
            <>
              <MultiImageUploader field="attachments" maxUploads={2} previewMaxWidth={100} onChange={setAttachments} />
              <Button color="positive" onClick={submit}>
                Submit
              </Button>
            </>
          )}
        </Form>
      </div>
    </>
  )
}
