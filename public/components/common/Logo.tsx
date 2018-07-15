import * as React from "react";
import { Fider } from "@fider/services";

interface LogoProps {
  size: 50 | 100 | 200;
}

export const LogoUrl = (size: 50 | 100 | 200): string | undefined => {
  const tenant = Fider.session.tenant;
  if (tenant && tenant.logoId > 0) {
    return `${Fider.settings.assetsURL}/images/${size}/${Fider.session.tenant.logoId}`;
  }
  return undefined;
};

export const Logo = (props: LogoProps) => {
  const tenant = Fider.session.tenant;
  if (tenant && tenant.logoId > 0) {
    return <img src={LogoUrl(props.size)} alt={tenant.name} />;
  }
  return null;
};
