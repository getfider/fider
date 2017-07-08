import * as React from 'react';
import md5 = require('md5');

interface GravatarProps {
  hash?: string;
  email?: string;
}

export const Gravatar = (props: GravatarProps) => {
  const hash = props.email ? md5(props.email) : props.hash || '';
  return <img className="ui avatar image" src={ 'https://www.gravatar.com/avatar/' + hash } />;
};
