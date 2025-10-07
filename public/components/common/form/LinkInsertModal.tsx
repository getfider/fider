import React, { useState, useEffect, useRef } from "react"

interface LinkInsertModalProps {
  isOpen: boolean
  onClose: () => void
  onInsertLink: (text: string, url: string) => void
  selectedText?: string
}

const LinkInsertModal = ({ isOpen, onClose, onInsertLink, selectedText = "" }: LinkInsertModalProps) => {
  const [text, setText] = useState("")
  const [url, setUrl] = useState("")
  const [errors, setErrors] = useState<{ text?: string; url?: string }>({})
  const textInputRef = useRef<HTMLInputElement>(null)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    // Clear previous errors
    setErrors({})

    // Custom validation
    const newErrors: { text?: string; url?: string } = {}

    if (!text.trim()) {
      newErrors.text = "Text is required"
    }

    if (!url.trim()) {
      newErrors.url = "URL is required"
    } else {
      // Basic URL validation
      try {
        new URL(url.startsWith("http://") || url.startsWith("https://") ? url : `https://${url}`)
      } catch {
        newErrors.url = "Please enter a valid URL"
      }
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors)
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
    setErrors({})
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
            <h3 id="link-modal-title">Insert Link</h3>
          </div>
          <div className="c-modal-content">
            <form onSubmit={handleSubmit}>
              <div className="c-form-field">
                <label htmlFor="link-text">
                  Text to display <span style={{ color: "red" }}>*</span>
                </label>
                <input
                  ref={textInputRef}
                  id="link-text"
                  type="text"
                  value={text}
                  onChange={(e) => {
                    setText(e.target.value)
                    if (errors.text) {
                      setErrors((prev) => ({ ...prev, text: undefined }))
                    }
                  }}
                  placeholder="Enter link text"
                  className={`c-input ${errors.text ? "c-input--error" : ""}`}
                />
                {errors.text && <div className="text-red-600 text-sm mt-1">{errors.text}</div>}
              </div>
              <div className="c-form-field">
                <label htmlFor="link-url">
                  URL <span style={{ color: "red" }}>*</span>
                </label>
                <input
                  id="link-url"
                  type="text"
                  value={url}
                  onChange={(e) => {
                    setUrl(e.target.value)
                    if (errors.url) {
                      setErrors((prev) => ({ ...prev, url: undefined }))
                    }
                  }}
                  placeholder="https://example.com"
                  className={`c-input ${errors.url ? "c-input--error" : ""}`}
                />
                {errors.url && <div className="text-red-600 text-sm mt-1">{errors.url}</div>}
              </div>
            </form>
          </div>
          <div className="c-modal-footer c-modal-footer--right">
            <div style={{ display: "flex", gap: "8px" }}>
              <button type="button" onClick={handleClose} className="c-button c-button--secondary">
                Cancel
              </button>
              <button type="submit" onClick={handleSubmit} className="c-button c-button--primary">
                Insert Link
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LinkInsertModal
