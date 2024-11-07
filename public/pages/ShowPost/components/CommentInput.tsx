import React, { useState, useRef } from "react"

import { Post, ImageUpload } from "@fider/models"
import { Avatar, UserName, Button, TextArea, Form, MultiImageUploader } from "@fider/components"
import { SignInModal } from "@fider/components"
import { User, UserRole, UserStatus } from "@fider/models/identity"

import { cache, actions, Failure, Fider } from "@fider/services"
import { useFider } from "@fider/hooks"
import { HStack } from "@fider/components/layout"
import { t, Trans } from "@lingui/macro"
import MentionSelector from "./MentionSelector"

interface CommentInputProps {
  post: Post
}

const CACHE_TITLE_KEY = "CommentInput-Comment-"

const users: User[] = [
  { id: 1, name: "Matt", role: UserRole.Administrator, avatarURL: "", status: UserStatus.Active },
  { id: 2, name: "Matt", role: UserRole.Administrator, avatarURL: "", status: UserStatus.Active },
  { id: 3, name: "Matt", role: UserRole.Administrator, avatarURL: "", status: UserStatus.Active },
]

export const CommentInput = (props: CommentInputProps) => {
  const getCacheKey = () => `${CACHE_TITLE_KEY}${props.post.id}`

  const fider = useFider()
  const inputRef = useRef<HTMLTextAreaElement>()
  const [content, setContent] = useState((fider.session.isAuthenticated && cache.session.get(getCacheKey())) || "")
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)

  const [dropdownVisible, setDropdownVisible] = useState<boolean>(false)
  const [filteredUsers, setFilteredUsers] = useState<User[]>([])
  const [mentionIndex, setMentionIndex] = useState<number | null>(null)
  const [cursorPosition, setCursorPosition] = useState<{ top: number; left: number } | null>(null)

  const commentChanged = (newContent: string, selectionStart?: number) => {
    cache.session.set(getCacheKey(), newContent)
    setContent(newContent)

    // Get the cursor position
    const textBeforeCursor = newContent.slice(0, selectionStart)

    // Check for "@" mention
    const lastAtIndex = textBeforeCursor.lastIndexOf("@")
    if (lastAtIndex >= 0) {
      const mentionQuery = textBeforeCursor.slice(lastAtIndex + 1)
      const matchedUsers = users.filter((user) => user.name.toLowerCase().startsWith(mentionQuery.toLowerCase()))

      setFilteredUsers(matchedUsers)
      setDropdownVisible(matchedUsers.length > 0)
      setMentionIndex(lastAtIndex)

      // Calculate position of dropdown
      if (inputRef.current) {
        const textareaRect = inputRef.current.getBoundingClientRect()
        const top = textareaRect.top + window.scrollY + inputRef.current.scrollTop - 160
        const left = textareaRect.left + window.scrollX - 50 // Adjust as needed
        setCursorPosition({ top, left })
      }
    } else {
      setDropdownVisible(false)
    }
  }

  console.log(dropdownVisible, filteredUsers, mentionIndex, cursorPosition)

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
      <HStack spacing={2} center={false} className="c-comment-input">
        {Fider.session.isAuthenticated && <Avatar user={Fider.session.user} />}
        <div className="flex-grow bg-gray-50 rounded-md p-2">
          <Form error={error}>
            {Fider.session.isAuthenticated && (
              <div className="mb-1">
                <UserName user={Fider.session.user} />
              </div>
            )}
            <TextArea
              placeholder={t({ id: "showpost.commentinput.placeholder", message: "Leave a comment" })}
              field="content"
              disabled={fider.isReadOnly}
              value={content}
              minRows={1}
              onChange={commentChanged}
              onFocus={handleOnFocus}
              onKeyDown={(e) => {
                if (e.key === "Escape") {
                  setDropdownVisible(false)
                }
              }}
              inputRef={inputRef}
            />

            {dropdownVisible && cursorPosition && <MentionSelector cursorPosition={cursorPosition} names={filteredUsers.map((user) => user.name)} />}

            {content && (
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
