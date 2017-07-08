import * as React from 'react';
import * as md from 'markdown-it';

const cm = md('commonmark', { html: false, linkify: true }).enable('linkify');
const linkify = md('zero', { linkify: true }).enable('linkify');

interface MultiLineText {
  text?: string;
  markdown?: boolean;
}

export const MultiLineText = (props: MultiLineText) => {
  if (!props.text) {
    return <p></p>;
  }

  if (props.markdown) {
    return <div dangerouslySetInnerHTML={{__html: cm.render(props.text)}}></div>;
  }

  return <div>{props.text.split('\n').map((item, i) =>
    <span dangerouslySetInnerHTML={{__html: linkify.renderInline(item) + '<br />'}}></span>
  )}</div>;
};
