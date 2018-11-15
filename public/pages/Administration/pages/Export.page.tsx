import React from "react";

import { Button, Form, Field } from "@fider/components";
import { FaRegFileExcel } from "react-icons/fa";
import { AdminBasePage } from "../components/AdminBasePage";

export default class ExportPage extends AdminBasePage<{}, {}> {
  public id = "p-admin-export";
  public name = "export";
  public icon = FaRegFileExcel;
  public title = "Export";
  public subtitle = "Download your data";

  public content() {
    return (
      <Form>
        <Field label="Export Posts">
          <p className="info">
            Use this button to download a CSV file with all posts in this site. This can be useful to analyse the data
            with an external tool or simply to back it up.
          </p>
        </Field>
        <Field>
          <Button color="positive" href="/admin/export/posts.csv">
            Download
          </Button>
        </Field>
      </Form>
    );
  }
}
