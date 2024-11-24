import React, { useMemo, useCallback, useRef, useEffect, useState, Fragment, ReactNode } from "react"
import { Editor, Transforms, Range, createEditor, Descendant, BaseEditor, BaseRange } from "slate"
import { HistoryEditor, withHistory } from "slate-history"
import { Slate, Editable, ReactEditor, withReact, useSelected, useFocused } from "slate-react"

export const IS_MAC = typeof navigator !== "undefined" && /Mac OS X/.test(navigator.userAgent)

import ReactDOM from "react-dom"
import { UserNames } from "@fider/models"
import { actions } from "@fider/services"

import "./SlateEditor.scss"

type CustomText = {
  bold?: boolean
  italic?: boolean
  code?: boolean
  text: string
}

type EmptyText = {
  text: string
}

interface RenderElementProps {
  attributes: React.HTMLAttributes<HTMLElement>
  element: CustomElement
  children: React.ReactNode
}

type MentionElement = {
  type: "mention"
  character: string
  children: CustomText[]
}

type ParagraphElement = {
  type: "paragraph"
  align?: string
  children: Descendant[]
}

type CustomElement = MentionElement | ParagraphElement

type CustomEditor = BaseEditor &
  ReactEditor &
  HistoryEditor & {
    nodeToDecorations?: Map<Element, Range[]>
  }

declare module "slate" {
  interface CustomTypes {
    Editor: CustomEditor
    Element: CustomElement
    Text: CustomText | EmptyText
    Range: BaseRange & {
      [key: string]: unknown
    }
  }
}

const Portal = ({ children }: { children?: ReactNode }) => {
  return typeof document === "object" ? ReactDOM.createPortal(children, document.body) : null
}

export const MentionExample = () => {
  const [users, setUsers] = useState<UserNames[]>([])
  const ref = useRef<HTMLDivElement | null>(null)
  const [target, setTarget] = useState<Range | undefined>()
  const [index, setIndex] = useState(0)
  const [search, setSearch] = useState("")

  // Fetch users when component mounts
  useEffect(() => {
    const loadUsers = async () => {
      const result = await actions.getTaggableUsers("")
      if (result.ok) {
        setUsers(result.data)
      }
    }
    loadUsers()
  }, [])

  const renderElement = useCallback((props: RenderElementProps) => <Element {...props} />, [])
  // const renderLeaf = useCallback((props: { attributes: any; leaf: any; children: any }) => <Leaf {...props} />, [])
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
            insertMention(editor, filteredUsers[index].name)
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

  return (
    <Slate
      editor={editor}
      initialValue={initialValue}
      onChange={() => {
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
        }

        setTarget(undefined)
      }}
    >
      <Editable className="slate-editor" renderElement={renderElement} onKeyDown={onKeyDown} placeholder="Enter some text..." />
      {target && filteredUsers.length > 0 && (
        <Portal>
          <div ref={ref} className="slate-editor--mentions" data-cy="mentions-portal">
            {filteredUsers.map((char, i) => (
              <div
                key={char.name}
                onClick={() => {
                  Transforms.select(editor, target)
                  insertMention(editor, char.name)
                  setTarget(undefined)
                }}
                style={{
                  padding: "6px 12px",
                  borderRadius: "3px",
                  cursor: "pointer",
                  background: i === index ? "#B4D5FF" : "transparent",
                }}
              >
                {char.name}
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

const insertMention = (editor: Editor, character: string) => {
  const mention: MentionElement = {
    type: "mention",
    character,
    children: [{ text: "" }],
  }
  Transforms.insertNodes(editor, mention)
  Transforms.move(editor)
}

const Element = (props: RenderElementProps) => {
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
const initialValue: Descendant[] = [
  {
    type: "paragraph",
    children: [
      {
        text: "This example shows how you might implement a simple ",
      },
      {
        text: "@-mentions",
        bold: true,
      },
      {
        text: " feature that lets users autocomplete mentioning a user by their username. Which, in this case means Star Wars characters. The ",
      },
      {
        text: "mentions",
        bold: true,
      },
      {
        text: " are rendered as ",
      },
      {
        text: "void inline elements",
        code: true,
      },
      {
        text: " inside the document.",
      },
    ],
  },
  {
    type: "paragraph",
    children: [
      { text: "Try mentioning characters, like " },
      {
        type: "mention",
        character: "R2-D2",
        children: [{ text: "", bold: true }],
      },
      { text: " or " },
      {
        type: "mention",
        character: "Mace Windu",
        children: [{ text: "" }],
      },
      { text: "!" },
    ],
  },
]

// export const Serialize = (nodes: Descendant[]): string => {
//   return nodes.map((n) => Node.string(n)).join("\n")
// }

// export const Deserialize = (value: string): Descendant[] => {
//   return value.split("\n").map((line) => {
//     return {
//       children: [{ text: line }],
//     } as Descendant
//   })
// }
