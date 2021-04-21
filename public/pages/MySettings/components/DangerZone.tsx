import React from "react"

import { Button, Modal, ButtonClickEvent } from "@fider/components"
import { actions, notify, navigator } from "@fider/services"

interface DangerZoneState {
  clicked: boolean
}

export class DangerZone extends React.Component<any, DangerZoneState> {
  constructor(props: any) {
    super(props)
    this.state = {
      clicked: false,
    }
  }

  public onClickDelete = async () => {
    this.setState({ clicked: true })
  }

  public onCancel = async () => {
    this.setState({ clicked: false })
  }

  public onConfirm = async (e: ButtonClickEvent) => {
    const response = await actions.deleteCurrentAccount()
    if (response.ok) {
      e.preventEnable()
      navigator.goHome()
    } else {
      notify.error("Failed to delete your account. Try again later")
    }
  }

  public render() {
    return (
      <div className="l-danger-zone">
        <Modal.Window isOpen={this.state.clicked} center={false} onClose={this.onCancel}>
          <Modal.Header>Delete account</Modal.Header>
          <Modal.Content>
            <p>
              When you choose to delete your account, we will erase all your personal information forever. The content you have published will remain, but it
              will be anonymised.
            </p>
            <p>
              This process is irreversible. <strong>Are you sure?</strong>
            </p>
          </Modal.Content>
          <Modal.Footer>
            <Button color="danger" size="tiny" onClick={this.onConfirm}>
              Confirm
            </Button>
            <Button color="cancel" size="tiny" onClick={this.onCancel}>
              Cancel
            </Button>
          </Modal.Footer>
        </Modal.Window>

        <h4>Delete account</h4>
        <p className="info">
          When you choose to delete your account, we will erase all your personal information forever. The content you have published will remain, but it will
          be anonymised.
        </p>
        <p className="info">This process is irreversible. Please be certain.</p>
        <Button color="danger" size="tiny" onClick={this.onClickDelete}>
          Delete My Account
        </Button>
      </div>
    )
  }
}
