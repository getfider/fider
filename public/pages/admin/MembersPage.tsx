import * as React from 'react';
import { Header, Footer, Button, Gravatar, UserName } from '@fider/components/common';

import { User, CurrentUser, UserRole } from '@fider/models';

import { inject, injectables } from '@fider/di';
import { Session, TenantService } from '@fider/services';

interface MembersPageState {
    administrators: User[];
    collaborators: User[];
    visitors: User[];
    filteredNewAdministrators: User[];
    filteredNewCollaborators: User[];
    newAdministratorFilter: string;
    newCollaboratorFilter: string;
}

export class MembersPage extends React.Component<{}, MembersPageState> {
    private allUsers: User[];
    private currentUser: CurrentUser;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.TenantService)
    public tenantService: TenantService;

    constructor(props: {}) {
        super(props);
        this.allUsers = this.session.get<User[]>('users') || [];
        this.currentUser = this.session.getCurrentUser()!;
        this.state = this.groupUsers();
    }

    private async changeRole(user: User, role: UserRole): Promise<any> {
        const response = await this.tenantService.changeRole(user.id, role);
        if (response.ok) {
            user.role = role;
            this.setState(this.groupUsers());
        }
    }

    private groupUsers(): MembersPageState {
        const usersByRole = this.allUsers.reduce<{ [key: number]: User[] }>((groups, x) => {
            groups[x.role] = [ x ].concat(groups[x.role] || []);
            return groups;
        }, {});

        return {
            administrators: usersByRole[UserRole.Administrator] || [],
            collaborators: usersByRole[UserRole.Collaborator] || [],
            visitors: usersByRole[UserRole.Visitor] || [],
            filteredNewAdministrators: [],
            filteredNewCollaborators: [],
            newAdministratorFilter: '',
            newCollaboratorFilter: '',
        };
    }

    private showUser(user: User, role: UserRole, addable: boolean, removable: boolean) {
        if (user.id === this.currentUser.id || this.currentUser.role !== UserRole.Administrator) {
            removable = false;
        }

        return <div key={user.id} className="item">
                   <Gravatar user={user} />
                   <div className="content">
                       <UserName user={user} />
                   </div>
                   <div className="right floated content">
                        {
                            removable &&
                            <Button size="mini" onClick={() => this.changeRole(user, UserRole.Visitor)} className="red showover">
                                <i className="remove icon"></i>
                                Remove
                            </Button>
                        }
                        {
                            addable &&
                            <Button size="mini" onClick={() => this.changeRole(user, role)} className="green showover">
                                <i className="add icon"></i>
                                Add
                            </Button>
                        }
                   </div>
               </div>;
    }

    private filterVisitors(property: string, text: string) {
        let filtered: User[] = [];
        if (text) {
            filtered = this.state.visitors.filter((x) => x.name.toLowerCase().indexOf(text.toLowerCase()) >= 0);
        }

        if (property === 'administrator') {
            this.setState({
                newAdministratorFilter: text,
                filteredNewAdministrators: filtered
            });
        } else if (property === 'collaborator') {
            this.setState({
                newCollaboratorFilter: text,
                filteredNewCollaborators: filtered
            });
        }
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
                                        { this.state.administrators.map((x) => this.showUser(x, UserRole.Administrator, false, true)) }
                                    </div>
                                    { this.currentUser.role === UserRole.Administrator && <div className="ui mini form">
                                        <p>Add new administrator</p>
                                        <div className="mini field">
                                            <input
                                                type="text"
                                                value={this.state.newAdministratorFilter}
                                                onChange={(x) => this.filterVisitors('administrator', x.currentTarget.value)}
                                                placeholder="Search users by name"
                                            />
                                        </div>
                                        <div className="ui middle aligned very relaxed selection list">
                                            { this.state.filteredNewAdministrators.map((x) => this.showUser(x, UserRole.Administrator, true, false)) }
                                        </div>
                                        {
                                            this.state.newAdministratorFilter &&
                                            this.state.filteredNewAdministrators.length === 0 &&
                                            <p className="info">No users to show.</p>
                                        }
                                    </div> }
                                </div>
                            </div>

                            <div className="eight wide computer sixteen wide mobile column">
                                <div className="ui segment">
                                    <h4 className="ui header">Collaborators</h4>
                                    <p className="info">Collaborators can edit and manage content, but not permissions and settings.</p>
                                    <div className="ui middle aligned very relaxed selection list">
                                        { this.state.collaborators.map((x) => this.showUser(x, UserRole.Collaborator, false, true)) }
                                    </div>
                                    { this.currentUser.role === UserRole.Administrator && <div className="ui mini form">
                                        <p>Add new collaborator</p>
                                        <div className="mini field">
                                            <input
                                                type="text"
                                                value={this.state.newCollaboratorFilter}
                                                onChange={(x) => this.filterVisitors('collaborator', x.currentTarget.value)}
                                                placeholder="Search users by name"
                                            />
                                        </div>
                                        <div className="ui middle aligned very relaxed selection list">
                                            { this.state.filteredNewCollaborators.map((x) => this.showUser(x, UserRole.Collaborator, true, false)) }
                                        </div>
                                        {
                                            this.state.newCollaboratorFilter &&
                                            this.state.filteredNewCollaborators.length === 0 &&
                                            <p className="info">No users to show.</p>
                                        }
                                    </div> }
                                </div>
                            </div>

                        </div>
                    </div>
                <Footer />
            </div>;
    }
}
