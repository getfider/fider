import * as React from 'react';
import { Header, Footer, Button, DisplayError } from '@fider/components/common';

import { User } from '@fider/models';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

export class MembersPage extends React.Component<{}, {}> {
    private allUsers: User[];

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);
        this.allUsers = this.session.get<User[]>('users') || [];
    }

    public render() {
      return <div>
                <Header />
                    <div className="page ui container">
                        <h1 className="ui header">Administration</h1>
                        <h4 className="ui header">Members</h4>

                        <div className="ui grid">
                            <div className="eight wide computer sixteen wide mobile column">
                                <div className="ui segment">
                                    <h4 className="ui header">Administrators</h4>
                                    <p className="info">Administrators have full access to edit and manage content, permissions and settings.</p>
                                    {
                                        this.allUsers.map((x) => {
                                            return <h1 key={x.id}>{ x.name }</h1>;
                                        })
                                    }
                                </div>
                            </div>
                            <div className="eight wide computer sixteen wide mobile column">
                                <div className="ui segment">
                                    <h4 className="ui header">Collaborators</h4>
                                    <p className="info">Collaborators can edit and manage content, but not permissions and settings.</p>
                                    {
                                        this.allUsers.map((x) => {
                                            return <h1 key={x.id}>{ x.name }</h1>;
                                        })
                                    }
                                </div>
                            </div>
                        </div>
                    </div>
                <Footer />
            </div>;
    }
}
