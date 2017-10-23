import * as React from 'react';
import { Header, Footer, Button, Gravatar, UserName } from '@fider/components/common';

import { User, UserRole } from '@fider/models';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

interface MembersPageState {
    administrators: User[];
    collaborators: User[];
    visitors: User[];
    newAdministrators: User[];
}

export class MembersPage extends React.Component<{}, MembersPageState> {
    private allUsers: User[];

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);
        this.allUsers = this.session.get<User[]>('users') || [];

        const usersByRole = this.allUsers.reduce<{ [key: number]: User[] }>((groups, x) => {
            groups[x.role] = [ x ].concat(groups[x.role] || []);
            return groups;
        }, {});

        this.state = {
            administrators: usersByRole[UserRole.Administrator] || [],
            collaborators: usersByRole[UserRole.Collaborator] || [],
            visitors: usersByRole[UserRole.Visitor] || [],
            newAdministrators: [],
        };
    }

    private showUser(user: User, addable: boolean, removable: boolean) {
        return <div key={user.id} className="item">
                   <Gravatar user={user} />
                   <div className="content">
                       <UserName user={user} />
                   </div>
                   <div className="right floated content">
                        {
                            removable &&
                            <button className="ui mini button red">
                                <i className="remove icon"></i>
                                Remove
                            </button>
                        }
                        {
                            addable &&
                            <button className="ui mini button green">
                                <i className="add icon"></i>
                                Add
                            </button>
                        }
                   </div>
               </div>;
    }

    private filterVisitors(text: string) {
        let newAdministrators: User[] = [];
        if (text) {
            newAdministrators = this.state.visitors.filter((x) => x.name.toLowerCase().indexOf(text.toLowerCase()) >= 0);
        }

        this.setState({
            newAdministrators
        });
    }

    public render() {
      return <div>
                <Header />
                    <div className="page ui container">
                        <h2 className="ui header">
                            <i className="circular users icon"></i>
                            <div className="content">
                              Members
                              <div className="sub header">Manage your account administrators and collaborators.</div>
                            </div>
                        </h2>

                        <div className="ui grid">
                            <div className="eight wide computer sixteen wide mobile column">
                                <div className="ui segment">
                                    <h4 className="ui header">Administrators</h4>
                                    <p className="info">Administrators have full access to edit and manage content, permissions and settings.</p>
                                    <div className="ui middle aligned very relaxed selection list">
                                        { this.state.administrators.map((x) => this.showUser(x, false, true)) }
                                    </div>
                                    <div className="ui mini form">
                                        <p>Add new administrators</p>
                                        <div className="mini field">
                                            <input type="text" onChange={(x) => this.filterVisitors(x.currentTarget.value)} placeholder="Search users by name"/>
                                        </div>
                                    </div>
                                    <div className="ui middle aligned very relaxed selection list">
                                        { this.state.newAdministrators.map((x) => this.showUser(x, true, false)) }
                                    </div>
                                </div>
                            </div>
                            <div className="eight wide computer sixteen wide mobile column">
                                <div className="ui segment">
                                    <h4 className="ui header">Collaborators</h4>
                                    <p className="info">Collaborators can edit and manage content, but not permissions and settings.</p>
                                    { this.state.collaborators.map((x) => this.showUser(x, false, true)) }
                                </div>
                            </div>
                        </div>
                    </div>
                <Footer />
            </div>;
    }
}
