import React from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { Fider, FiderContext } from "./services/fider";
import { IconContext } from "react-icons";
import { Footer, Header } from "./components";
import HomePage from "./pages/Home/Home.page";
import ShowPostPage from "./pages/ShowPost/ShowPost.page";
import SignInPage from "./pages/SignIn/SignIn.page";
import SignUpPage from "./pages/SignUp/SignUp.page";
import UIToolkitPage from "./pages/UI/UIToolkit.page";

interface PageConfiguration {
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

const route = (path: string, component: any, showHeader: boolean = true): PageConfiguration => {
  path = path.replace("/", "/").replace(":number", "\\d+").replace(":string", ".+").replace("*", "/?.*");

  const regex = new RegExp(`^${path}$`);
  return { regex, component, showHeader };
};

const pathRegex = [route("", HomePage), route("/posts/:number*", ShowPostPage), route("/signin", SignInPage, false), route("/signup", SignUpPage, false), route("/-/ui", UIToolkitPage)];

export const resolveRootComponent = (path: string): PageConfiguration => {
  if (path.length > 0 && path.charAt(path.length - 1) === "/") {
    path = path.substring(0, path.length - 1);
  }
  for (const entry of pathRegex) {
    if (entry && entry.regex.test(path)) {
      return entry;
    }
  }
  throw new Error(`Component not found for route ${path}.`);
};

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
