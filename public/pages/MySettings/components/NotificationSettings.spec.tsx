import React from "react"
import { render } from "@testing-library/react"
import { Fider, FiderContext } from "@fider/services"
import { NotificationSettings } from "./NotificationSettings"

const session = (tenant: { isPro: boolean; isModerationEnabled: boolean }, user: { isCollaborator: boolean; role: string }) =>
  Fider.initialize({
    settings: { environment: "development", oauth: [] },
    tenant: { allowedSchemes: "", ...tenant },
    user: { name: "Jon Snow", ...user },
  })

const admin = { isCollaborator: true, role: "administrator" }
const collaborator = { isCollaborator: true, role: "collaborator" }
const visitor = { isCollaborator: false, role: "visitor" }
const moderated = { isPro: true, isModerationEnabled: true }

const showsAwaitingApproval = (fider: ReturnType<typeof Fider.initialize>) => {
  render(
    <FiderContext.Provider value={fider}>
      <NotificationSettings userSettings={{}} settingsChanged={() => undefined} />
    </FiderContext.Provider>
  )
  // Translations are mocked out in tests, so query by structure rather than text.
  return document.querySelector(".c-notification-settings__awaiting-approval") !== null
}

describe("<NotificationSettings /> Awaiting Approval", () => {
  test("shows for administrators on a moderated Pro site", () => {
    expect(showsAwaitingApproval(session(moderated, admin))).toBe(true)
  })

  test("shows for collaborators on a moderated Pro site", () => {
    expect(showsAwaitingApproval(session(moderated, collaborator))).toBe(true)
  })

  test("hidden from visitors", () => {
    expect(showsAwaitingApproval(session(moderated, visitor))).toBe(false)
  })

  test("hidden when content moderation is off", () => {
    expect(showsAwaitingApproval(session({ isPro: true, isModerationEnabled: false }, admin))).toBe(false)
  })

  test("hidden when the site isn't Pro", () => {
    expect(showsAwaitingApproval(session({ isPro: false, isModerationEnabled: true }, admin))).toBe(false)
  })
})
