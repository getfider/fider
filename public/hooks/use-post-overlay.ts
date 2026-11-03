import { useEffect, useLayoutEffect, useState, useRef, useCallback } from "react"

interface UsePostOverlayOptions {
  basePath: string
  onPostClosed?: (postNumber: number) => void
}

export function usePostOverlay({ basePath, onPostClosed }: UsePostOverlayOptions) {
  const [selectedPostId, setSelectedPostId] = useState<number | null>(null)

  // Where the list was when a post was opened, so it can be restored on close
  const listScrollY = useRef(0)
  const listSearch = useRef("")
  const isPostDirty = useRef(false)
  const shownPostId = useRef<number | null>(null)

  const onPostClosedRef = useRef(onPostClosed)
  onPostClosedRef.current = onPostClosed

  const handlePostClick = useCallback((postNumber: number, slug: string) => {
    listScrollY.current = window.scrollY
    listSearch.current = window.location.search
    // Keep the scroll position on the list's history entry too, for the browser back button
    window.history.replaceState({ ...window.history.state, scrollY: window.scrollY }, "")
    window.history.pushState({ selectedPostId: postNumber }, "", `/posts/${postNumber}/${slug}`)
    setSelectedPostId(postNumber)
  }, [])

  const handleCloseOverlay = useCallback(() => {
    window.history.pushState({ scrollY: listScrollY.current }, "", `${basePath}${listSearch.current}`)
    setSelectedPostId(null)
  }, [basePath])

  const setIsPostDirty = useCallback((dirty: boolean) => {
    isPostDirty.current = dirty
  }, [])

  // Runs before the browser paints, so the page never shows at the wrong scroll position
  useLayoutEffect(() => {
    const previousPostId = shownPostId.current
    shownPostId.current = selectedPostId
    if (selectedPostId === previousPostId) {
      return
    }

    window.scrollTo(0, selectedPostId === null ? listScrollY.current : 0)

    if (previousPostId !== null && isPostDirty.current) {
      onPostClosedRef.current?.(previousPostId)
    }
    isPostDirty.current = false
  }, [selectedPostId])

  useEffect(() => {
    const handlePopState = (e: PopStateEvent) => {
      const path = window.location.pathname
      if (path === basePath || path === basePath.replace(/\/$/, "")) {
        if (typeof e.state?.scrollY === "number") {
          listScrollY.current = e.state.scrollY
        }
        setSelectedPostId(null)
      } else {
        const match = path.match(/^\/posts\/(\d+)/)
        if (match) {
          setSelectedPostId(parseInt(match[1], 10))
        }
      }
    }

    window.addEventListener("popstate", handlePopState)
    return () => window.removeEventListener("popstate", handlePopState)
  }, [basePath])

  return {
    selectedPostId,
    handlePostClick,
    handleCloseOverlay,
    setIsPostDirty,
  }
}
