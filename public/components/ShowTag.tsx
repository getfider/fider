import * as React from 'react';
import { Tag } from '@fider/models';

interface TagProps {
  tag: Tag;
  size?: 'mini' | 'tiny' | 'small' | 'normal' | 'large';
  circular?: boolean;
}

const getRGB = (color: string) => {
  const r = color.substring(1, 3);
  const g = color.substring(3, 5);
  const b = color.substring(5, 7);

  return {
    R: parseInt(r, 16),
    G: parseInt(g, 16),
    B: parseInt(b, 16)
  };
};

const idealTextColor = (color: string) => {
  const nThreshold = 105;
  const components = getRGB(color);
  const bgDelta = (components.R * 0.299) + (components.G * 0.587) + (components.B * 0.114);

  return (bgDelta > 186) ? '#000000' : '#ffffff';
};

export const ShowTag = (props: TagProps) => {
  const formatClass = props.circular === true ? 'empty circular' : '';
  const sizeClass = props.size || 'normal';

  return (
    <div
      title={`${props.tag.name}${!props.tag.isPublic ? ' (Private)' : ''}`}
      className={`ui label ${sizeClass} ${formatClass}`}
      style={{
        backgroundColor: `#${props.tag.color}`,
        color: idealTextColor(props.tag.color)
      }}
    >
      {!props.tag.isPublic && !props.circular && <i className="lock icon"/>}
      {props.circular ? '' : props.tag.name || 'Tag'}
    </div>
  );
};
