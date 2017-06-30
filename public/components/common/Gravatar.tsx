import * as React from 'react';
import md5 = require('md5');

export const Gravatar = (props: {email?: string}) => {
  const hash = props.email ? md5(props.email) : '';
  return <img className="ui avatar image" src={ 'https://www.gravatar.com/avatar/' + hash } />;
};