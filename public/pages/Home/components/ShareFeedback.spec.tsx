import React from "react"
import { screen, fireEvent, render } from "@testing-library/react"
import { setupModalRoot } from "@fider/services/testing"
import { Fider, FiderContext } from "@fider/services"
import { ShareFeedback } from "./ShareFeedback"

const fider = Fider.initialize({
  settings: { environment: "development", oauth: [] },
  tenant: { allowedSchemes: "" },
  user: { name: "Jon Snow" },
})

beforeEach(() => {
  setupModalRoot()
  window.localStorage.clear()
  window.sessionStorage.clear()
})

const renderWithTitle = (title: string) => {
  render(
    <FiderContext.Provider value={fider}>
      <ShareFeedback isOpen={true} placeholder="" onClose={() => undefined} tags={[]} />
    </FiderContext.Provider>
  )
  fireEvent.change(screen.getByRole("textbox", { name: /title/i }), { target: { value: title } })
}

// Translations are mocked out in tests, so the button has no text to query by.
const submitButton = () => document.querySelector(".c-share-feedback-signin button")

describe("<ShareFeedback /> submit button", () => {
  test.each([
    ["ASCII title of 10 bytes", "abcdefghij"],
    ["Chinese title of 4 characters (12 bytes)", "深色模式"],
    ["Japanese title of 5 characters (15 bytes)", "ダークモード"],
  ])("is shown for %s", (_, title) => {
    renderWithTitle(title)
    expect(submitButton()).not.toBeNull()
  })

  test.each([
    ["ASCII title of 9 bytes", "abcdefghi"],
    ["Chinese title of 3 characters (9 bytes)", "農民曆"],
    ["title padded with whitespace", "  abc   def  "],
  ])("is hidden for %s", (_, title) => {
    renderWithTitle(title)
    expect(submitButton()).toBeNull()
  })
})
