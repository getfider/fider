import "./Dropdown.scss"

import React, { createContext, useContext, useEffect, useRef, useState } from "react"
import { classSet } from "@fider/services"

interface DropdownListItemProps {
  href?: string
  onClick?: () => void
  className?: string
  children: React.ReactNode
}

const ListItem = (props: DropdownListItemProps) => {
  const ctx = useContext(DropdownContext)
  const handleClick = () => {
    if (props.onClick) {
      props.onClick()
    }

    ctx?.close()
  }

  if (props.href) {
    return (
      <a href={props.href} className={`c-dropdown__listitem ${props.className}`}>
        {props.children}
      </a>
    )
  }

  return (
    <div onClick={handleClick} className={`c-dropdown__listitem ${props.className}`}>
      {props.children}
    </div>
  )
}

const Divider = () => {
  return <hr className="c-dropdown__divider" />
}

interface DropdownProps {
  renderHandle: JSX.Element
  position?: "left" | "right"
  onToggled?: (isOpen: boolean) => void
  children: React.ReactNode
  wide?: boolean
  fullsceenSm?: boolean
}

interface DropdownContextFuncs {
  close(): void
}

export const DropdownContext = createContext<DropdownContextFuncs | null>(null)
DropdownContext.displayName = "DropdownContext"

export const Dropdown = (props: DropdownProps) => {
  const node = useRef<HTMLDivElement | null>(null)
  const [isOpen, setIsOpen] = useState(false)
  const position = props.position || "right"

  const changeToggleState = (newState: boolean) => {
    setIsOpen(newState)
    if (props.onToggled) {
      props.onToggled(newState)
    }
  }

  const toggleIsOpen = () => {
    changeToggleState(!isOpen)
  }

  const close = () => {
    changeToggleState(false)
  }

  const handleClick = (e: MouseEvent) => {
    if (node.current && node.current.contains(e.target as Node)) {
      return
    }

    close()
  }

  useEffect(() => {
    document.addEventListener("mousedown", handleClick)

    return () => {
      document.removeEventListener("mousedown", handleClick)
    }
  }, [])

  const listClassName = classSet({
    "c-dropdown__list--wide": props.wide,
    "c-dropdown__list shadow-lg": true,
    "c-dropdown__list--fullscreen-small": props.fullsceenSm,
    [`c-dropdown__list--${position}`]: position === "left",
  })

  return (
    <DropdownContext.Provider value={{ close }}>
      <div ref={node} className="c-dropdown">
        {/* Only render the handle when closed */}
        {!isOpen && (
          <button type="button" className="c-dropdown__handle" onClick={toggleIsOpen}>
            {props.renderHandle}
          </button>
        )}
        {isOpen && (
          <div className={listClassName} style={{ position: "static" }}>
            {props.children}
          </div>
        )}
      </div>
    </DropdownContext.Provider>
  )
}

Dropdown.ListItem = ListItem
Dropdown.Divider = Divider
