import * as React from "react";
import { Tenant } from "@fider/models";
import { page } from "@fider/services";

interface LogoProps {
  url?: string;
  tenant?: Tenant;
  size?: 50 | 100 | 200;
}

export const Logo = (props: LogoProps) => {
  if (props.url) {
    return <img src={props.url} />;
  }

  const size = props.size || 200;
  if (props.tenant && props.tenant.logoId > 0) {
    return (
      <img
        src={`${window.props.settings.tenantAssetsBaseURL}/logo/${size}/${props.tenant.logoId}`}
        alt={props.tenant.name}
      />
    );
  }

  return null;
};
