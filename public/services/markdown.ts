import * as md from "markdown-it";

const fullMarkdown = md("commonmark", {
  html: false,
  breaks: true,
  linkify: true
}).enable(["linkify", "strikethrough"]);

const simpleMarkdown = md("commonmark", { html: false, breaks: true, linkify: true })
  .enable(["linkify", "strikethrough"])
  .disable(["heading", "image"]);

const linkOpen = (tokens: md.Token[], idx: number, options: any, env: any, self: md.Renderer) => {
  const aIndex = tokens[idx].attrIndex("target");

  if (aIndex < 0) {
    tokens[idx].attrPush(["target", "_blank"]);
  } else {
    tokens[idx].attrs[aIndex][1] = "_blank";
  }

  return self.renderToken(tokens, idx, options);
};

fullMarkdown.renderer.rules.link_open = linkOpen;
simpleMarkdown.renderer.rules.link_open = linkOpen;

export const full = (input: string): string => {
  return fullMarkdown.render(input).trim();
};

export const simple = (input: string): string => {
  return simpleMarkdown.render(input).trim();
};
