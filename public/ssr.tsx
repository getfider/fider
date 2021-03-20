import React from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { Fider, FiderContext } from "./services/fider";
import { IconContext } from "react-icons";
import { Footer, Header } from "./components";
import { resolveRootComponent } from "./router";

function ssrRender(url: string, pathname: string, args: any) {
  const fider = Fider.initialize({ ...args });
  const config = resolveRootComponent(pathname);
  window.location.href = url;

  return renderToStaticMarkup(
    <FiderContext.Provider value={fider}>
      <IconContext.Provider value={{ className: "icon" }}>
        {config.showHeader && <Header />}
        {React.createElement(config.component, args.props)}
        {config.showHeader && <Footer />}
      </IconContext.Provider>
    </FiderContext.Provider>
  );
}

(globalThis as any).ssrRender = ssrRender;
