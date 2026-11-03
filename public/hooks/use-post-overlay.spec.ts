import { renderHook, act } from "@testing-library/react"
import { usePostOverlay } from "./use-post-overlay"

describe("usePostOverlay", () => {
  let scrollTo: jest.SpyInstance

  beforeEach(() => {
    window.history.replaceState(null, "", "/?view=trending")
    Object.defineProperty(window, "scrollY", { value: 1200, configurable: true })
    scrollTo = jest.spyOn(window, "scrollTo").mockImplementation(() => undefined)
  })

  afterEach(() => {
    scrollTo.mockRestore()
  })

  test("opening a post scrolls to the top and closing restores the list position synchronously", () => {
    const { result } = renderHook(() => usePostOverlay({ basePath: "/" }))

    act(() => result.current.handlePostClick(5, "my-post"))
    expect(result.current.selectedPostId).toBe(5)
    expect(window.location.pathname).toBe("/posts/5/my-post")
    expect(scrollTo).toHaveBeenLastCalledWith(0, 0)

    act(() => result.current.handleCloseOverlay())
    expect(result.current.selectedPostId).toBeNull()
    expect(window.location.pathname + window.location.search).toBe("/?view=trending")
    expect(scrollTo).toHaveBeenLastCalledWith(0, 1200)
  })

  test("browser back restores the scroll position saved on the list's history entry", () => {
    const { result } = renderHook(() => usePostOverlay({ basePath: "/" }))

    act(() => result.current.handlePostClick(5, "my-post"))

    act(() => {
      window.history.replaceState(null, "", "/")
      window.dispatchEvent(new PopStateEvent("popstate", { state: { scrollY: 800 } }))
    })
    expect(result.current.selectedPostId).toBeNull()
    expect(scrollTo).toHaveBeenLastCalledWith(0, 800)
  })

  test("onPostClosed is only called when the post changed", () => {
    const onPostClosed = jest.fn()
    const { result } = renderHook(() => usePostOverlay({ basePath: "/", onPostClosed }))

    act(() => result.current.handlePostClick(5, "my-post"))
    act(() => result.current.handleCloseOverlay())
    expect(onPostClosed).not.toHaveBeenCalled()

    act(() => result.current.handlePostClick(6, "other-post"))
    act(() => result.current.setIsPostDirty(true))
    act(() => result.current.handleCloseOverlay())
    expect(onPostClosed).toHaveBeenCalledWith(6)
  })
})
