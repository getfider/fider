import React from "react";
import { TenantStatus } from "@fider/models";
import { Message } from "./Message";
import { useFider } from "@fider/hooks";

export const TenantStatusInfo = () => {
  const fider = useFider();

  if (!fider.isBillingEnabled() || fider.session.tenant.status !== TenantStatus.Locked) {
    return null;
  }

  return (
    <div className="container">
      <Message type="error">
        This site is locked due to lack of a subscription. Visit the <a href="/admin/billing">Billing</a> settings to
        update it.
      </Message>
    </div>
  );
};
