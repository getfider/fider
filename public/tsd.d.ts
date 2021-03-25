export {}

declare global {
  interface Window {
    ga?: (cmd: string, evt: string, args?: any) => void
    set: (key: string, value: any) => void
  }

  let __webpack_nonce__: string
  let __webpack_public_path__: string
}
