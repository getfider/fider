import * as React from "react";

interface LogoProps {
  url?: string;
  size?: 50 | 100 | 200;
}

export const Logo = (props: LogoProps) => {
  if (props.url) {
    return <img src={props.url} />;
  }

  const size = props.size || 200;
  if (page.tenant && page.tenant.logoId > 0) {
    return (
      <img src={`${page.settings.tenantAssetsBaseURL}/logo/${size}/${page.tenant.logoId}`} alt={page.tenant.name} />
    );
  }

  return null;
};
