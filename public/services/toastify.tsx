import { I18nProvider } from "@lingui/react"
import React from "react"
import ReactDOM from "react-dom/client"
import { i18n } from "@lingui/core"

import { ToastContainer, toast, ToastContent, ToastOptions } from "react-toastify"
import "react-toastify/dist/ReactToastify.css"

let hasContainer = false

const setup = () => {
  if (!hasContainer) {
    hasContainer = true
    const rootElement = document.getElementById("root-toastify")
    if (rootElement) {
      const root = ReactDOM.createRoot(rootElement)
      root.render(
        <I18nProvider i18n={i18n}>
          <ToastContainer position={toast.POSITION.TOP_RIGHT} />
        </I18nProvider>
      )
    }
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
