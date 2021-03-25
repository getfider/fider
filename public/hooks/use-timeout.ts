import { useRef, useEffect } from "react"

type CallbackFunction = () => void

export function useTimeout(callback: CallbackFunction, delay: number) {
  const savedCallback = useRef<CallbackFunction>()

  useEffect(() => {
    savedCallback.current = callback
  })

  useEffect(() => {
    function tick() {
      if (savedCallback.current) {
        savedCallback.current()
      }
    }
    const timer = window.setTimeout(tick, delay)
    return function cleanup() {
      window.clearTimeout(timer)
    }
  }, [delay])
}
