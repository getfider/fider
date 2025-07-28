import "./Modal.scss"

import React, { useEffect, useRef } from "react"
import ReactDOM from "react-dom"
import { classSet } from "@fider/services"
import { Icon } from "@fider/components"
import IconX from "@fider/assets/images/heroicons-x.svg"

interface ModalWindowProps {
  children?: React.ReactNode
  className?: string
  isOpen: boolean
  size?: "small" | "large" | "fluid" | "fullscreen"
  canClose?: boolean
  center?: boolean
  onClose: () => void
}

interface ModalFooterProps {
  align?: "left" | "center" | "right"
  children?: React.ReactNode
}

const ModalWindow: React.FunctionComponent<ModalWindowProps> = ({ size = "small", canClose = true, center = true, ...props }) => {
  const root = useRef<HTMLElement>(document.getElementById("root-modal"))

  useEffect(() => {
    if (props.isOpen) {
      document.body.style.overflow = "hidden"
      document.addEventListener("keydown", keyDown, false)
    } else {
      document.body.style.overflow = ""
      document.removeEventListener("keydown", keyDown, false)
    }
  }, [props.isOpen])

  const swallow = (evt: React.MouseEvent<HTMLDivElement>) => {
    evt.stopPropagation()
  }

  const keyDown = (event: KeyboardEvent) => {
    if (event.code === "Escape") {
      close()
    }
  }

  const close = () => {
    if (canClose) {
      props.onClose()
    }
  }

  if (!props.isOpen || !root.current) {
    return null
  }

  const className = classSet({
    "c-modal-window": true,
    [`${props.className}`]: !!props.className,
    "c-modal-window--center": center,
    [`c-modal-window--${size}`]: true,
  })

  return ReactDOM.createPortal(
    <div aria-disabled={true} className="c-modal-dimmer" onClick={close}>
      <div className="c-modal-scroller">
        <div className={className} data-testid="modal" onClick={swallow}>
          {props.children}
        </div>
      </div>
    </div>,
    root.current
  )
}
const Header = (props: { children: React.ReactNode }) => <div className="c-modal-header">{props.children}</div>
const Content = (props: { children: React.ReactNode }) => <div className="c-modal-content">{props.children}</div>
const Footer = (props: ModalFooterProps) => {
  const align = props.align || "right"
  const className = classSet({
    "c-modal-footer": true,
    [`c-modal-footer--${align}`]: true,
  })
  return <div className={className}>{props.children}</div>
}

export const CloseIcon = ({ closeModal }: { closeModal: () => void }) => (
  <Icon sprite={IconX} height="30" width="30" onClick={closeModal} className="c-modal-closeicon" />
)

export const Modal = {
  Window: ModalWindow,
  Header,
  Content,
  Footer,
}
