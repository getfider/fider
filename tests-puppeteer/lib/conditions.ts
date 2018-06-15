import { WebComponent } from "./components";
import { NewablePage, Page, Browser } from ".";

export interface WaitCondition {
  function: (...args: any[]) => boolean;
  args?: any[];
}

export const elementIsVisible = (target: string | WebComponent): WaitCondition => {
  const selector = typeof target === "string" ? target : target.selector;
  return {
    function: (query: string) => {
      const node = document.querySelector(query);
      if (!node) {
        return false;
      }
      const style = window.getComputedStyle(node);
      return style && style.display !== "none" && style.visibility !== "hidden" && style.opacity !== "0";
    },
    args: [selector]
  };
};

// export function pageHasLoaded<T extends Page>(page: NewablePage<T>): WaitCondition {
//   return (browser: Browser) => {
//     const condition = new page(browser).loadCondition();
//     return condition(browser);
//   };
// }