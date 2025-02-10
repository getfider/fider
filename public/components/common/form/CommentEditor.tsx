import React, { useMemo, useCallback, useRef, useEffect, useState, Fragment, ReactNode } from "react"
import { Element, Editor, Transforms, Range, createEditor, Descendant, BaseEditor, BaseRange, Node } from "slate"
import { HistoryEditor, withHistory } from "slate-history"
import { Slate, Editable, ReactEditor, withReact, useSelected, useFocused, RenderPlaceholderProps } from "slate-react"

export const IS_MAC = typeof navigator !== "undefined" && /Mac OS X/.test(navigator.userAgent)

import ReactDOM from "react-dom"
import { UserNames } from "@fider/models"
import { actions } from "@fider/services"

import "./CommentEditor.scss"
import { useFider } from "@fider/hooks"

export type TextType = { text: string }

interface RenderElementProps {
  attributes: React.HTMLAttributes<HTMLElement>
  element: CustomElement
  children: React.ReactNode
}

type MentionElement = {
  type: "mention"
  character: string
  children: TextType[]
}

export type ParagraphElement = {
  type: "paragraph"
  children: Descendant[]
}

type CustomElement = MentionElement | ParagraphElement

type CustomEditor = BaseEditor &
  ReactEditor &
  HistoryEditor & {
    nodeToDecorations?: Map<Element, Range[]>
  }

const emptyValue: Descendant[] = [
  {
    type: "paragraph",
    children: [{ text: "" }],
  },
]

declare module "slate" {
  interface CustomTypes {
    Editor: CustomEditor
    Element: CustomElement
    Text: TextType
    Range: BaseRange & {
      [key: string]: unknown
    }
  }
}

const Portal = ({ children }: { children?: ReactNode }) => {
  return typeof document === "object" ? ReactDOM.createPortal(children, document.body) : null
}

interface CommentEditorProps {
  initialValue?: string
  placeholder?: string
  onChange?: (value: string) => void
  onFocus?: React.FocusEventHandler<HTMLDivElement>
  className?: string
}

const Placeholder = ({ attributes, children }: RenderPlaceholderProps) => {
  return (
    <span
      {...attributes}
      className="slate-editor--placeholder"
      style={{
        position: "absolute",
        opacity: 0.3,
        userSelect: "none",
        pointerEvents: "none",
      }}
    >
      {children}
    </span>
  )
}

export const CommentEditor: React.FunctionComponent<CommentEditorProps> = (props) => {
  const [users, setUsers] = useState<UserNames[]>([])
  const ref = useRef<HTMLDivElement | null>(null)
  const [target, setTarget] = useState<Range | undefined>()
  const [index, setIndex] = useState(0)
  const [search, setSearch] = useState("")
  const fider = useFider()

  // Fetch users when component mounts
  useEffect(() => {
    const loadUsers = async () => {
      const result = await actions.getTaggableUsers("")
      if (result.ok) {
        setUsers(result.data)
      }
    }
    if (fider.session.isAuthenticated) {
      loadUsers()
    }
  }, [])

  const renderElement = useCallback((props: RenderElementProps) => <SlateElement {...props} />, [])
  const editor = useMemo(() => withMentions(withReact(withHistory(createEditor()))), [])

  const filteredUsers = users.filter((user) => user.name.toLowerCase().startsWith(search.toLowerCase())).slice(0, 10)

  const onKeyDown = useCallback(
    (event: React.KeyboardEvent<HTMLDivElement>) => {
      if (target && filteredUsers.length > 0) {
        switch (event.key) {
          case "ArrowDown": {
            event.preventDefault()
            const prevIndex = index >= filteredUsers.length - 1 ? 0 : index + 1
            setIndex(prevIndex)
            break
          }
          case "ArrowUp": {
            event.preventDefault()
            const nextIndex = index <= 0 ? filteredUsers.length - 1 : index - 1
            setIndex(nextIndex)
            break
          }
          case "Tab":
          case "Enter":
            event.preventDefault()
            Transforms.select(editor, target)
            insertMention(editor, filteredUsers[index])
            setTarget(undefined)
            break
          case "Escape":
            event.preventDefault()
            setTarget(undefined)
            break
        }
      }
    },
    [filteredUsers, editor, index, target]
  )

  // Where to show the mentions portal
  useEffect(() => {
    if (target && filteredUsers.length > 0 && ref.current) {
      const el = ref.current
      const domRange = ReactEditor.toDOMRange(editor, target)
      const rect = domRange.getBoundingClientRect()
      el.style.top = `${rect.top + window.pageYOffset + 24}px`
      el.style.left = `${rect.left + window.pageXOffset}px`
    }
  }, [filteredUsers.length, editor, index, search, target])

  const initialValue = props.initialValue ? deserialize(props.initialValue) : emptyValue

  return (
    <Slate
      editor={editor}
      initialValue={initialValue}
      onChange={(descendant) => {
        const { selection } = editor

        if (selection && Range.isCollapsed(selection)) {
          const [start] = Range.edges(selection)
          const wordBefore = Editor.before(editor, start, { unit: "word" })
          const before = wordBefore && Editor.before(editor, wordBefore)
          const beforeRange = before && Editor.range(editor, before, start)
          const beforeText = beforeRange && Editor.string(editor, beforeRange)
          const beforeMatch = beforeText && beforeText.match(/^@(\w+)$/)
          const after = Editor.after(editor, start)
          const afterRange = Editor.range(editor, start, after)
          const afterText = Editor.string(editor, afterRange)
          const afterMatch = afterText.match(/^(\s|$)/)

          if (beforeMatch && afterMatch) {
            setTarget(beforeRange)
            setSearch(beforeMatch[1])
            setIndex(0)
            return
          }

          console.log(JSON.stringify(descendant))

          props.onChange && props.onChange(serialize(descendant))
        }

        setTarget(undefined)
      }}
    >
      <Editable
        readOnly={false}
        className="slate-editor"
        renderElement={renderElement}
        onKeyDown={onKeyDown}
        onFocus={props.onFocus}
        placeholder={props.placeholder}
        renderPlaceholder={Placeholder}
      />
      {target && filteredUsers.length > 0 && (
        <Portal>
          <div ref={ref} className="slate-editor--mentions" data-cy="mentions-portal">
            {filteredUsers.map((user, i) => (
              <div
                key={user.id}
                onClick={() => {
                  Transforms.select(editor, target)
                  insertMention(editor, user)
                  setTarget(undefined)
                }}
                style={{
                  padding: "6px 12px",
                  borderRadius: "3px",
                  cursor: "pointer",
                  background: i === index ? "#B4D5FF" : "transparent",
                }}
              >
                {user.name}
              </div>
            ))}
          </div>
        </Portal>
      )}
    </Slate>
  )
}

