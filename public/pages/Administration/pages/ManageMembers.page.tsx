import "./ManageMembers.page.scss";

import * as React from "react";
import { Button, Gravatar, UserName, Segment, ListItem, List, Input, Form } from "@fider/components/common";
import { User, CurrentUser, UserRole } from "@fider/models";
import { actions } from "@fider/services";
import { AdminBasePage, UserListItem } from "../components";

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

  private onRoleChanged = () => {
    this.setState(this.groupUsers());
  };

  private showUser(user: User, role: UserRole, addable: boolean, removable: boolean) {
    if (user.id === Fider.session.user.id || Fider.session.user.role !== UserRole.Administrator) {
      removable = false;
    }

    return (
      <UserListItem
        key={user.id}
        role={role}
        user={user}
        addable={addable}
        removable={removable}
        onChange={this.onRoleChanged}
      />
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

  private handleSearch = {
    administrator: (query: string) => {
      this.filterVisitors("administrator", query);
    },
    collaborator: (query: string) => {
      this.filterVisitors("collaborator", query);
    }
  };

  public content() {
    return (
      <div className="row">
        <div className="col-lg-6">
          <Segment>
            <h4>Administrators</h4>
            <p className="info">
              Administrators have full access to edit and manage content, permissions and settings.
            </p>
            <List hover={true}>
              {this.state.administrators.map(x => this.showUser(x, UserRole.Administrator, false, true))}
            </List>
            {Fider.session.user.isAdministrator && (
              <Form size="mini">
                <Input
                  label="Add new administrator"
                  field="new-administrator"
                  value={this.state.newAdministratorFilter}
                  onChange={this.handleSearch.administrator}
                  placeholder="Search users by name"
                />
                <List hover={true}>
                  {this.state.filteredNewAdministrators.map(x => this.showUser(x, UserRole.Administrator, true, false))}
                </List>
                {this.state.newAdministratorFilter &&
                  this.state.filteredNewAdministrators.length === 0 && <p className="info">No users to show.</p>}
              </Form>
            )}
          </Segment>
        </div>

        <div className="col-lg-6">
          <Segment>
            <h4>Collaborators</h4>
            <p className="info">Collaborators can edit and manage content, but not permissions and settings.</p>
            <List hover={true}>
              {this.state.collaborators.map(x => this.showUser(x, UserRole.Collaborator, false, true))}
            </List>
            {Fider.session.user.isAdministrator && (
              <Form size="mini">
                <Input
                  label="Add new collaborator"
                  field="new-collaborator"
                  value={this.state.newCollaboratorFilter}
                  onChange={this.handleSearch.collaborator}
                  placeholder="Search users by name"
                />
                <List hover={true}>
                  {this.state.filteredNewCollaborators.map(x => this.showUser(x, UserRole.Collaborator, true, false))}
                </List>
                {this.state.newCollaboratorFilter &&
                  this.state.filteredNewCollaborators.length === 0 && <p className="info">No users to show.</p>}
              </Form>
            )}
          </Segment>
        </div>
      </div>
    );
  }
}
