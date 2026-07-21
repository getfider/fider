import { useEffect, useState, useRef, useCallback } from "react"

interface UsePostOverlayOptions {
  basePath: string
  onPostClosed?: (postNumber: number) => void
}

export function usePostOverlay({ basePath, onPostClosed }: UsePostOverlayOptions) {
  const [selectedPostId, setSelectedPostId] = useState<number | null>(null)
  const [savedScrollPosition, setSavedScrollPosition] = useState<number>(0)
  const [isPostDirty, setIsPostDirty] = useState(false)
  const [savedSearch, setSavedSearch] = useState("")
  const [lastOpenedPostId, setLastOpenedPostId] = useState<number | null>(null)

  const onPostClosedRef = useRef(onPostClosed)
  onPostClosedRef.current = onPostClosed

  const handlePostClick = useCallback((postNumber: number, slug: string) => {
    setSavedScrollPosition(window.scrollY)
    setSavedSearch(window.location.search)
    setSelectedPostId(postNumber)
    setLastOpenedPostId(postNumber)
    setIsPostDirty(false)
    window.history.pushState({ selectedPostId: postNumber }, "", `/posts/${postNumber}/${slug}`)
  }, [])

  const handleCloseOverlay = useCallback(() => {
    setSelectedPostId(null)
    window.history.pushState({}, "", `${basePath}${savedSearch}`)
  }, [basePath, savedSearch])

  useEffect(() => {
    if (selectedPostId === null && lastOpenedPostId !== null) {
      if (isPostDirty && onPostClosedRef.current) {
        onPostClosedRef.current(lastOpenedPostId)
        setIsPostDirty(false)
      }
      setLastOpenedPostId(null)

      setTimeout(() => {
        window.scrollTo(0, savedScrollPosition)
      }, 0)
    }
  }, [selectedPostId, lastOpenedPostId, isPostDirty, savedScrollPosition])

  useEffect(() => {
    const handlePopState = () => {
      const path = window.location.pathname
      if (path === basePath || path === basePath.replace(/\/$/, "")) {
        setSelectedPostId(null)
      } else if (path.startsWith("/posts/")) {
        setSavedScrollPosition(window.scrollY)
        setSavedSearch(window.location.search)
        const match = path.match(/\/posts\/(\d+)/)
        if (match) {
          const postNumber = parseInt(match[1], 10)
          setSelectedPostId(postNumber)
          setLastOpenedPostId(postNumber)
        }
      }
    }

    window.addEventListener("popstate", handlePopState)
    return () => window.removeEventListener("popstate", handlePopState)
  }, [basePath, savedScrollPosition])

  return {
    selectedPostId,
    handlePostClick,
    handleCloseOverlay,
    setIsPostDirty,
  }
}
