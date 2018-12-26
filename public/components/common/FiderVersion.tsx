import React from "react";
import { Fider } from "@fider/services";

export const FiderVersion = () => {
  return (
    <p className="info center hidden-sm hidden-md">
      {!Fider.isBillingEnabled() && (
        <>
          Support our{" "}
          <a target="_blank" href="http://opencollective.com/fider">
            OpenCollective
          </a>
          <br />
        </>
      )}
      Fider v{Fider.settings.version}
    </p>
  );
};
