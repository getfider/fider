import { Editor } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import Link from "@tiptap/extension-link"
import React, { useState, useRef, useEffect } from "react"
import { EditorContent, useEditor } from "@tiptap/react"
import { Markdown } from "tiptap-markdown"
import Placeholder from "@tiptap/extension-placeholder"
import Document from "@tiptap/extension-document"
import Paragraph from "@tiptap/extension-paragraph"
import Text from "@tiptap/extension-text"
import HardBreak from "@tiptap/extension-hard-break"
import { i18n } from "@lingui/core"
import { useAllowedProtocols } from "@fider/hooks"

import "./CommentEditor.scss"

// At the top of the file, add imports for your icons
import IconH2 from "@fider/assets/images/heroicons-h2.svg"
import IconH3 from "@fider/assets/images/heroicons-h3.svg"
import IconItalic from "@fider/assets/images/heroicons-italic.svg"
import IconBold from "@fider/assets/images/heroicons-bold.svg"
import IconStrike from "@fider/assets/images/heroicons-strike.svg"
import IconCode from "@fider/assets/images/heroicons-code.svg"
import IconAt from "@fider/assets/images/heroicons-at.svg"
import IconOrderedList from "@fider/assets/images/heroicons-orderedlist.svg"
import IconBulletList from "@fider/assets/images/heroicons-bulletlist.svg"
import IconPhotograph from "@fider/assets/images/heroicons-photograph.svg"
import IconLink from "@fider/assets/images/heroicons-link.svg"
import { DisplayError, hasError, Icon, ValidationContext } from "@fider/components"
import { fileToBase64 } from "@fider/services"
import { generateBkey } from "@fider/services/bkey"
import { ImageUpload } from "@fider/models"
import { CustomImage } from "./CustomImage"

import suggestion from "./suggestion"
import { CustomMention } from "./CustomMention"
import LinkInsertModal from "./LinkInsertModal"
import { Trans } from "@lingui/react/macro"
import { classSet } from "@fider/services"

