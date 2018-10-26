import React from "react";
import { User, UserRole } from "@fider/models";
import { ListItem, Gravatar, UserName, Button } from "@fider/components";
import { actions } from "@fider/services";

interface UserListItemProps {
  user: User;
  role: UserRole;
  removable: boolean;
  addable: boolean;
  onChange: () => void;
}

export class UserListItem extends React.Component<UserListItemProps, {}> {
  private promote = async () => {
    await this.changeRole(this.props.role);
  };

  private demote = async () => {
    await this.changeRole(UserRole.Visitor);
  };

  private changeRole = async (role: UserRole) => {
    const response = await actions.changeUserRole(this.props.user.id, role);
    if (response.ok) {
      this.props.user.role = role;
      this.props.onChange();
    }
  };

  public render() {
    return (
      <ListItem>
        <Gravatar user={this.props.user} />
        <div className="content">
          <UserName user={this.props.user} />
        </div>
        {this.props.removable && (
          <Button size="tiny" color="danger" onClick={this.demote} className="right showover">
            <i className="remove icon" />
            Remove
          </Button>
        )}
        {this.props.addable && (
          <Button size="tiny" color="positive" onClick={this.promote} className="right showover">
            <i className="add icon" />
            Add
          </Button>
        )}
      </ListItem>
    );
  }
}
