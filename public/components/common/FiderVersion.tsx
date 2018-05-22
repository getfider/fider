import * as React from "react";
import { page } from "@fider/services";

export const FiderVersion = () => {
  return (
    <p className="info center hidden-sm hidden-md">
      Support our{" "}
      <a target="_blank" href="http://opencollective.com/fider">
        OpenCollective
      </a>
      <br />
      Fider v{page.systemSettings().version}
    </p>
  );
};
