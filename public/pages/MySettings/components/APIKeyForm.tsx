import React from "react";
import { Button } from "@fider/components";
import { actions } from "@fider/services";
import { withTranslation, WithTranslation } from "react-i18next";

interface APIKeyFormState {
  apiKey?: string;
}
class InternalAPIKeyForm extends React.Component<WithTranslation, APIKeyFormState> {
  constructor(props: WithTranslation) {
    super(props);
    this.state = {};
  }

  private regenerate = async () => {
    const result = await actions.regenerateAPIKey();
    if (result.ok) {
      this.setState({ apiKey: result.data.apiKey });
    }
  };

  private showAPIKey() {
    const { t } = this.props;
    return (
      <>
        <p className="info">
          {t("mySettings.newApiKey")}
          <code>{this.state.apiKey}</code>
        </p>
        <p className="info" />
      </>
    );
  }

  public render() {
    const { t } = this.props;
    return (
      <div className="l-api-key">
        <h4>API Key</h4>
        <p className="info">{t("mySettings.apiKeyStoreMessage")}</p>
        <p className="info" dangerouslySetInnerHTML={{ __html: t("mySettings.apiKeyDoc") }} />
        <p>
          <Button size="tiny" onClick={this.regenerate}>
            {t("mySettings.generateAPIKey")}
          </Button>
        </p>
        {this.state.apiKey && this.showAPIKey()}
      </div>
    );
  }
}

export const APIKeyForm = withTranslation()(InternalAPIKeyForm);
