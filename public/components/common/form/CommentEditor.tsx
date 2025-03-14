// src/Tiptap.tsx
import { Editor } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import React, { useState } from "react"
import { EditorContent, useEditor } from "@tiptap/react"
import { Markdown } from "tiptap-markdown"
import Placeholder from "@tiptap/extension-placeholder"
import Document from "@tiptap/extension-document"
import Paragraph from "@tiptap/extension-paragraph"
import Text from "@tiptap/extension-text"

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
// import IconMarkdown from "@fider/assets/images/heroicons-h2.svg"
import IconBulletList from "@fider/assets/images/heroicons-bulletlist.svg"
import { Icon } from "@fider/components"

import suggestion from "./suggestion"
import { CustomMention } from "./CustomMention"
import { Trans } from "@lingui/react/macro"

// define your extension array
const MenuBar = ({ editor, isMarkdownMode, toggleMarkdownMode }: { editor: Editor | null; isMarkdownMode: boolean; toggleMarkdownMode: () => void }) => {
  if (!editor) {
    return null
  }

  return (
    <div className="c-editor-toolbar">
      <div className="c-editor-button-group">
        {/* Only show formatting buttons when not in markdown mode */}
        {!isMarkdownMode && (
          <>
            <button
              type="button"
              title="Heading 2"
              onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
              className={`c-editor-button ${editor.isActive("heading", { level: 2 }) ? "is-active" : ""}`}
            >
              <Icon sprite={IconH2} width="18" height="18" />
            </button>
            <button
              type="button"
              title="Heading 3"
              onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
              className={`c-editor-button ${editor.isActive("heading", { level: 3 }) ? "is-active" : ""}`}
            >
              <Icon sprite={IconH3} />
            </button>
            <button
              type="button"
              title="Bold"
              onClick={() => editor.chain().focus().toggleBold().run()}
              className={`c-editor-button ${editor.isActive("bold") ? "is-active" : ""}`}
            >
              <Icon sprite={IconBold} />
            </button>
            <button
              type="button"
              title="Italic"
              onClick={() => editor.chain().focus().toggleItalic().run()}
              className={`c-editor-button ${editor.isActive("italic") ? "is-active" : ""}`}
            >
              <Icon sprite={IconItalic} />
            </button>
            <button
              type="button"
              title="Strikethrough"
              onClick={() => editor.chain().focus().toggleStrike().run()}
              className={`c-editor-button ${editor.isActive("strike") ? "is-active" : ""}`}
            >
              <Icon sprite={IconStrike} />
            </button>
            <button
              type="button"
              title="BulletList"
              onClick={() => editor.chain().focus().toggleBulletList().run()}
              className={`c-editor-button ${editor.isActive("bulletList") ? "is-active" : ""}`}
            >
              <Icon sprite={IconBulletList} />
            </button>
            <button
              type="button"
              title="OrderedList"
              onClick={() => editor.chain().focus().toggleOrderedList().run()}
              className={`c-editor-button ${editor.isActive("orderedList") ? "is-active" : ""}`}
            >
              <Icon sprite={IconOrderedList} />
            </button>
            <button
              type="button"
              title="Code"
              onClick={() => editor.chain().focus().toggleCodeBlock().run()}
              className={`c-editor-button ${editor.isActive("codeBlock") ? "is-active" : ""}`}
            >
              <Icon sprite={IconCode} />
            </button>
            <button
              type="button"
              title="Mention"
              onClick={() => {
                // Insert @ character at current position and trigger mention suggestion
                editor.chain().focus().insertContent("@").run()
              }}
              className="c-editor-button"
            >
              <Icon sprite={IconAt} />
            </button>
          </>
        )}
        <button
          type="button"
          title={isMarkdownMode ? "Rich Text Mode" : "Markdown Mode"}
          onClick={toggleMarkdownMode}
          className={`c-editor-button ${isMarkdownMode ? "is-active" : ""} ml-auto text-xs`}
        >
          <span className="c-editor-button-text">
            <Trans id="editor.markdownmode">Markdown editor</Trans>
          </span>
        </button>
      </div>
    </div>
  )
}
interface CommentEditorProps {
  initialValue: string | null
  placeholder?: string
  onChange?: (value: string) => void
  onFocus?: () => void
}

const markdownToHtml = (markdownString: string) => {
  return markdownString
    .split("\n\n")
    .map((line: string) => `<p>${line}</p>`)
    .join("")
}

const Tiptap: React.FunctionComponent<CommentEditorProps> = (props) => {
  const [isRawMarkdownMode, setIsRawMarkdownMode] = useState(false)

  const getIntialContent = () => {
    if (isRawMarkdownMode) {
      return markdownToHtml(props.initialValue ?? "")
    } else {
      return props.initialValue ?? ""
    }
  }

  const [editorContent, setEditorContent] = useState(getIntialContent())

  const toggleMarkdownMode = () => {
    if (editor) {
      // Store current content before switching
      let currentContent
      if (isRawMarkdownMode) {
        currentContent = editor.getText()
      } else {
        currentContent = markdownToHtml(editor.storage.markdown.getMarkdown())
      }
      // Destroy current editor
      editor.destroy()
      setIsRawMarkdownMode(!isRawMarkdownMode)
      setEditorContent(currentContent)
    }
  }

  const updated = ({ editor }: { editor: Editor; transaction: any }): void => {
    const markdown = isRawMarkdownMode ? editor.getText() : editor.storage.markdown.getMarkdown()
    props.onChange && props.onChange(markdown)
  }

  const extensions = isRawMarkdownMode
    ? [
        // Minimal extensions for markdown mode
        Document,
        Paragraph,
        Text,
        Placeholder.configure({
          placeholder: props.placeholder ?? "Write your comment here...",
          emptyEditorClass: "tiptap-is-empty",
        }),
      ]
    : [
        StarterKit,
        Markdown.configure({
          html: true,
        }),
        CustomMention.configure({
          HTMLAttributes: {
            class: "mention",
          },
          suggestion,
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
      },
    },
    [isRawMarkdownMode, editorContent]
  ) // Re-initialize when mode changes

  return (
    <div className="fider-tiptap-editor">
      <MenuBar editor={editor} isMarkdownMode={isRawMarkdownMode} toggleMarkdownMode={toggleMarkdownMode} />
      <EditorContent editor={editor} />
    </div>
  )
}

const CommentEditor = React.memo(Tiptap, (prevProps, nextProps) => {
  return prevProps.placeholder === nextProps.placeholder
})

export default CommentEditor
