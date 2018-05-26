import * as React from "react";

import { CurrentUser, UserSettings } from "@fider/models";
import { Toggle, Segment, Segments, Field, Button, Modal } from "@fider/components";

interface DangerZoneProps {
  user: CurrentUser;
}

interface DangerZoneState {
  clicked: boolean;
}

export class DangerZone extends React.Component<DangerZoneProps, DangerZoneState> {
  constructor(props: DangerZoneProps) {
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

  public render() {
    return (
      <div className="l-danger-zone">
        <Modal.Window canClose={true} isOpen={this.state.clicked} center={false} onClose={this.onCancel}>
          <Modal.Header>Delete account</Modal.Header>
          <Modal.Content>
            <p>
              When you choose to delete your account, we will erase all your personal information forever. The content
              you have published, votes and comments will remain, but they will be anonymised.
            </p>
            <p>
              This process is irreversible. <strong>Are you sure?</strong>
            </p>
          </Modal.Content>
          <Modal.Footer>
            <Button color="danger" size="tiny">
              Confirm
            </Button>
            <Button size="tiny" onClick={this.onCancel}>
              Cancel
            </Button>
          </Modal.Footer>
        </Modal.Window>
        <h4>Delete account</h4>
        <p className="info">
          When you choose to delete your account, we will erase all your personal information forever. The content you
          have published, votes and comments will remain, but they will be anonymised.
        </p>
        <p className="info">This process is irreversible. Please be certain.</p>
        <Button color="danger" size="tiny" onClick={this.onClickDelete}>
          Delete My Account
        </Button>
      </div>
    );
  }
}
