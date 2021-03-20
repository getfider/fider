export const isTouch = (): boolean => {
  return typeof window === 'undefined' ? false :  "ontouchstart" in window || navigator.maxTouchPoints > 0 || navigator.msMaxTouchPoints > 0;
};
