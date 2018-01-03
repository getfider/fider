import * as React from 'react';
import { SignInControl } from '@fider/components/common';
import { AuthSettings } from '@fider/models';

interface SignInModalProps {
  auth: AuthSettings;
}

export class SignInModal extends React.Component<SignInModalProps, {}> {
  public render() {
    return (
      <div id="signin-modal" className="ui modal small">
          <div className="header">
            Sign in to raise your voice
          </div>
          <div className="content">
              <SignInControl auth={this.props.auth} signInByEmail={true} />
          </div>
      </div>
    );
  }
}
