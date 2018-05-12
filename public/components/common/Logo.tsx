import * as React from "react";
import { Tenant } from "@fider/models";

interface LogoProps {
  url?: string;
  tenant?: Tenant;
}

export const Logo = (props: LogoProps) => {
  if (props.url) {
    return <img src={props.url} />;
  }
  if (props.tenant && props.tenant.logoId > 0) {
    return <img src={`/logo/50/${props.tenant.logoId}`} alt={props.tenant.name} />;
  }
  return null;
};
