import * as React from 'react';
import {SocialSignInList } from '@fider/components/common';

export class SignInModal extends React.Component<{}, {}> {
    public render() {
        return <div id="signin-modal" className="ui modal tiny">
                    <div className="header">
                      Log in to raise your voice.
                    </div>
                    <div className="content">
                        <p className="info">We'll never post to any of your accounts.</p>
                        <SocialSignInList orientation="horizontal" size="normal" />
                    </div>
               </div>;
    }
}
