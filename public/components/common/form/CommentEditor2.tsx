// src/Tiptap.tsx
import { Editor } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import React from "react"
import { EditorContent, useEditor } from "@tiptap/react"
import { Markdown } from "tiptap-markdown"
import Placeholder from "@tiptap/extension-placeholder"
// import Mention from "@tiptap/extension-mention"

import "./CommentEditor.scss"

// At the top of the file, add imports for your icons
import IconH2 from "@fider/assets/images/heroicons-h2.svg"
import IconH3 from "@fider/assets/images/heroicons-h3.svg"
import IconItalic from "@fider/assets/images/heroicons-italic.svg"
import IconBold from "@fider/assets/images/heroicons-bold.svg"
import IconStrike from "@fider/assets/images/heroicons-strike.svg"
import IconCode from "@fider/assets/images/heroicons-code.svg"
import IconOrderedList from "@fider/assets/images/heroicons-orderedlist.svg"
import IconBulletList from "@fider/assets/images/heroicons-bulletlist.svg"
import { Icon } from "@fider/components"

import suggestion from "./suggestion"
import { CustomMention } from "./CustomMention"
// import CustomMarkdown from "./MentionMarkdown"

// define your extension array

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
      </div>
    </div>
  )
}

interface CommentEditorProps2 {
  initialValue: string | null
  placeholder?: string
  onChange?: (value: string) => void
  onFocus?: () => void
}

const Tiptap: React.FunctionComponent<CommentEditorProps2> = (props) => {
  const updated = ({ editor }: { editor: Editor; transaction: any }): void => {
    const markdown = editor.storage.markdown.getMarkdown()
    props.onChange && props.onChange(markdown)
  }
  const initialValue = props.initialValue ?? ""

  const extensions = [
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

  const editor = useEditor({
    extensions,
    content: initialValue,
    onUpdate: updated,
    onFocus: () => {
      if (props.onFocus) {
        props.onFocus()
      }
    },
    editorProps: {
      attributes: {
        class: "no-focus",
      },
    },
  })

  return (
    <div className="fider-tiptap-editor">
      <MenuBar editor={editor} />
      <EditorContent editor={editor} />
    </div>
  )
}

const MemoizedTiptap = React.memo(Tiptap, (prevProps, nextProps) => {
  return prevProps.placeholder === nextProps.placeholder
})

export default MemoizedTiptap