const withMentions = (editor: CustomEditor) => {
  const { isInline, isVoid, markableVoid } = editor

  editor.isInline = (element: CustomElement) => {
    return element.type === "mention" ? true : isInline(element)
  }

  editor.isVoid = (element: CustomElement) => {
    return element.type === "mention" ? true : isVoid(element)
  }

  editor.markableVoid = (element: CustomElement) => {
    return element.type === "mention" || markableVoid(element)
  }

  return editor
}

const insertMention = (editor: Editor, user: UserNames) => {
  const newMention = { ...user, isNew: true }
  const mention: MentionElement = {
    type: "mention",
    character: user.name,
    children: [{ text: "@" + JSON.stringify(newMention) }],
  }
  Transforms.insertNodes(editor, mention)
  Transforms.move(editor)
}

const SlateElement = (props: RenderElementProps) => {
  const { attributes, children, element } = props
  switch (element.type) {
    case "mention":
      return <Mention {...(props as { attributes: React.HTMLAttributes<HTMLSpanElement>; children: React.ReactNode; element: MentionElement })} />
    default:
      return <p {...attributes}>{children}</p>
  }
}

const Mention = ({
  attributes,
  children,
  element,
}: {
  attributes: React.HTMLAttributes<HTMLSpanElement>
  children: React.ReactNode
  element: MentionElement
}) => {
  const selected = useSelected()
  const focused = useFocused()
  const style: React.CSSProperties = {
    boxShadow: selected && focused ? "0 0 0 2px #B4D5FF" : "none",
  }
  // See if our empty text child has any styling marks applied and apply those
  return (
    <span {...attributes} className="slate-editor--mention" contentEditable={false} data-cy={`mention-${element.character.replace(" ", "-")}`} style={style}>
      {/* Prevent Chromium from interrupting IME when moving the cursor */}
      {/* 1. span + inline-block 2. div + contenteditable=false */}
      <div contentEditable={false}>
        {IS_MAC ? (
          // Mac OS IME https://github.com/ianstormtaylor/slate/issues/3490
          <Fragment>
            {children}@{element.character}
          </Fragment>
        ) : (
          // Others like Android https://github.com/ianstormtaylor/slate/pull/5360
          <Fragment>
            @{element.character}
            {children}
          </Fragment>
        )}
      </div>
    </span>
  )
}

const serialize = (nodes: Descendant[]): string => {
  return nodes.map((n) => Node.string(n)).join("\n")
}

const deserialize = (markdown: string): Descendant[] => {
  return markdown.split("\n").map((line) => {
    const children: Descendant[] = []
    const regex = /@{\\?"id\\?":\d+,\\?"name\\?":\\?"[^"]+\\?"(,\\?"isNew\\?":[^}]+)?}/g
    let lastIndex = 0

    let match
    while ((match = regex.exec(line)) !== null) {
      // Add text before the mention
      if (match.index > lastIndex) {
        children.push({ text: line.slice(lastIndex, match.index) })
      }

      // Handle mention
      try {
        const jsonStr = match[0].replace(/\\/g, "").slice(1)
        const mentionData = JSON.parse(jsonStr)
        children.push({
          type: "mention",
          character: mentionData.name,
          children: [{ text: match[0] }],
        })
      } catch (err) {
        console.error("Error parsing mention:", err)
        // Just add the text as a normal paragraph
        children.push({ text: line })
      }

      lastIndex = match.index + match[0].length
    }

    // Add remaining text after last mention
    if (lastIndex <= line.length) {
      children.push({ text: line.slice(lastIndex) })
    }

    return {
      type: "paragraph",
      children,
    }
  })
}
