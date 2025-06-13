import React, { useState, useEffect, useRef } from "react"
import { Button, ButtonClickEvent, Input, Form, TextArea, MultiImageUploader } from "@fider/components"
import { SignInModal } from "@fider/components"
import { cache, actions, classSet, Failure, querystring } from "@fider/services"
import { ImageUpload, Tag } from "@fider/models"
import { useFider } from "@fider/hooks"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
// import { CACHE_ATTACHMENT_KEY } from "./ShareFeedback"
import { TagsSelect } from "@fider/components/common/TagsSelect"

interface PostInputProps {
  placeholder: string
  onTitleChanged: (title: string) => void
  tags: Tag[]
}

const CACHE_TITLE_KEY = "PostInput-Title"
const CACHE_DESCRIPTION_KEY = "PostInput-Description"
const CACHE_TAGS_KEY = "PostInput-Tags"

export const PostInput = (props: PostInputProps) => {
  const getCachedValue = (key: string): string => {
    if (fider.session.isAuthenticated) {
      return cache.session.get(key) || ""
    }
    return ""
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

  const fider = useFider()
  const canEditTags = fider.session.isAuthenticated && fider.settings.postWithTags && props.tags.length > 0
  const titleRef = useRef<HTMLInputElement>()
  const [title, setTitle] = useState(getCachedValue(CACHE_TITLE_KEY))
  const [description, setDescription] = useState(getCachedValue(CACHE_DESCRIPTION_KEY))
  const [tags, setTags] = useState(getTagsCachedValue())
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
      const result = await actions.createPost(
        title,
        description,
        attachments,
        tags.map((tag) => tag.slug)
      )
      if (result.ok) {
        clearError()
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY, CACHE_TAGS_KEY)
        location.href = `/posts/${result.data.number}/${result.data.slug}`
        event.preventEnable()
      } else if (result.error) {
        setError(result.error)
      }
    }
  }

  const handleTagsChanged = (newTags: Tag[]) => {
    cache.session.set(CACHE_TAGS_KEY, newTags.map((tag) => tag.slug).join(","))
    setTags(newTags)
  }

  const details = () => (
    <>
      <TextArea
        field="description"
        onChange={handleDescriptionChange}
        value={description}
        minRows={5}
        placeholder={i18n._("home.postinput.description.placeholder", { message: "Describe your suggestion (optional)" })}
      />
      {canEditTags && (
        <div className={classSet({ "c-form-field": true })}>
          <TagsSelect tags={props.tags} selectionChanged={handleTagsChanged} selected={tags} canEdit={true} />
        </div>
      )}
      <MultiImageUploader field="attachments" maxUploads={3} onChange={setAttachments} noPadding />
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
