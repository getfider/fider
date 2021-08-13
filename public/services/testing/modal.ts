export const setupModalRoot = () => {
  let portalRoot = document.getElementById("root-modal")
  if (portalRoot) {
    document.body.removeChild(portalRoot)
  }

  portalRoot = document.createElement("div")
  portalRoot.setAttribute("id", "root-modal")
  document.body.appendChild(portalRoot)
}
