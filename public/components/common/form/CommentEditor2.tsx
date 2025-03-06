// src/Tiptap.tsx
import { Editor } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import React from "react"
import { EditorContent, useEditor } from "@tiptap/react"
import { Markdown } from "tiptap-markdown"
import Mention from "@tiptap/extension-mention"

import "./CommentEditor.scss"

// At the top of the file, add imports for your icons
import IconH2 from "@fider/assets/images/heroicons-h2.svg"
import IconH3 from "@fider/assets/images/heroicons-h3.svg"
import IconItalic from "@fider/assets/images/heroicons-italic.svg"
import IconBold from "@fider/assets/images/heroicons-bold.svg"
import IconStrike from "@fider/assets/images/heroicons-strike.svg"
import { Icon } from "@fider/components"

import suggestion from "./suggestion"

// define your extension array
const extensions = [
  StarterKit,
  Markdown.configure({
    html: true,
  }),
  Mention.configure({
    HTMLAttributes: {
      class: "mention",
    },
    suggestion,
    renderText: (props) => {
      return "mushrooms"
    },
  }),
]

const MenuBar = ({ editor }: { editor: Editor | null }) => {
  if (!editor) {
    return null
  }

  return (
    <div className="c-editor-toolbar">
      <div className="c-editor-button-group">
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
      </div>
    </div>
  )
}

interface CommentEditorProps2 {
  initialValue: string | null
  placeholder?: string
  onChange?: (value: string) => void
}

const Tiptap: React.FunctionComponent<CommentEditorProps2> = (props) => {
  const updated = ({ editor }: { editor: Editor; transaction: any }): void => {
    const markdown = editor.storage.markdown.getMarkdown()
    console.log(markdown)
    props.onChange && props.onChange(markdown)
  }
  const initialValue = props.initialValue ?? "Hello"

  const editor = useEditor({
    extensions,
    content: initialValue,
    onUpdate: updated,
  })

  return (
    <div className="fider-tiptap-editor">
      <MenuBar editor={editor} />
      <EditorContent editor={editor} />
    </div>
  )
}

const MemoizedTiptap = React.memo(Tiptap)

export default MemoizedTiptap