const MenuBar = ({
  editor,
  isMarkdownMode,
  toggleMarkdownMode,
  disabled,
  onImageUpload,
  onLinkClick,
}: {
  editor: Editor | null
  isMarkdownMode: boolean
  disabled: boolean
  toggleMarkdownMode: () => void
  onImageUpload: (file: File) => Promise<void>
  onLinkClick: (selectedText: string) => void
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null)

  if (!editor) {
    return null
  }

  const handleImageClick = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click()
    }
  }

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0]
      await onImageUpload(file)

      // Reset the input value so the same file can be selected again
      if (fileInputRef.current) {
        fileInputRef.current.value = ""
      }
    }
  }

  return (
    <div className="c-editor-toolbar">
      <div className="c-editor-button-group">
        {/* Only show formatting buttons when not in markdown mode */}
        {!isMarkdownMode && (
          <>
            <button
              disabled={disabled}
              type="button"
              title="Heading 2"
              onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
              className={`c-editor-button ${editor.isActive("heading", { level: 2 }) ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconH2} width="18" height="18" />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Heading 3"
              onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
              className={`c-editor-button ${editor.isActive("heading", { level: 3 }) ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconH3} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Bold"
              onClick={() => editor.chain().focus().toggleBold().run()}
              className={`c-editor-button ${editor.isActive("bold") ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconBold} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Italic"
              onClick={() => editor.chain().focus().toggleItalic().run()}
              className={`c-editor-button ${editor.isActive("italic") ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconItalic} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Strikethrough"
              onClick={() => editor.chain().focus().toggleStrike().run()}
              className={`c-editor-button ${editor.isActive("strike") ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconStrike} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="BulletList"
              onClick={() => editor.chain().focus().toggleBulletList().run()}
              className={`c-editor-button ${editor.isActive("bulletList") ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconBulletList} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="OrderedList"
              onClick={() => editor.chain().focus().toggleOrderedList().run()}
              className={`c-editor-button ${editor.isActive("orderedList") ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconOrderedList} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Code"
              onClick={() => editor.chain().focus().toggleCodeBlock().run()}
              className={`c-editor-button ${editor.isActive("codeBlock") ? "is-active" : ""} ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconCode} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Mention"
              onClick={() => {
                // Get the current cursor position
                const { from } = editor.state.selection
                // Get the character before the cursor
                const textBefore = editor.state.doc.textBetween(Math.max(0, from - 1), from)
                // Insert space before @ only if the previous character isn't a space or the cursor is at the beginning
                if (from === 0 || textBefore === " " || textBefore === "\n") {
                  editor.chain().focus().insertContent("@").run()
                } else {
                  editor.chain().focus().insertContent(" @").run()
                }
              }}
              className={`c-editor-button ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconAt} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Insert Image"
              onClick={handleImageClick}
              className={`c-editor-button ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconPhotograph} />
            </button>
            <button
              disabled={disabled}
              type="button"
              title="Insert Link"
              onClick={() => {
                // Get the text that was highlighted when the user clicked the link button
                const { from, to } = editor.state.selection
                const selectedText = editor.state.doc.textBetween(from, to)
                onLinkClick(selectedText)
              }}
              className={`c-editor-button ${disabled ? "is-disabled" : ""}`}
            >
              <Icon sprite={IconLink} />
            </button>
            <input ref={fileInputRef} type="file" accept="image/*" style={{ display: "none" }} onChange={handleFileChange} />
          </>
        )}
        <button
          disabled={disabled}
          type="button"
          title={isMarkdownMode ? "Rich Text Mode" : "Markdown Mode"}
          onClick={toggleMarkdownMode}
          className={`no-focus c-markdown-toggle ${isMarkdownMode ? "is-active" : ""} ${disabled ? "is-disabled" : ""} ml-auto text-xs`}
        >
          <span className="c-editor-button-text">
            {isMarkdownMode ? (
              <Trans id="editor.richtextmode">Switch to rich text editor</Trans>
            ) : (
              <Trans id="editor.markdownmode">Switch to markdown editor</Trans>
            )}
          </span>
        </button>
      </div>
    </div>
  )
}
interface CommentEditorProps {
  initialValue: string | null
  placeholder?: string
  onChange?: (value: string, plainText?: string) => void
  onFocus?: () => void
  disabled: boolean
  field: string
  onImageUploaded?: (upload: ImageUpload) => void
  onGetImageSrc?: (bkey: string) => string
  maxAttachments?: number
  maxImageSizeKB?: number
}

const markdownToHtml = (markdownString: string) => {
  const converted = markdownString
    .split("\n\n")
    .map((line: string) => `<p>${line}</p>`)
    .join("")
    .replace(/\\\n/g, "<br>")
    .replace(/\n/g, "<br>")

  console.log("markdownToHtml done: ", converted)

  return converted
}

const Tiptap: React.FunctionComponent<CommentEditorProps> = (props) => {
  const [isRawMarkdownMode, setIsRawMarkdownMode] = useState(false)
  const [imageUploads, setImageUploads] = useState<ImageUpload[]>([])
  const [isLinkModalOpen, setIsLinkModalOpen] = useState(false)
  const [selectedText, setSelectedText] = useState("")
  const allowedProtocols = useAllowedProtocols()

  // Use a ref instead of state for tracking document images
  // This avoids the async state update issue and prevents unnecessary re-renders
  const documentImagesRef = useRef<Map<string, boolean>>(new Map())

  const getIntialContent = () => {
    if (isRawMarkdownMode) {
      return markdownToHtml(props.initialValue ?? "")
    } else {
      return props.initialValue ?? ""
    }
  }

  const [editorContent, setEditorContent] = useState("")

  const toggleMarkdownMode = () => {
    if (editor) {
      // Store current content before switching
      let currentContent
      if (isRawMarkdownMode) {
        console.log("Getting text from markdown editor")
        currentContent = editor.getText()
      } else {
        console.log("Getting markdown from rich text editor")
        currentContent = editor.storage.markdown.getMarkdown()
      }
      console.log("Current content before toggle:", JSON.stringify(currentContent))
      // Destroy current editor
      editor.destroy()
      setIsRawMarkdownMode(!isRawMarkdownMode)
      setEditorContent(currentContent)
    }
  }

  // Handle image deletion
  const handleImageRemove = async (bkey: string) => {
    // Create an ImageUpload object with remove flag set to true
    const removeUpload: ImageUpload = {
      bkey,
      remove: true,
    }

    // Call the parent component's onImageUploaded prop with the removeUpload object
    if (props.onImageUploaded) {
      props.onImageUploaded(removeUpload)
    }
  }

  // Track all images in the document using the ref
  const trackDocumentImages = (editor: Editor) => {
    if (!editor) return

    // Create a new map for the current state
    const newImages = new Map<string, boolean>()

    // Find all image nodes in the document
    editor.state.doc.descendants((node) => {
      if (node.type.name === "customImage" || node.type.name === "image") {
        const bkey = node.attrs.bkey
        if (bkey) {
          newImages.set(bkey, true)
        }
      }
      return true
    })

    // Store the previous images for comparison
    const prevImages = new Map(documentImagesRef.current)

    // Update the ref with the new images
    documentImagesRef.current = newImages

    // Check for removed images
    checkForRemovedImages(prevImages, newImages)
  }

  // Check for removed images by comparing previous and current document images
  const checkForRemovedImages = (prevImages: Map<string, boolean>, currentImages: Map<string, boolean>) => {
    // Find images that were in the previous state but not in the current state
    prevImages.forEach((_, bkey) => {
      if (!currentImages.has(bkey)) {
        // This image was removed
        handleImageRemove(bkey)
      }
    })
  }

  const updated = ({ editor }: { editor: Editor; transaction: any }): void => {
    // Get the current markdown content
    const markdown = isRawMarkdownMode ? editor.getText() : editor.storage.markdown.getMarkdown()

    // Pass markdown to parent (plain text will be generated in ShareFeedback if needed)
    props.onChange && props.onChange(markdown)

    // Track the current state of images in the document
    // This will also check for removed images
    trackDocumentImages(editor)

    // Also pass any image uploads to the parent component
    if (props.onImageUploaded && imageUploads.length > 0) {
      imageUploads.forEach((upload) => {
        props.onImageUploaded && props.onImageUploaded(upload)
      })
      setImageUploads([])
    }
  }

  const validateImageUpload = (file: File): string => {
    // Default max size is 5MB (5 * 1024 KB)
    const maxSizeKB = props.maxImageSizeKB || 5 * 1024

    // Check file size
    if (file.size > maxSizeKB * 1024) {
      return i18n._({
        id: "validation.custom.maximagesize",
        values: { kilobytes: maxSizeKB },
        message: "The image size must be smaller than {kilobytes}KB.",
      })
    }

    // Check max attachments if specified
    if (props.maxAttachments) {
      if (documentImagesRef.current.size >= props.maxAttachments) {
        return i18n._({
          id: "validation.custom.maxattachments",
          values: { number: props.maxAttachments },
          message: "A maximum of {number} attachments are allowed.",
        })
      }
    }

    return "" // No error
  }

  const handleImageUpload = async (file: File) => {
    // Validate the image upload
    const errorMessage = validateImageUpload(file)
    if (errorMessage) {
      alert(errorMessage)
      return
    }

    try {
      const base64 = await fileToBase64(file)

      // Generate a bkey for this image that matches the server-side format
      const bkey = generateBkey(file.name)

      // Create an ImageUpload object to be sent to the server
      const newUpload: ImageUpload = {
        bkey,
        upload: {
          fileName: file.name,
          content: base64,
          contentType: file.type,
        },
        remove: false,
      }

      // Insert the image into the editor with the server-provided bkey
      if (editor) {
        editor
          .chain()
          .focus()
          .setImage({
            src: `data:${file.type};base64,${base64}`,
            alt: file.name,
            ...({ bkey } as any),
          })
          .run()

        // Add the bkey to the upload object for the parent component
        newUpload.bkey = bkey

        // Manually pass the upload to the parent component
        // since the updated() handler might not fire immediately
        if (props.onImageUploaded) {
          props.onImageUploaded(newUpload)
        }

        // Update the document images ref
        documentImagesRef.current.set(bkey, true)
      }
    } catch (error) {
      console.error("Error uploading image:", error)
    }
  }

  const handleInsertLink = (text: string, url: string) => {
    if (!editor) return

    editor.chain().focus().insertContent(`<a href="${url}" target="_blank" rel="noopener nofollow">${text}</a>`).run()
  }

  // Handle keyboard shortcuts
  const handleKeyDown = (event: KeyboardEvent) => {
    if (!editor) return

    const isMac = navigator.userAgent.toUpperCase().indexOf("MAC") >= 0
    const isCmdK = (isMac && event.metaKey && event.key === "k") || (!isMac && event.ctrlKey && event.key === "k")

    if (isCmdK) {
      event.preventDefault()
      // Get the text that was highlighted when the user pressed Cmd+K
      const { from, to } = editor.state.selection
      const selectedText = editor.state.doc.textBetween(from, to)
      setSelectedText(selectedText)
      setIsLinkModalOpen(true)
    }
  }

  const extensions = isRawMarkdownMode
    ? [
        // Minimal extensions for markdown mode
        Document,
        Paragraph,
        Text,
        HardBreak,
        Placeholder.configure({
          placeholder: props.placeholder ?? "Write your comment here...",
          emptyEditorClass: "tiptap-is-empty",
        }),
      ]
    : [
        StarterKit,
        Link.configure({
          openOnClick: true,
          autolink: true,
          defaultProtocol: "https",
          protocols: allowedProtocols,
          HTMLAttributes: {
            class: "text-link",
            target: "_blank",
            rel: "noopener nofollow",
          },
        }),
        Markdown.configure({
          html: false,
          breaks: true,
        }),
        CustomMention.configure({
          HTMLAttributes: {
            class: "mention",
          },
          suggestion,
        }),
        CustomImage.configure({
          HTMLAttributes: {},
          allowBase64: true,
          onImageUpload: (upload) => {
            if (props.onImageUploaded) {
              // Initialize other required properties
              props.onImageUploaded(upload)
            }
          },
          onImageRemove: (id) => {
            // This is called when an image is removed from the editor
            handleImageRemove(id)
          },
          onGetImageSrc: (bkey) => {
            if (props.onGetImageSrc) {
              const imageSrc = props.onGetImageSrc(bkey)
              if (imageSrc) {
                return "data:image/jpeg;base64," + imageSrc
              }
            }
            return ""
          },
        }),
        Placeholder.configure({
          placeholder: props.placeholder ?? "Write your comment here...",
          emptyEditorClass: "tiptap-is-empty",
        }),
      ]

  const editor = useEditor(
    {
      extensions,
      content: editorContent,
      onUpdate: updated,
      onFocus: () => {
        if (props.onFocus) {
          props.onFocus()
        }
      },
      editorProps: {
        attributes: {
          class: isRawMarkdownMode ? "markdown-mode no-focus" : "no-focus",
        },
        handlePaste: (view, event) => {
          // Check if the clipboard has files
          if (event.clipboardData && event.clipboardData.files && event.clipboardData.files.length > 0) {
            // Get the first file (assuming it's an image)
            const file = event.clipboardData.files[0]

            // Check if it's an image
            if (file.type.startsWith("image/")) {
              // Prevent the default paste behavior
              event.preventDefault()

              // Validate and upload the image
              const errorMessage = validateImageUpload(file)
              if (errorMessage) {
                alert(errorMessage)
                return true
              }

              // Upload the image
              handleImageUpload(file)
              return true
            }
          }

          // Check for pasted image data URLs
          if (event.clipboardData) {
            const items = event.clipboardData.items

            for (let i = 0; i < items.length; i++) {
              if (items[i].type.indexOf("image") !== -1) {
                // Get the image as a blob
                const blob = items[i].getAsFile()

                if (blob) {
                  // Prevent the default paste behavior
                  event.preventDefault()

                  // Validate and upload the image
                  const errorMessage = validateImageUpload(blob)
                  if (errorMessage) {
                    alert(errorMessage)
                    return true
                  }

                  // Upload the image
                  handleImageUpload(blob)
                  return true
                }
              }
            }
          }

          return false
        },
      },
    },
    [isRawMarkdownMode, editorContent]
  ) // Re-initialize when mode changes

  // Initialize document images when editor is ready
  useEffect(() => {
    if (editor) {
      trackDocumentImages(editor)
    }
  }, [editor])

  // Add keyboard event listener for shortcuts
  useEffect(() => {
    if (editor) {
      const editorElement = editor.view.dom
      editorElement.addEventListener("keydown", handleKeyDown)

      return () => {
        editorElement.removeEventListener("keydown", handleKeyDown)
      }
    }
  }, [editor])

  return (
    <ValidationContext.Consumer>
      {(ctx) => (
        <div>
          <div
            className={classSet({
              "fider-tiptap-editor": true,
              "m-error": hasError(props.field, ctx.error),
            })}
          >
            <MenuBar
              disabled={props.disabled}
              editor={editor}
              isMarkdownMode={isRawMarkdownMode}
              toggleMarkdownMode={toggleMarkdownMode}
              onImageUpload={handleImageUpload}
              onLinkClick={(text) => {
                setSelectedText(text)
                setIsLinkModalOpen(true)
              }}
            />
            <EditorContent editor={editor} data-testid="tiptap-editor" />
          </div>
          <DisplayError fields={[props.field]} error={ctx.error} />
          <LinkInsertModal isOpen={isLinkModalOpen} onClose={() => setIsLinkModalOpen(false)} onInsertLink={handleInsertLink} selectedText={selectedText} />
        </div>
      )}
    </ValidationContext.Consumer>
  )
}

const CommentEditor = React.memo(Tiptap, (prevProps, nextProps) => {
  return prevProps.placeholder === nextProps.placeholder
})

export default CommentEditor
