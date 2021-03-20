import React from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import Home from "@fider/pages/Home/Home.page"
import { Fider, FiderContext } from './services/fider';
import { IconContext } from 'react-icons';

function doWork(f: any, props: any) {
  let fider = Fider.initialize({...f });

  return renderToStaticMarkup(
    <FiderContext.Provider value={fider}>
      <IconContext.Provider value={{ className: "icon" }}>
        <Home {...props} />
      </IconContext.Provider>
    </FiderContext.Provider>
  )
}

(globalThis as any).doWork = doWork