import * as React from "react";
import { Button, Gravatar, UserName } from "@fider/components/common";
import { User, CurrentUser, UserRole } from "@fider/models";
import { actions } from "@fider/services";
import { AdminBasePage } from "../components";

interface ManageMembersPageState {
  administrators: User[];
  collaborators: User[];
  visitors: User[];
  filteredNewAdministrators: User[];
  filteredNewCollaborators: User[];
  newAdministratorFilter: string;
  newCollaboratorFilter: string;
}

interface ManageMembersPageProps {
  user: CurrentUser;
  users: User[];
}

export class ManageMembersPage extends AdminBasePage<ManageMembersPageProps, ManageMembersPageState> {
  public id = "p-admin-members";
  public name = "members";
  public icon = "users";
  public title = "Members";
  public subtitle = "Manage your site administrators and collaborators";

  constructor(props: ManageMembersPageProps) {
    super(props);
    this.state = this.groupUsers();
  }

  private async changeRole(user: User, role: UserRole): Promise<any> {
    const response = await actions.changeUserRole(user.id, role);
    if (response.ok) {
      user.role = role;
      this.setState(this.groupUsers());
    }
  }

  private groupUsers(): ManageMembersPageState {
    const usersByRole = this.props.users.reduce<{ [key: number]: User[] }>((groups, x) => {
      groups[x.role] = [x].concat(groups[x.role] || []);
      return groups;
    }, {});

    return {
      administrators: usersByRole[UserRole.Administrator] || [],
      collaborators: usersByRole[UserRole.Collaborator] || [],
      visitors: usersByRole[UserRole.Visitor] || [],
      filteredNewAdministrators: [],
      filteredNewCollaborators: [],
      newAdministratorFilter: "",
      newCollaboratorFilter: ""
    };
  }

  private showUser(user: User, role: UserRole, addable: boolean, removable: boolean) {
    if (user.id === this.props.user.id || this.props.user.role !== UserRole.Administrator) {
      removable = false;
    }

    return (
      <div key={user.id} className="item">
        <Gravatar user={user} />
        <div className="content">
          <UserName user={user} />
        </div>
        <div className="right floated content">
          {removable && (
            <Button
              size="tiny"
              color="danger"
              onClick={() => this.changeRole(user, UserRole.Visitor)}
              className="showover"
            >
              <i className="remove icon" />Remove
            </Button>
          )}
          {addable && (
            <Button size="tiny" color="positive" onClick={() => this.changeRole(user, role)} className="showover">
              <i className="add icon" />Add
            </Button>
          )}
        </div>
      </div>
    );
  }

  private filterVisitors(property: string, text: string) {
    let filtered: User[] = [];
    if (text) {
      filtered = this.state.visitors.filter(x => x.name.toLowerCase().indexOf(text.toLowerCase()) >= 0);
    }

    if (property === "administrator") {
      this.setState({
        newAdministratorFilter: text,
        filteredNewAdministrators: filtered
      });
    } else if (property === "collaborator") {
      this.setState({
        newCollaboratorFilter: text,
        filteredNewCollaborators: filtered
      });
    }
  }

  public content() {
    return (
      <div className="ui grid">
        <div className="eight wide computer sixteen wide mobile column">
          <div className="ui segment">
            <h4 className="ui header">Administrators</h4>
            <p className="info">
              Administrators have full access to edit and manage content, permissions and settings.
            </p>
            <div className="ui middle aligned very relaxed selection list">
              {this.state.administrators.map(x => this.showUser(x, UserRole.Administrator, false, true))}
            </div>
            {this.props.user.role === UserRole.Administrator && (
              <div className="ui mini form">
                <p>Add new administrator</p>
                <div className="mini field">
                  <input
                    type="text"
                    value={this.state.newAdministratorFilter}
                    onChange={x => this.filterVisitors("administrator", x.currentTarget.value)}
                    placeholder="Search users by name"
                  />
                </div>
                <div className="ui middle aligned very relaxed selection list">
                  {this.state.filteredNewAdministrators.map(x => this.showUser(x, UserRole.Administrator, true, false))}
                </div>
                {this.state.newAdministratorFilter &&
                  this.state.filteredNewAdministrators.length === 0 && <p className="info">No users to show.</p>}
              </div>
            )}
          </div>
        </div>

        <div className="eight wide computer sixteen wide mobile column">
          <div className="ui segment">
            <h4 className="ui header">Collaborators</h4>
            <p className="info">Collaborators can edit and manage content, but not permissions and settings.</p>
            <div className="ui middle aligned very relaxed selection list">
              {this.state.collaborators.map(x => this.showUser(x, UserRole.Collaborator, false, true))}
            </div>
            {this.props.user.role === UserRole.Administrator && (
              <div className="ui mini form">
                <p>Add new collaborator</p>
                <div className="mini field">
                  <input
                    type="text"
                    value={this.state.newCollaboratorFilter}
                    onChange={x => this.filterVisitors("collaborator", x.currentTarget.value)}
                    placeholder="Search users by name"
                  />
                </div>
                <div className="ui middle aligned very relaxed selection list">
                  {this.state.filteredNewCollaborators.map(x => this.showUser(x, UserRole.Collaborator, true, false))}
                </div>
                {this.state.newCollaboratorFilter &&
                  this.state.filteredNewCollaborators.length === 0 && <p className="info">No users to show.</p>}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }
}
