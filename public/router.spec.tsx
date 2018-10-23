import { resolveRootComponent } from "./router";
import * as Loadable from "react-loadable";

[{ path: "", shouldError: false }, { path: "thispathshouldn'tmatchanything", shouldError: true }].forEach(x => {
  test(`Router should resolve correct component for path '${x.path}'`, () => {
    if (x.shouldError) {
      expect(() => resolveRootComponent(x.path)).toThrowError();
    } else {
      const page = resolveRootComponent(x.path);
      expect(page.component).toBeTruthy();
    }
  });
});
