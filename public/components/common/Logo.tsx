import * as React from "react";
import { Fider } from "@fider/services";

interface LogoProps {
  url?: string;
  size?: 50 | 100 | 200;
}

export const Logo = (props: LogoProps) => {
  const tenant = Fider.session.tenant;
  if (props.url) {
    return <img src={props.url} />;
  }

  const size = props.size || 200;
  if (tenant && tenant.logoId > 0) {
    return <img src={`${Fider.settings.assetsURL}/logo/${size}/${tenant.logoId}`} alt={tenant.name} />;
  }

  return null;
};
