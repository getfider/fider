import * as React from "react";
import { Fider, uploadedImageUrl } from "@fider/services";

type Size = 24 | 50 | 100 | 200;

interface LogoProps {
  size: Size;
}

export const TenantLogoUrl = (size: Size): string | undefined => {
  const tenant = Fider.session.tenant;
  if (tenant && tenant.logoId > 0) {
    return uploadedImageUrl(tenant.logoId, size);
  }
  return undefined;
};

export const TenantLogo = (props: LogoProps) => {
  const tenant = Fider.session.tenant;
  if (tenant && tenant.logoId > 0) {
    return <img src={TenantLogoUrl(props.size)} alt={tenant.name} />;
  }
  return null;
};
