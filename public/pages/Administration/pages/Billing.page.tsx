import React from "react";

import { FaFileInvoice } from "react-icons/fa";
import { AdminBasePage } from "../components/AdminBasePage";

export default class BillingPage extends AdminBasePage<{}, {}> {
  public id = "p-admin-billing";
  public name = "billing";
  public icon = FaFileInvoice;
  public title = "Billing";
  public subtitle = "Manage your subscription";

  public content() {
    return <h4>Work in progress</h4>;
  }
}
