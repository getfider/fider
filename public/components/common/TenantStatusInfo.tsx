import React from "react";
import { Fider } from "@fider/services";
import { TenantStatus } from "@fider/models";
import { Message } from "./Message";

export const TenantStatusInfo = () => {
  if (!Fider.isBillingEnabled() || Fider.session.tenant.status !== TenantStatus.Locked) {
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
