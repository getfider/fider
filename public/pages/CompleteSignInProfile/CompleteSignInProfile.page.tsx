import * as React from "react";

import { HomePage, HomePageProps, SignInPage } from "../";
import { Modal, Button, Form2, Input } from "@fider/components";
import { page, actions, Failure } from "@fider/services";

interface CompleteSignInProfilePageState {
  name: string;
  error?: Failure;
}

export class CompleteSignInProfilePage extends React.Component<HomePageProps, CompleteSignInProfilePageState> {
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
      this.setState({ error: result.error });
    }
  }

  public render() {
    return (
      <>
        <Modal.Window canClose={false} isOpen={true}>
          <Modal.Header>Complete your profile</Modal.Header>
          <Modal.Content>
            <p>Because this is your first sign in, please input your display name.</p>
            <Form2 error={this.state.error}>
              <Input
                field="name"
                onChange={name => this.setState({ name })}
                maxLength={100}
                placeholder="Your display name"
                suffix={
                  <Button onClick={() => this.submit()} color="positive" disabled={this.state.name === ""}>
                    Submit
                  </Button>
                }
              />
            </Form2>
          </Modal.Content>
        </Modal.Window>
        {this.props.tenant.isPrivate
          ? React.createElement(SignInPage, this.props)
          : React.createElement(HomePage, this.props)}
      </>
    );
  }
}
