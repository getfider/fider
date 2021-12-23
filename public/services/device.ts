export const isTouch = (): boolean => {
  return "ontouchstart" in window || navigator.maxTouchPoints > 0
}
