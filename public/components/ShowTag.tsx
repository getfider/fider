import * as React from 'react';

interface TagProps {
  name: string;
  color: string;
  isPublic: boolean;
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

  return ((components.R * 0.299) + (components.G * 0.587) + (components.B * 0.114) > 186) ? '#000000' : '#ffffff';
};

export const ShowTag = (props: TagProps) => {
  return (
    <div
      title={`${props.name} ${!props.isPublic && '(Private)'}`}
      className="ui label"
      style={{
        backgroundColor: `#${props.color}`,
        color: idealTextColor(props.color)
      }}
    >
      {!props.isPublic && <i className="lock icon"/>}
      {props.name || 'Tag'}
    </div>
  );
};
