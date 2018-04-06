import * as React from "react";

import { HomePage, HomePageProps } from "../";
import { Modal, Form, Button } from "@fider/components/common";
import { page, actions } from "@fider/services";

interface CompleteSignInProfilePageState {
  name: string;
}

export class CompleteSignInProfilePage extends React.Component<HomePageProps, CompleteSignInProfilePageState> {
  private form!: Form;
  private key: string;

  constructor(props: HomePageProps) {
    super(props);
    this.key = page.getQueryString("k");
    this.state = {
      name: ""
    };
  }

  private async submit() {
    const result = await actions.completeProfile(this.key, this.state.name);
    if (result.ok) {
      location.href = "/";
    } else if (result.error) {
      this.form.setFailure(result.error);
    }
  }

  public render() {
    return (
      <>
        <Modal.Window canClose={false} isOpen={true}>
          <Modal.Header>Complete your profile</Modal.Header>
          <Modal.Content>
            <p>Because this is your first sign in, please input your display name.</p>
            <Form
              ref={f => {
                this.form = f!;
              }}
            >
              <div className="ui small action fluid input">
                <input
                  onChange={e => this.setState({ name: e.currentTarget.value })}
                  type="text"
                  maxLength={100}
                  placeholder="Your display name"
                  className="small"
                />
                <Button onClick={() => this.submit()} size="small" color="green" disabled={this.state.name === ""}>
                  Submit
                </Button>
              </div>
            </Form>
          </Modal.Content>
        </Modal.Window>
        {React.createElement(HomePage, this.props)}
      </>
    );
  }
}
