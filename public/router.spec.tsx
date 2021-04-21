import { resolveRootComponent } from "./router"
import * as Pages from "@fider/AsyncPages"
;[
  { path: "", expected: Pages.AsyncHomePage },
  { path: "/posts/123", expected: Pages.AsyncShowPostPage },
  { path: "/posts/123/the-slug", expected: Pages.AsyncShowPostPage },
  { path: "/posts" },
  { path: "/admin", expected: Pages.AsyncGeneralSettingsPage },
  { path: "/admin/invitations", expected: Pages.AsyncInvitationsPage },
  { path: "/oauth/_name/echo", expected: Pages.AsyncOAuthEchoPage },
].forEach((x) => {
  test(`Router should resolve correct component for path '${x.path}'`, () => {
    if (x.expected) {
      const page = resolveRootComponent(x.path)
      expect(page.component).toEqual(x.expected)
    } else {
      expect(() => resolveRootComponent(x.path)).toThrowError()
    }
  })
})
