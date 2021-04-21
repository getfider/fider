import "./ManageMembers.page.scss"

import React from "react"
import { Segment, List, Input, ListItem, Avatar, UserName, DropDown, DropDownItem } from "@fider/components/common"
import { User, UserRole, UserStatus } from "@fider/models"
import { AdminBasePage } from "../components/AdminBasePage"
import { FaUsers, FaTimes, FaSearch } from "react-icons/fa"
import { actions, Fider } from "@fider/services"

interface ManageMembersPageState {
  query: string
  users: User[]
  visibleUsers: User[]
}

interface ManageMembersPageProps {
  users: User[]
}

interface UserListItemProps {
  user: User
  onAction: (actionName: string, user: User) => Promise<void>
}

const UserListItem = (props: UserListItemProps) => {
  const admin = props.user.role === UserRole.Administrator && <span className="staff">administrator</span>
  const collaborator = props.user.role === UserRole.Collaborator && <span className="staff">collaborator</span>
  const blocked = props.user.status === UserStatus.Blocked && <span className="blocked">blocked</span>
  const isVisitor = props.user.role === UserRole.Visitor

  const renderEllipsis = () => {
    return (
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" width="16" height="16" focusable="false">
        <path d="M3 9.5A1.5 1.5 0 114.5 8 1.5 1.5 0 013 9.5zM11.5 8A1.5 1.5 0 1013 6.5 1.5 1.5 0 0011.5 8zm-5 0A1.5 1.5 0 108 6.5 1.5 1.5 0 006.5 8z"></path>
      </svg>
    )
  }

  const actionSelected = (item: DropDownItem) => {
    props.onAction(item.value, props.user)
  }

  return (
    <ListItem>
      <Avatar user={props.user} />
      <div className="l-user-details">
        <UserName user={props.user} />
        <span>
          {admin} {collaborator} {blocked}
        </span>
      </div>
      {Fider.session.user.id !== props.user.id && Fider.session.user.isAdministrator && (
        <DropDown
          className="l-user-actions"
          inline={true}
          highlightSelected={false}
          style="simple"
          items={[
            !blocked && (!!collaborator || isVisitor) && { label: "Promote to Administrator", value: "to-administrator" },
            !blocked && (!!admin || isVisitor) && { label: "Promote to Collaborator", value: "to-collaborator" },
            !blocked && (!!collaborator || !!admin) && { label: "Demote to Visitor", value: "to-visitor" },
            isVisitor && !blocked && { label: "Block User", value: "block" },
            isVisitor && !!blocked && { label: "Unblock User", value: "unblock" },
          ]}
          renderControl={renderEllipsis}
          onChange={actionSelected}
        />
      )}
    </ListItem>
  )
}

export default class ManageMembersPage extends AdminBasePage<ManageMembersPageProps, ManageMembersPageState> {
  public id = "p-admin-members"
  public name = "members"
  public icon = FaUsers
  public title = "Members"
  public subtitle = "Manage your site administrators and collaborators"

  constructor(props: ManageMembersPageProps) {
    super(props)

    const users = this.props.users.sort(this.sortByStaff)
    this.state = {
      query: "",
      users,
      visibleUsers: users.slice(0, 10),
    }
  }

  private showMore = (event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void => {
    event.preventDefault()
    this.setState({
      visibleUsers: this.state.users.slice(0, this.state.visibleUsers.length + 10),
    })
  }

  private clearSearch = () => {
    this.handleSearchFilterChanged("")
  }

  private handleSearchFilterChanged = (query: string) => {
    const users = this.props.users.filter((x) => x.name.toLowerCase().indexOf(query.toLowerCase()) >= 0).sort(this.sortByStaff)
    this.setState({ query, users, visibleUsers: users.slice(0, 10) })
  }

  private handleAction = async (actionName: string, user: User) => {
    const changeRole = async (role: UserRole) => {
      const result = await actions.changeUserRole(user.id, role)
      if (result.ok) {
        user.role = role
      }
      this.handleSearchFilterChanged(this.state.query)
    }

    const changeStatus = async (status: UserStatus) => {
      const action = status === UserStatus.Blocked ? actions.blockUser : actions.unblockUser
      const result = await action(user.id)
      if (result.ok) {
        user.status = status
      }
      this.forceUpdate()
    }

    if (actionName === "to-collaborator") {
      await changeRole(UserRole.Collaborator)
    } else if (actionName === "to-visitor") {
      await changeRole(UserRole.Visitor)
    } else if (actionName === "to-administrator") {
      await changeRole(UserRole.Administrator)
    } else if (actionName === "block") {
      await changeStatus(UserStatus.Blocked)
    } else if (actionName === "unblock") {
      await changeStatus(UserStatus.Active)
    }
  }

  private sortByStaff = (left: User, right: User) => {
    if (right.role === left.role) {
      if (left.name < right.name) {
        return -1
      } else if (left.name > right.name) {
        return 1
      }
      return 0
    }

    if (right.role !== UserRole.Visitor) {
      return 1
    }
    return -1
  }

  public content() {
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
          {this.state.visibleUsers.length === 0 && <span>No users found.</span>}
          <List divided={true}>
            {this.state.visibleUsers.map((user) => (
              <UserListItem key={user.id} user={user} onAction={this.handleAction} />
            ))}
          </List>
        </Segment>
        <p className="info">
          {!this.state.query && (
            <>
              Showing {this.state.visibleUsers.length} of {this.state.users.length} registered users
            </>
          )}
          {this.state.query && (
            <>
              Showing {this.state.visibleUsers.length} of {this.state.users.length} users matching &apos;{this.state.query}&apos;
            </>
          )}
          {this.state.visibleUsers.length < this.state.users.length && (
            <a className="l-show-more" onTouchEnd={this.showMore} onClick={this.showMore}>
              view more
            </a>
          )}
        </p>
        <ul className="l-legend info">
          <li>
            <strong>&middot; Administrators</strong>have full access to edit and manage content, permissions and all site settings.
          </li>
          <li>
            <strong>&middot; Collaborators</strong> can edit and manage content, but not permissions and settings.
          </li>
          <li>
            <strong>&middot; Blocked</strong> users are unable to log into this site.
          </li>
        </ul>
      </>
    )
  }
}
