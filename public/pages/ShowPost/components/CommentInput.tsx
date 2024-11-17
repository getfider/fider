import React, { useState } from "react"

import { Post, ImageUpload } from "@fider/models"
import { Avatar, UserName, Button, Form, MultiImageUploader } from "@fider/components"
import { SignInModal } from "@fider/components"

import { cache, actions, Failure, Fider } from "@fider/services"
import { useFider } from "@fider/hooks"
import { HStack } from "@fider/components/layout"
import { t, Trans } from "@lingui/macro"

import { SlateEditor, Serialize } from "@fider/components"
import { Descendant } from "slate"
import { Paragraph } from "@fider/components/common/form/SlateEditor"

interface CommentInputProps {
  post: Post
}

const CACHE_TITLE_KEY = "CommentInput-Comment-"

export const CommentInput = (props: CommentInputProps) => {
  const getCacheKey = () => `${CACHE_TITLE_KEY}${props.post.id}`

  const getContentFromCache = () => {
    const cacheVal = cache.session.get(getCacheKey())
    if (cacheVal) {
      return JSON.parse(cacheVal)
    }
  }

  const fider = useFider()
  // const inputRef = useRef<HTMLTextAreaElement>()
  const [content, setContent] = useState<Descendant[]>((fider.session.isAuthenticated && getContentFromCache()) || null)
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)

  const hideModal = () => setIsSignInModalOpen(false)
  const clearError = () => setError(undefined)

  const submit = async () => {
    clearError()

    const editorText = Serialize(content)
    const result = await actions.createComment(props.post.number, editorText, attachments)
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

  const commentChanged = (newContent: Descendant[]) => {
    setContent(newContent)
    cache.session.set(getCacheKey(), JSON.stringify(newContent))
  }

  const hasContent = content?.length > 0 && (content[0] as Paragraph).children[0].text !== ""

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <HStack spacing={2} center={false} className="c-comment-input">
        {Fider.session.isAuthenticated && <Avatar user={Fider.session.user} />}
        <div className="flex-grow bg-gray-50 rounded-md p-2">
          <Form error={error}>
            {Fider.session.isAuthenticated && (
              <div className="mb-1">
                <UserName user={Fider.session.user} />
              </div>
            )}
            <SlateEditor
              initialValue={content}
              onChange={commentChanged}
              onFocus={editorFocused}
              disabled={!Fider.session.isAuthenticated}
              placeholder={t({ id: "showpost.commentinput.placeholder", message: "Leave a comment" })}
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
