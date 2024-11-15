import React, { useState } from "react"

import { Post, ImageUpload } from "@fider/models"
import { Avatar, UserName, Button, Form, MultiImageUploader } from "@fider/components"
import { SignInModal } from "@fider/components"

import { cache, actions, Failure, Fider } from "@fider/services"
import { useFider } from "@fider/hooks"
import { HStack } from "@fider/components/layout"
import { t, Trans } from "@lingui/macro"

import { createEditor } from "slate"
// Import the Slate components and React plugin.
import { Slate, Editable, withReact } from "slate-react"
//
// TypeScript users only add this code
import { BaseEditor, Descendant, Node } from "slate"
import { ReactEditor } from "slate-react"

type paragraph = { type: "paragraph"; children: text[] }
type text = { text: string }

declare module "slate" {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor
    Element: paragraph
    Text: text
  }
}

interface CommentInputProps {
  post: Post
}

const CACHE_TITLE_KEY = "CommentInput-Comment-"

const initialValue: Descendant[] = [
  {
    type: "paragraph",
    children: [{ text: "" }],
  },
]

export const CommentInput = (props: CommentInputProps) => {
  const getCacheKey = () => `${CACHE_TITLE_KEY}${props.post.id}`

  const [editor] = useState(() => withReact(createEditor()))

  const fider = useFider()
  // const inputRef = useRef<HTMLTextAreaElement>()
  // const [content, setContent] = useState((fider.session.isAuthenticated && cache.session.get(getCacheKey())) || "")
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)
  const [hasContent, setHasContent] = useState(false)

  const hideModal = () => setIsSignInModalOpen(false)
  const clearError = () => setError(undefined)

  const submit = async () => {
    clearError()

    const editorText = serialize(editor.children)
    console.log("editorText", editorText)
    const result = await actions.createComment(props.post.number, editorText, attachments)
    if (result.ok) {
      cache.session.remove(getCacheKey())
      location.reload()
    } else {
      setError(result.error)
    }
  }

  const serialize = (nodes: Descendant[]): string => {
    return nodes.map((n) => Node.string(n)).join("\n")
  }

  const handleOnFocus = () => {
    console.log("focus")
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true)
    }
  }

  const commentChanged = (newContent: Descendant[]) => {
    console.log("newContent", newContent, newContent.toString())
    const contentExists = newContent.length > 0 && (newContent[0] as paragraph).children[0].text !== ""
    setHasContent(contentExists)
  }

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
            <Slate editor={editor} initialValue={initialValue} onChange={commentChanged}>
              <Editable onFocus={handleOnFocus} placeholder={t({ id: "showpost.commentinput.placeholder", message: "Leave a comment" })} />
            </Slate>
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
