import * as React from "react";

import { Button, Form, Field } from "@fider/components";
import { AdminBasePage } from "../components";

export class AuthenticationPage extends AdminBasePage<{}, {}> {
  public id = "p-admin-authentication";
  public name = "authentication";
  public icon = "sign in alternate";
  public title = "Authentication";
  public subtitle = "Manager your site authentication";

  public content() {
    return (
      <>
        <h3>Hello World</h3>
      </>
    );
  }
}
