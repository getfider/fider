import React, { useState, useEffect, useRef } from "react"
import { Button, ButtonClickEvent, Input, Form, TextArea, MultiImageUploader } from "@fider/components"
import { SignInModal } from "@fider/components"
import { cache, actions, Failure, navigator, querystring } from "@fider/services"
import { ImageUpload, Tag } from "@fider/models"
import { useFider } from "@fider/hooks"
import { t, Trans } from "@lingui/macro"
import { TagsFilter } from "./TagsFilter"

interface PostInputProps {
  placeholder: string
  onTitleChanged: (title: string) => void
  tags: Tag[]
}

const CACHE_TITLE_KEY = "PostInput-Title"
const CACHE_DESCRIPTION_KEY = "PostInput-Description"

export const PostInput = (props: PostInputProps) => {
  const getCachedValue = (key: string): string => {
    if (fider.session.isAuthenticated) {
      return cache.session.get(key) || ""
    }
    return ""
  }

  const fider = useFider()
  const titleRef = useRef<HTMLInputElement>()
  const [title, setTitle] = useState(getCachedValue(CACHE_TITLE_KEY))
  const [description, setDescription] = useState(getCachedValue(CACHE_DESCRIPTION_KEY))
  const [tags, setTags] = useState<string[]>(querystring.getArray("tags"))
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)

  useEffect(() => {
    props.onTitleChanged(title)
  }, [title])

  const handleTitleFocus = () => {
    if (!fider.session.isAuthenticated && titleRef.current) {
      titleRef.current.blur()
      setIsSignInModalOpen(true)
    }
  }

  const handleTitleChange = (value: string) => {
    cache.session.set(CACHE_TITLE_KEY, value)
    setTitle(value)
    props.onTitleChanged(value)
  }

  const hideModal = () => setIsSignInModalOpen(false)
  const clearError = () => setError(undefined)

  const handleDescriptionChange = (value: string) => {
    cache.session.set(CACHE_DESCRIPTION_KEY, value)
    setDescription(value)
  }

  const submit = async (event: ButtonClickEvent) => {
    if (title) {
      const result = await actions.createPost(title, description, attachments, tags)
      if (result.ok) {
        clearError()
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY)
        location.href = `/posts/${result.data.number}/${result.data.slug}`
        event.preventEnable()
      } else if (result.error) {
        setError(result.error)
      }
    }
  }

  const handleTagsChanged = (newTags: string[]) => {
    setTags(newTags)

    navigator.replaceState(
      querystring.stringify({
        view: querystring.get("view"),
        query: querystring.get("query"),
        tags: newTags,
        limit: querystring.getNumber("limit"),
      })
    )
  }

  const details = () => (
    <>
      <TextArea
        field="description"
        onChange={handleDescriptionChange}
        value={description}
        minRows={5}
        placeholder={t({ id: "home.postinput.description.placeholder", message: "Describe your suggestion (optional)" })}
      />
      <TagsFilter tags={props.tags} selectionChanged={handleTagsChanged} selected={tags} />
      <MultiImageUploader field="attachments" maxUploads={3} onChange={setAttachments} />
      <Button type="submit" variant="primary" onClick={submit}>
        <Trans id="action.submit">Submit</Trans>
      </Button>
    </>
  )

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <Form error={error}>
        <Input
          field="title"
          disabled={fider.isReadOnly}
          noTabFocus={!fider.session.isAuthenticated}
          inputRef={titleRef}
          onFocus={handleTitleFocus}
          maxLength={100}
          value={title}
          onChange={handleTitleChange}
          placeholder={props.placeholder}
        />
        {title && details()}
      </Form>
    </>
  )
}
