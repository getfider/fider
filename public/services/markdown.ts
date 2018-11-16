import marked from "marked";

marked.setOptions({
  headerIds: false,
  xhtml: true,
  smartLists: true,
  gfm: true,
  breaks: true
});

const link = (href: string, title: string, text: string) => {
  const titleAttr = title ? ` title=${title}` : "";
  return `<a href="${href}"${titleAttr} target="_blank">${text}</a>`;
};

const simpleRenderer = new marked.Renderer();
simpleRenderer.heading = (text, level, raw) => `<p>${raw}</p>`;
simpleRenderer.image = (href, title, text) => "";
simpleRenderer.link = link;

const fullRenderer = new marked.Renderer();
fullRenderer.link = link;

export const full = (input: string): string => {
  return marked(input, { renderer: fullRenderer }).trim();
};

export const simple = (input: string): string => {
  return marked(input, { renderer: simpleRenderer }).trim();
};
