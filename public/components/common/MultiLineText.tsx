import * as React from 'react';
import * as md from 'markdown-it';

const cm = md('commonmark', { html: false, linkify: true }).enable('linkify');

interface MultiLineText {
  className?: string;
  text?: string;
  markdown?: boolean;
}

export const MultiLineText = (props: MultiLineText) => {
  if (!props.text) {
    return <p></p>;
  }

  if (props.markdown) {
    return <div className={props.className} dangerouslySetInnerHTML={{__html: cm.render(props.text)}}></div>;
  }

  return <div className={props.className}>{props.text.split('\n').map((item, i) =>
    <span dangerouslySetInnerHTML={{__html: cm.renderInline(item) + '<br />'}}></span>
  )}</div>;
};
