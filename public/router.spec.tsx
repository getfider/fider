import { resolveRootComponent } from "./router";
import { HomePage, ShowIdeaPage, InvitationsPage, GeneralSettingsPage } from "@fider/pages";

[
  { path: "", expected: HomePage },
  { path: "/ideas/123", expected: ShowIdeaPage },
  { path: "/ideas/123/the-slug", expected: ShowIdeaPage },
  { path: "/ideas" },
  { path: "/admin", expected: GeneralSettingsPage },
  { path: "/admin/invitations", expected: InvitationsPage }
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
