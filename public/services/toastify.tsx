import React from "react"
import ReactDOM from "react-dom"

import { ToastContainer, toast, ToastContent, ToastOptions } from "react-toastify"
import "react-toastify/dist/ReactToastify.css"

let hasContainer = false

const setup = () => {
  if (!hasContainer) {
    hasContainer = true
    ReactDOM.render(<ToastContainer position={toast.POSITION.TOP_RIGHT} toastClassName="c-toast" />, document.getElementById("root-toastify"))
  }
}

export const success = (content: ToastContent, options?: ToastOptions) => {
  setup()
  toast.success(content, options)
}

export const error = (content: ToastContent, options?: ToastOptions) => {
  setup()
  toast.error(content, options)
}
