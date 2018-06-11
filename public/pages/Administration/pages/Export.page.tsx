import "./PrivacySettings.page.scss";

import * as React from "react";

import { CurrentUser } from "@fider/models";
import { Button, Form, Field } from "@fider/components";
import { AdminBasePage } from "../components";

export class ExportPage extends AdminBasePage<{}, {}> {
  public id = "p-admin-export";
  public name = "export";
  public icon = "file excel outline";
  public title = "Export";
  public subtitle = "Download your data";

  public content() {
    return (
      <Form>
        <Field label="Export Ideas">
          <p className="info">
            Use this button to download a CSV file with all ideas in this site. This can be useful to analyse the data
            with an external tool or simply to back it up.
          </p>
        </Field>
        <Field>
          <Button href="/admin/export/ideas.csv">Download</Button>
        </Field>
      </Form>
    );
  }
}
