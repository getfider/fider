import * as React from 'react';
import { SignInControl } from '@fider/components/common';

export class SignInModal extends React.Component<{}, {}> {
    public render() {
        return <div id="signin-modal" className="ui modal small">
                    <div className="header">
                      Sign in to raise your voice
                    </div>
                    <div className="content">
                        <SignInControl />
                    </div>
               </div>;
    }
}
