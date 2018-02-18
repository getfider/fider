import * as React from 'react';
import * as md from 'markdown-it';

const full = md('commonmark', { html: false, breaks: true, linkify: true }).enable(['linkify', 'strikethrough']);
const simple = md('commonmark', { html: false, breaks: true, linkify: true }).enable(['linkify', 'strikethrough']).disable([ 'heading', 'image' ]);

interface MultiLineText {
  className?: string;
  text?: string;
  style: 'full' | 'simple';
}

export const MultiLineText = (props: MultiLineText) => {
  if (!props.text) {
    return <p />;
  }

  const func = props.style === 'full' ? full : simple;
  return (
    <div
      className={`markdown-body ${props.className || ''}`}
      dangerouslySetInnerHTML={{__html: func.render(props.text)}}
    />
  );
};
