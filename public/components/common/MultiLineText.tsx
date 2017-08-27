import * as React from 'react';
import * as md from 'markdown-it';

const simple = md('commonmark', { html: false, breaks: true, linkify: true }).disable('heading').disable('image');
const full = md('commonmark', { html: false, breaks: true, linkify: true });

interface MultiLineText {
  className?: string;
  text?: string;
  style: 'full' | 'simple';
}

export const MultiLineText = (props: MultiLineText) => {
  if (!props.text) {
    return <p></p>;
  }

  const func = props.style === 'full' ? full : simple;
  return <div className={ `markdown-body ${props.className || ''}` } dangerouslySetInnerHTML={{__html: func.render(props.text)}}></div>;
};
