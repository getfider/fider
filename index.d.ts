declare interface Window {
  ga?: (cmd: string, evt: string, args?: any) => void
  set: (key: string, value: any) => void
}

interface SpriteSymbol {
  id: string
  viewBox: string
  url: string
}

declare let __webpack_nonce__: string
declare let __webpack_public_path__: string

declare module "*.svg" {
  const content: SpriteSymbol
  export default content
}
