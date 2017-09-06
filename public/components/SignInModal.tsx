import * as React from 'react';
import { SocialSignInList } from '@fider/components/common';

export class SignInModal extends React.Component<{}, {}> {
    public render() {
        return <div id="signin-modal" className="ui modal small">
                    <div className="header">
                      Log in to raise your voice.
                    </div>
                    <div className="content">
                        <SocialSignInList size="normal" />
                    </div>
               </div>;
    }
}
