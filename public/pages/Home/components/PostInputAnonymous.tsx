import React, { useState, useEffect, useRef } from "react"
import { Button, ButtonClickEvent, Input, Form, TextArea, MultiImageUploader, Icon } from "@fider/components"
import { cache, actions, Failure } from "@fider/services"
import { ImageUpload } from "@fider/models"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"
import IconAttach from "@fider/assets/images/heroicons-paperclip.svg"

interface PostInputProps {
  placeholder: string
}

const CACHE_TITLE_KEY = "PostInput-Title"
const CACHE_DESCRIPTION_KEY = "PostInput-Description"

export const PostInputAnonymous = (props: PostInputProps) => {
  const fider = useFider()
  const getCachedValue = (key: string): string => {
    if (fider.session.isAuthenticated) {
      return cache.session.get(key) || ""
    }
    return ""
  }

  const titleRef = useRef<HTMLInputElement>()
  const [title, setTitle] = useState(getCachedValue(CACHE_TITLE_KEY))
  const [description, setDescription] = useState(getCachedValue(CACHE_DESCRIPTION_KEY))
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)
  const [titleManuallyEdited, setTitleManuallyEdited] = useState(false)

  // Auto-populate title with first 80 characters of description if title hasn't been manually edited
  useEffect(() => {
    if (!titleManuallyEdited) {
      let newlineIndex = Math.min(description.indexOf("\n"), 80)
      if (newlineIndex == -1) {
        newlineIndex = 80
      }
      const autoTitle = description.substring(0, newlineIndex)
      handleTitleChange(autoTitle, false)
    }
  }, [description, titleManuallyEdited])
  const handleTitleChange = (value: string, isManualEdit = true) => {
    cache.session.set(CACHE_TITLE_KEY, value)
    setTitle(value)

    // If this is a manual edit (not auto-generated from description),
    // mark the title as manually edited so we stop auto-populating
    if (isManualEdit) {
      setTitleManuallyEdited(true)
    }
  }

  const clearError = () => setError(undefined)

  const handleDescriptionChange = (value: string) => {
    cache.session.set(CACHE_DESCRIPTION_KEY, value)
    setDescription(value)
  }

  const submit = async (event: ButtonClickEvent) => {
    if (title) {
      const result = await actions.createAnonymousPost(title, description, attachments)
      if (result.ok) {
        clearError()
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY)
        location.href = `/draft/${result.data.code}`
        event.preventEnable()
      } else if (result.error) {
        setError(result.error)
      }
    }
  }

  return (
    <>
      <Form error={error}>
        <TextArea
          label={props.placeholder}
          field="description"
          onChange={handleDescriptionChange}
          value={description}
          minRows={5}
          placeholder={i18n._({
            id: "newpost.modal.description.placeholder",
            message: "Tell us about your idea. Explain it fully, don't hold back, the more information the better.",
          })}
        />
        <MultiImageUploader
          field="attachments"
          maxUploads={3}
          onChange={setAttachments}
          addImageButton={
            <a className="flex items-center clickable">
              <Icon sprite={IconAttach} height="18" width="18" />
              <span className="ml-1">
                <Trans id="newpost.modal.addimage">Add Images</Trans>
              </span>
            </a>
          }
        />
        <Input
          field="title"
          inputRef={titleRef}
          maxLength={255}
          label={i18n._({ id: "newpost.modal.title.label", message: "Give your suggestion a title" })}
          value={title}
          onChange={handleTitleChange}
          placeholder={i18n._({ id: "newpost.modal.title.placeholder", message: "Something short and snappy, sum it up in a few words" })}
        />

        <Button type="submit" variant="primary" onClick={submit}>
          <Trans id="action.submit">Submit</Trans>
        </Button>
      </Form>
    </>
  )
}
