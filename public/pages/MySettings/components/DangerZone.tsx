import React from "react";

import { Button, Modal, ButtonClickEvent } from "@fider/components";
import { actions, notify, navigator } from "@fider/services";
import { useTranslation } from "react-i18next";

interface DangerZoneState {
  clicked: boolean;
}

export class DangerZone extends React.Component<{}, DangerZoneState> {
  constructor(props: {}) {
    super(props);
    this.state = {
      clicked: false
    };
  }

  public onClickDelete = async () => {
    this.setState({ clicked: true });
  };

  public onCancel = async () => {
    this.setState({ clicked: false });
  };

  public onConfirm = async (e: ButtonClickEvent) => {
    const { t } = useTranslation();
    const response = await actions.deleteCurrentAccount();
    if (response.ok) {
      e.preventEnable();
      navigator.goHome();
    } else {
      notify.error(t("mySettings.deleteAccountFailed"));
    }
  };

  public render() {
    const { t } = useTranslation();
    return (
      <div className="l-danger-zone">
        <Modal.Window isOpen={this.state.clicked} center={false} onClose={this.onCancel}>
          <Modal.Header>{t('mySettings.deleteAccount')}</Modal.Header>
          <Modal.Content>
            <p>
              When you choose to delete your account, we will erase all your personal information forever. The content
              you have published will remain, but it will be anonymised.
            </p>
            <p>
              This process is irreversible. <strong>Are you sure?</strong>
            </p>
          </Modal.Content>
          <Modal.Footer>
            <Button color="danger" size="tiny" onClick={this.onConfirm}>
              {t('common.button.confirm')}
            </Button>
            <Button color="cancel" size="tiny" onClick={this.onCancel}>
              {t('common.button.cancel')}
            </Button>
          </Modal.Footer>
        </Modal.Window>

        <h4>{t('mySettings.deleteAccount')}</h4>
        <p className="info">

        </p>
        <p className="info"> {t('mySettings.deleteAccountMessage.part2')}
        </p>
        <Button color="danger" size="tiny" onClick={this.onClickDelete}>
          {t('mySettings.deleteMyAccount')}
        </Button>
      </div>
    );
  }
}
