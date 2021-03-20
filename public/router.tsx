interface PageConfiguration {
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

export const route = (path: string, component: any, showHeader: boolean = true): PageConfiguration => {
  path = path.replace("/", "/").replace(":number", "\\d+").replace(":string", ".+").replace("*", "/?.*");

  const regex = new RegExp(`^${path}$`);
  return { regex, component, showHeader };
};

export const resolveRootComponent = (path: string, routes: PageConfiguration[]): PageConfiguration => {
  if (path.length > 0 && path.charAt(path.length - 1) === "/") {
    path = path.substring(0, path.length - 1);
  }
  for (const entry of routes) {
    if (entry && entry.regex.test(path)) {
      return entry;
    }
  }
  throw new Error(`Component not found for route ${path}.`);
};
