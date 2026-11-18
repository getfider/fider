import React from "react"
import { screen, fireEvent, render } from "@testing-library/react"
import { setupModalRoot } from "@fider/services/testing"
import { Fider, FiderContext } from "@fider/services"
import { ShareFeedback } from "./ShareFeedback"

const signedIn = Fider.initialize({
  settings: { environment: "development", oauth: [] },
  tenant: { allowedSchemes: "" },
  user: { name: "Jon Snow" },
})

beforeEach(() => {
  setupModalRoot()
  window.localStorage.clear()
})

const renderWithTitle = (title: string, fider = signedIn) => {
  render(
    <FiderContext.Provider value={fider}>
      <ShareFeedback isOpen={true} placeholder="" onClose={() => undefined} tags={[]} />
    </FiderContext.Provider>
  )
  fireEvent.change(screen.getByRole("textbox", { name: /title/i }), { target: { value: title } })
}

// Translations are mocked out in tests, so query by structure rather than text.
const submitButton = () => document.querySelector(".c-share-feedback-signin button")
const hint = () => document.querySelector(".c-share-feedback__title-hint")

// The server counts UTF-8 bytes (minimum 10), so 4 CJK characters (12 bytes) are enough.
describe("<ShareFeedback /> submit button", () => {
  test.each([["abcdefghij"], ["深色模式"], ["ダークモード"]])("is enabled for %s", (title) => {
    renderWithTitle(title)
    expect(submitButton()).toBeEnabled()
    expect(hint()).toBeNull()
  })

  test.each([[""], ["abcdefghi"], ["農民曆"], ["  abc   def  "]])("is disabled with a hint for %p", (title) => {
    renderWithTitle(title)
    expect(submitButton()).toBeDisabled()
    expect(hint()).not.toBeNull()
  })
})

describe("<ShareFeedback /> signed out", () => {
  const signedOut = () => Fider.initialize({ settings: { environment: "development", oauth: [] }, tenant: { allowedSchemes: "", isEmailAuthAllowed: true } })
  const signIn = () => document.querySelector(".c-share-feedback-signin")

  test("sign-in is locked while the title is too short", () => {
    renderWithTitle("深色", signedOut())
    expect(signIn()).toHaveAttribute("inert")
    expect(hint()).not.toBeNull()
  })

  test("sign-in unlocks once the title is long enough", () => {
    renderWithTitle("深色模式", signedOut())
    expect(signIn()).not.toHaveAttribute("inert")
    expect(hint()).toBeNull()
  })
})
