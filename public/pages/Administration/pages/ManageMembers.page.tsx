import "./ManageMembers.page.scss";

import React from "react";
import { Segment, List, Input, ListItem, Gravatar, UserName, DropDown } from "@fider/components/common";
import { User, UserRole, UserStatus } from "@fider/models";
import { AdminBasePage } from "../components/AdminBasePage";
import { FaUsers, FaEllipsisH, FaTimes, FaSearch } from "react-icons/fa";

interface ManageMembersPageState {
  query: string;
  users: User[];
  visibleUsers: User[];
}

interface ManageMembersPageProps {
  users: User[];
}

const UserListItem = (props: { user: User }) => {
  const admin = props.user.role === UserRole.Administrator && <span className="staff">administrator</span>;
  const collaborator = props.user.role === UserRole.Collaborator && <span className="staff">collaborator</span>;
  const blocked = props.user.status === UserStatus.Blocked && <span className="blocked">blocked</span>;
  const isVisitor = props.user.role === UserRole.Visitor;

  const renderEllipsis = () => {
    return <FaEllipsisH />;
  };

  return (
    <ListItem>
      <Gravatar user={props.user} />
      <div className="l-user-details">
        <UserName user={props.user} />
        <span>
          {admin} {collaborator} {blocked}
        </span>
      </div>
      <DropDown
        className="l-user-actions"
        direction="left"
        highlightSelected={false}
        items={[
          (!!collaborator || isVisitor) && { label: "Promote to Administrator", value: "to-administrator" },
          (!!admin || isVisitor) && { label: "Promote to Collaborator", value: "to-collaborator" },
          (!!collaborator || !!admin) && { label: "Demote to Visitor", value: "to-visitor" },
          !blocked && { label: "Block User", value: "block" },
          !!blocked && { label: "Unblock User", value: "unblock" }
        ]}
        renderText={renderEllipsis}
      />
    </ListItem>
  );
};

export class ManageMembersPage extends AdminBasePage<ManageMembersPageProps, ManageMembersPageState> {
  public id = "p-admin-members";
  public name = "members";
  public icon = FaUsers;
  public title = "Members";
  public subtitle = "Manage your site administrators and collaborators";

  constructor(props: ManageMembersPageProps) {
    super(props);
    this.state = {
      query: "",
      users: this.props.users,
      visibleUsers: this.props.users.slice(0, 10)
    };
  }

  private showMore = (event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void => {
    event.preventDefault();
    this.setState({
      visibleUsers: this.state.users.slice(0, this.state.users.length + 10)
    });
  };

  private clearSearch = () => {
    this.setState({ query: "" });
  };

  private handleSearchFilterChanged = (query: string) => {
    const users = this.props.users.filter(x => x.name.toLowerCase().indexOf(query.toLowerCase()) >= 0);
    this.setState({ query, users, visibleUsers: users.slice(0, 10) });
  };

  private chunks = () => {
    const [col1, col2] = [[], []] as User[][];
    this.state.visibleUsers.forEach((user, index) => {
      if (index % 2 === 0) {
        col1.push(user);
      } else {
        col2.push(user);
      }
    });
    return [col1, col2];
  };

  public content() {
    const [col1, col2] = this.chunks();

    return (
      <>
        <Input
          field="query"
          icon={this.state.query ? FaTimes : FaSearch}
          onIconClick={this.state.query ? this.clearSearch : undefined}
          placeholder="Search for users by name..."
          value={this.state.query}
          onChange={this.handleSearchFilterChanged}
        />
        <Segment>
          <div className="row">
            <div className="col-lg-6 col-left">
              <List divided={true}>
                {col1.map(user => (
                  <UserListItem key={user.id} user={user} />
                ))}
              </List>
            </div>
            <div className="col-lg-6 col-right">
              <List divided={true}>
                {col2.map(user => (
                  <UserListItem key={user.id} user={user} />
                ))}
              </List>
            </div>
          </div>
        </Segment>
        <p className="info">
          {!this.state.query && (
            <>
              Showing {this.state.visibleUsers.length} of {this.state.users.length} registered users
            </>
          )}
          {this.state.query && (
            <>
              Showing {this.state.visibleUsers.length} of {this.state.users.length} users matching '{this.state.query}'
            </>
          )}
          {this.state.visibleUsers.length < this.state.users.length && (
            <a className="l-show-more" onTouchEnd={this.showMore} onClick={this.showMore}>
              view more
            </a>
          )}
        </p>
      </>
    );
  }
}
