import { resolveRootComponent } from "./router";
import { HomePage, ShowPostPage, InvitationsPage, GeneralSettingsPage, OAuthEchoPage } from "@fider/pages";

[
  { path: "", expected: HomePage },
  { path: "/posts/123", expected: ShowPostPage },
  { path: "/posts/123/the-slug", expected: ShowPostPage },
  { path: "/posts" },
  { path: "/admin", expected: GeneralSettingsPage },
  { path: "/admin/invitations", expected: InvitationsPage },
  { path: "/oauth/_name/echo", expected: OAuthEchoPage }
].forEach(x => {
  test(`Router should resolve correct component for path '${x.path}'`, () => {
    if (x.expected) {
      const page = resolveRootComponent(x.path);
      expect(page.component).toEqual(x.expected);
    } else {
      expect(() => resolveRootComponent(x.path)).toThrowError();
    }
  });
});
