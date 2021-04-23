// NOT USED, BUT WE MIGHT NEED IN FUTURE

// import "./Select2.scss"

// import React, { useEffect, useRef, useState } from "react"
// import IconSelector from "@fider/assets/images/heroicons-selector.svg"
// import { HStack } from "@fider/components/layout"
// import { Icon } from "@fider/components"

// interface Select2Props<T> {
//   items: T[]
//   itemKey: keyof T
//   onSelect: (item: T) => void
//   renderItem: (item: T) => JSX.Element
//   renderHandle: () => JSX.Element
// }

// export const Select2 = <T extends unknown>(props: Select2Props<T>) => {
//   const node = useRef<HTMLDivElement | null>(null)
//   const [isOpen, setIsOpen] = useState(false)

//   const toggleIsOpen = () => {
//     setIsOpen(!isOpen)
//   }

//   const handleClick = (e: MouseEvent) => {
//     if (node.current && node.current.contains(e.target as Node)) {
//       return
//     }
//     setIsOpen(false)
//   }

//   useEffect(() => {
//     document.addEventListener("mousedown", handleClick)

//     return () => {
//       document.removeEventListener("mousedown", handleClick)
//     }
//   }, [])

//   const handleChange = (i: T) => () => {
//     props.onSelect(i)
//     setIsOpen(false)
//   }

//   return (
//     <div ref={node} className="c-select2">
//       <button type="button" className="c-select2__handle w-full" onClick={toggleIsOpen}>
//         <HStack justify="between">
//           {props.renderHandle()}
//           <Icon sprite={IconSelector} className="h-5 w-5 text-gray-500" />
//         </HStack>
//       </button>
//       {isOpen && (
//         <ul className="c-select2__list shadow-lg w-full">
//           {props.items.map((i) => {
//             const key = (i[props.itemKey] as unknown) as string
//             return (
//               <li onClick={handleChange(i)} key={key}>
//                 {props.renderItem(i)}
//               </li>
//             )
//           })}
//         </ul>
//       )}
//     </div>
//   )
// }
