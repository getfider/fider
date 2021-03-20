import React from 'react';
import { renderToString } from 'react-dom/server';
import Home from "@fider/pages/Home/Home.page"
import { Fider, FiderContext } from './services/fider';
import { IconContext } from 'react-icons';

export function doWork(f: any, props: any) {
  let fider = Fider.initialize({...f });

  return {
    html: renderToString(
      <FiderContext.Provider value={fider}>
        <IconContext.Provider value={{ className: "icon" }}>
          <Home {...props} />
        </IconContext.Provider>
      </FiderContext.Provider>
    ),
  };
}