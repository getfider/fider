import "./PrivacySettings.page.scss";

import * as React from "react";

import { CurrentUser } from "@fider/models";
import { Button } from "@fider/components/common";
import { AdminBasePage } from "../components";

interface ExportPageProps {
  user: CurrentUser;
}

export class ExportPage extends AdminBasePage<ExportPageProps, {}> {
  public id = "p-admin-export";
  public name = "export";
  public icon = "file excel outline";
  public title = "Export";
  public subtitle = "Download your data";

  public content() {
    return (
      <div className="ui form">
        <div className="field">
          <label htmlFor="private">Export Ideas</label>
          <p className="info">
            Use this button to download a CSV file with all ideas in this site. This can be useful to analyse the data
            with an external tool or simply to back it up.
          </p>
          <Button href="/admin/export/ideas.csv">Download</Button>
        </div>
      </div>
    );
  }
}
