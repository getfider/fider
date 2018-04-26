import { resolveRootComponent } from "./router";
import { HomePage } from "@fider/pages";

test("Router should resolve correct component", () => {
  const page = resolveRootComponent("");
  expect(page.component).toEqual(HomePage);
});
