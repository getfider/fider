import React, { useState, useEffect, useRef } from "react"
import { Input } from "./Input"
import { Form } from "./Form"
import { Failure } from "@fider/services"
import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"

interface LinkInsertModalProps {
  isOpen: boolean
  onClose: () => void
  onInsertLink: (text: string, url: string) => void
  selectedText?: string
}

const LinkInsertModal = ({ isOpen, onClose, onInsertLink, selectedText = "" }: LinkInsertModalProps) => {
  const [text, setText] = useState("")
  const [url, setUrl] = useState("")
  const [error, setError] = useState<Failure | undefined>(undefined)
  const textInputRef = useRef<HTMLInputElement>(null)

  // Create an error item with the given field, message ID, and message
  const createError = (field: string, messageId: string, message: string) => ({
    field,
    message: i18n._({ id: messageId, message }),
  })

  const handleSubmit = () => {
    // Clear previous errors
    setError(undefined)

    // Custom validation
    const errorItems: { field?: string; message: string }[] = []

    if (!text.trim()) {
      errorItems.push(createError("text", "linkmodal.text.required", "Text is required"))
    }

    if (!url.trim()) {
      errorItems.push(createError("url", "linkmodal.url.required", "URL is required"))
    } else {
      // Basic URL validation
      try {
        new URL(url.startsWith("http://") || url.startsWith("https://") ? url : `https://${url}`)
      } catch {
        errorItems.push(createError("url", "linkmodal.url.invalid", "Please enter a valid URL"))
      }
    }

    if (errorItems.length > 0) {
      setError({ errors: errorItems })
      return
    }

    // Form is valid, handle submission
    const urlWithProtocol = url.startsWith("http://") || url.startsWith("https://") ? url : `https://${url}`
    onInsertLink(text.trim(), urlWithProtocol)
    setText("")
    setUrl("")
    onClose()
  }

  const handleClose = () => {
    setText("")
    setUrl("")
    setError(undefined)
    onClose()
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Escape") {
      handleClose()
    }
  }

  // Focus management and prefill text
  useEffect(() => {
    if (isOpen) {
      // Prefill with the highlighted text when the user clicked the link button
      setText(selectedText)
      if (textInputRef.current) {
        textInputRef.current.focus()
      }
    }
  }, [isOpen, selectedText])

  if (!isOpen) return null

  return (
    <div className="c-modal-dimmer" role="dialog" aria-modal="true" aria-labelledby="link-modal-title" onKeyDown={handleKeyDown}>
      <div className="c-modal-scroller">
        <div className="c-modal-window c-modal-window--small">
          <div className="c-modal-header">
            <h3 id="link-modal-title">
              <Trans id="linkmodal.title">Insert Link</Trans>
            </h3>
          </div>
          <div className="c-modal-content">
            <Form error={error}>
              <Input
                field="text"
                label={i18n._({ id: "linkmodal.text.label", message: "Text to display" })}
                value={text}
                onChange={setText}
                placeholder={i18n._({ id: "linkmodal.text.placeholder", message: "Enter link text" })}
                inputRef={textInputRef}
              />
              <Input
                field="url"
                label={i18n._({ id: "linkmodal.url.label", message: "URL" })}
                value={url}
                onChange={setUrl}
                placeholder={i18n._({ id: "linkmodal.url.placeholder", message: "https://example.com" })}
              />
            </Form>
          </div>
          <div className="c-modal-footer c-modal-footer--right">
            <div style={{ display: "flex", gap: "8px" }}>
              <button type="button" onClick={handleClose} className="c-button c-button--secondary">
                <Trans id="action.cancel">Cancel</Trans>
              </button>
              <button type="button" onClick={handleSubmit} className="c-button c-button--primary">
                <Trans id="linkmodal.insert">Insert Link</Trans>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LinkInsertModal
