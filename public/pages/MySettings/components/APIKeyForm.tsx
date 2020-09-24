import React from "react";
import { Button } from "@fider/components";
import { actions } from "@fider/services";
import { useTranslation } from "react-i18next";

interface APIKeyFormState {
  apiKey?: string;
}

export class APIKeyForm extends React.Component<{}, APIKeyFormState> {
  constructor(props: {}) {
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
    const { t } = useTranslation();
    return (
      <>
        <p className="info" dangerouslySetInnerHTML={
          { __html: t('mySettings.newApiKey', { apiKey: this.state.apiKey }) }
        }>
        </p>
        <p className="info"></p>
      </>
    );
  }

  public render() {
    const { t } = useTranslation();
    return (
      <div className="l-api-key">
        <h4>API Key</h4>
        <p className="info">
          {t('mySettings.apiKeyStoreMessage')}
        </p>
        <p className="info" dangerouslySetInnerHTML={
          { __html: t('mySettings.apiKeyDoc') }
        }>
        </p>
        <p>
          <Button size="tiny" onClick={this.regenerate}>
            {t('mySettings.generateAPIKey')}
          </Button>
        </p>
        {this.state.apiKey && this.showAPIKey()}
      </div>
    );
  }
}
