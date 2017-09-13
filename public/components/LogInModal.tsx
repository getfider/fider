import * as React from 'react';
import { LogInControl } from '@fider/components/common';

export class LogInModal extends React.Component<{}, {}> {
    public render() {
        return <div id="login-modal" className="ui modal small">
                    <div className="header">
                      Log in to raise your voice
                    </div>
                    <div className="content">
                        <LogInControl />
                    </div>
               </div>;
    }
}
