// import React, { useState } from "react"
// import { createEditor } from "slate"
// // Import the Slate components and React plugin.
// import { Slate, Editable, withReact } from "slate-react"
// //
// // TypeScript users only add this code
// import { BaseEditor, Descendant, Node } from "slate"
// import { ReactEditor } from "slate-react"
// import "./SlateEditor.scss"

// export type Paragraph = { type: "paragraph"; children: Text[] }
// export type Text = { text: string }

// declare module "slate2" {
//   interface CustomTypes {
//     Editor: BaseEditor & ReactEditor
//     Element: Paragraph
//     Text: Text
//   }
// }

// const emptyValue: Descendant[] = [
//   {
//     type: "paragraph",
//     children: [{ text: "" }],
//   },
// ]

// interface SlateEditorProps {
//   initialValue?: Descendant[]
//   disabled?: boolean
//   placeholder?: string
//   onChange?: (value: Descendant[]) => void
//   onFocus?: React.FocusEventHandler<HTMLDivElement>
//   className?: string
// }

// export const SlateEditor: React.FunctionComponent<SlateEditorProps> = (props) => {
//   const [editor] = useState(() => withReact(createEditor()))
//   console.log(props.initialValue)

//   return (
//     <Slate editor={editor} initialValue={props.initialValue || emptyValue} onChange={props.onChange}>
//       <Editable className="slate-editor" onFocus={props.onFocus} placeholder={props.placeholder} readOnly={props.disabled} />
//     </Slate>
//   )
// }

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
