const toastify = () => import(/* webpackChunkName: "toastify" */ "./toastify")

export const success = (content: string | JSX.Element) => {
  return toastify().then((toast) => {
    toast.success(content)
  })
}

export const error = (content: string | JSX.Element) => {
  return toastify().then((toast) => {
    toast.error(content)
  })
}
